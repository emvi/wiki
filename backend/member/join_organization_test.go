package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"strings"
	"testing"
)

func TestJoinOrganizationFailure(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	createTestInvitation(t, orga, "test@user.com", false)

	input := []struct {
		UserId   hide.ID
		Username string
		Code     string
	}{
		{0, "", ""},
		{user.ID, "", ""},
		{user.ID, "new_user", "invalid"},
		{user.ID, "new_user", "code"},
	}
	expected := []error{
		errs.UserNotFound,
		errs.InvitationNotFound,
		errs.InvitationNotFound,
		errs.IsMemberAlready,
	}

	for i, in := range input {
		if err := JoinOrganization(in.UserId, JoinOrganizationData{in.Username, in.Code, ""}, testMailSender); err != expected[i] {
			t.Fatalf("Expected error '%v', but was: %v", expected[i], err)
		}
	}
}

func TestJoinOrganizationSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	user := testutil.CreateUserWithoutOrganization(t, 321, "new@user.com")
	createTestInvitation(t, orga, "new@user.com", true)

	if err := JoinOrganization(user.ID, JoinOrganizationData{"new_user", "code", ""}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user.ID)

	if member == nil ||
		!member.Active ||
		member.Username != "new_user" ||
		!member.RecommendationMail ||
		!member.ShowCreateButton ||
		!member.ShowActionButtons ||
		!member.ShowNavigation {
		t.Fatalf("Member not as expected: %v", member)
	}

	all := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupAllName)
	readonly := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupReadOnlyName)

	if all == nil || readonly == nil {
		t.Fatal("Default user group must exist")
	}

	if all.MemberCount != 2 || readonly.MemberCount != 2 {
		t.Fatalf("Default groups member count must have been updated, but was: %v %v", all.MemberCount, readonly.MemberCount)
	}

	allMember := model.GetUserGroupMemberByGroupIdAndUserId(all.ID, user.ID)

	if allMember == nil {
		t.Fatal("User must be member of default group all")
	}

	readonlyMember := model.GetUserGroupMemberByGroupIdAndUserId(readonly.ID, user.ID)

	if readonlyMember == nil {
		t.Fatal("User must be member of default group readonly")
	}

	testutil.AssertFeedCreated(t, orga, "joined_organization")
}

func TestJoinOrganizationExistingMemberSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	user := testutil.CreateUserWithoutOrganization(t, 321, "new@user.com")
	member := &model.OrganizationMember{OrganizationId: orga.ID,
		UserId:     user.ID,
		LanguageId: lang.ID,
		Username:   "old_username",
		Active:     false}

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}

	feed := testutil.CreateFeed(t, orga, user, lang, true)
	createTestInvitation(t, orga, "new@user.com", false)

	if err := JoinOrganization(user.ID, JoinOrganizationData{"new_user", "code", ""}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	member = model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user.ID)

	if member == nil || !member.Active || member.Username != "new_user" {
		t.Fatalf("Expected existing member to be active and username updated, but was: %v", member)
	}

	feedAccess := model.GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotification(orga.ID, user.ID, feed.ID, true)

	if !feedAccess.Read {
		t.Fatal("Feed must be marked read when rejoining organization")
	}
}

func TestJoinOrganizationMemberSettings(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	user2 := testutil.CreateUserWithoutOrganization(t, 321, "user2@test.com")
	testutil.CreateInvitation(t, orga, "user2@test.com", "code", false)
	data := JoinOrganizationData{
		Username: "user2",
		Code:     "code",
	}

	if err := JoinOrganization(user2.ID, data, testMailSender); err != nil {
		t.Fatalf("Expected user to join organization, but was: %v", err)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user2.ID)

	if member == nil ||
		member.SendNotificationsInterval != defaultSendNotificationIntervalDays ||
		!member.RecommendationMail {
		t.Fatalf("Member not as expected: %v", member)
	}
}

func TestCheckValidUsername(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)

	input := []string{
		"",
		"ab",
		"wayTooooooooooLooooooooooooooooooooong",
		".invalid",
		"invalid_",
		"in..valid",
		"in__valid",
		"in--valid",
		"inv/alid",
	}
	expected := []error{
		errs.UsernameTooShort,
		errs.UsernameTooShort,
		errs.UsernameTooLong,
		errs.UsernameInvalid,
		errs.UsernameInvalid,
		errs.UsernameInvalid,
		errs.UsernameInvalid,
		errs.UsernameInvalid,
		errs.UsernameInvalid,
	}

	for i, in := range input {
		if err := checkValidUsername(orga.ID, in); err != expected[i] {
			t.Fatalf("Expected error '%v', but was: %v", errs.UsernameInvalid, err)
		}
	}

	if err := checkValidUsername(orga.ID, "testuser1"); err != errs.UsernameInUse {
		t.Fatalf("Expected error '%v', but was: %v", errs.UsernameInUse, err)
	}

	if err := checkValidUsername(orga.ID, "012valid_Us-ern.ame"); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}
}

func TestJoinOrganizationReadOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	user := testutil.CreateUserWithoutOrganization(t, 321, "new@user.com")
	createTestInvitation(t, orga, "new@user.com", true)

	if err := JoinOrganization(user.ID, JoinOrganizationData{"new_user", "code", ""}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUsername(orga.ID, "new_user")

	if member == nil {
		t.Fatal("Member must have been created")
	}

	if !member.ReadOnly {
		t.Fatal("Member must be read only")
	}
}

func TestJoinOrganizationUpdateGroupMemberCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	user2 := testutil.CreateUser(t, orga, 321, "new@user.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user2, false)

	if err := RemoveMember(orga, user.ID, user2.ID, false); err != nil {
		t.Fatal(err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 0 {
		t.Fatalf("User group must not have members anymore")
	}

	createTestInvitation(t, orga, user2.Email, false)

	if err := JoinOrganization(user2.ID, JoinOrganizationData{user2.OrganizationMember.Username, "code", ""}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 1 {
		t.Fatalf("User group must have one member")
	}
}

func TestJoinOrganizationInvitationLink(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	orga.InvitationCode = null.NewString("invitationcode", true)
	orga.InvitationReadOnly = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	testutil.CreateLang(t, orga, "en", "English", true)
	user := testutil.CreateUserWithoutOrganization(t, 321, "new@user.com")

	if err := JoinOrganization(user.ID, JoinOrganizationData{"new_user", "", "foobar"}, testMailSender); err != errs.OrganizationNotFound {
		t.Fatalf("Organization must not have been found, but was: %v", err)
	}

	if err := JoinOrganization(user.ID, JoinOrganizationData{"new_user", "", "invitationcode"}, testMailSender); err != nil {
		t.Fatalf("Must have joined organization, but was: %v", err)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUsername(orga.ID, "new_user")

	if member == nil {
		t.Fatal("Member must have been created")
	}

	if !member.ReadOnly {
		t.Fatal("Member must be read only")
	}
}

func TestInformAdminsMemberJoined(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, constants.GroupAdminName)
	testutil.CreateUserGroupMember(t, group, admin, false)
	user := testutil.CreateUserWithoutOrganization(t, 321, "new@user.com")
	mailer := func(subject, msgHTML, from string, to ...string) error {
		if len(subject) == 0 || len(msgHTML) == 0 {
			t.Fatalf("Mail subject and content not as expected: %v %v", subject, msgHTML)
		}

		// sender is the receiver as we use the default sender in production, so "to" won't be set
		if from != "test@user.com" {
			t.Fatalf("Sender not as expected: %v", from)
		}

		if len(to) != 0 {
			t.Fatalf("Receiver must not be set, but was: %v", to)
		}

		return nil
	}
	informAdminsMemberJoined(orga.ID, user, mailer)
}

func TestRenderMemberJoinedMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, constants.GroupAdminName)
	testutil.CreateUserGroupMember(t, group, admin, false)
	user := testutil.CreateUserWithoutOrganization(t, 321, "new@user.com")
	user.Firstname = "Max"
	user.Lastname = "Mustermann"

	if err := model.SaveUser(nil, user, false); err != nil {
		t.Fatal(err)
	}

	subject, body, err := renderMemberJoinedMail(orga, user)

	if err != nil {
		t.Fatalf("Mail must have been rendered, but was: %v", err)
	}

	if subject != "Max Mustermann joined your organization at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, user.Firstname) ||
		!strings.Contains(body, user.Lastname) ||
		!strings.Contains(body, orga.Name) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["title-1"])) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["title-2"])) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["joined-1"])) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["joined-2"])) ||
		!strings.Contains(body, string(memberJoinedMailI18n["en"]["text"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func createTestInvitation(t *testing.T, orga *model.Organization, email string, readOnly bool) {
	invitation := &model.Invitation{OrganizationId: orga.ID,
		Email:    email,
		Code:     "code",
		ReadOnly: readOnly}

	if err := model.SaveInvitation(nil, invitation); err != nil {
		t.Fatal(err)
	}
}
