package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"strings"
	"testing"
)

func TestInviteMemberNotAdmin(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)

	input := []struct {
		Orga   *model.Organization
		UserId hide.ID
		Emails []string
	}{
		{orga, 0, []string{}},
	}

	for _, in := range input {
		if err := InviteMember(in.Orga, in.UserId, InviteMemberData{in.Emails, false}, testMailSender); err != errs.PermissionDenied {
			t.Fatalf("Expected error '%v', but was: %v", errs.PermissionDenied, err)
		}
	}
}

func TestInviteMemberInvalidEmail(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	input := []struct {
		Orga   *model.Organization
		UserId hide.ID
		Emails []string
	}{
		{orga, user.ID, []string{}},
		{orga, user.ID, []string{"invalid", "invalid"}},
		{orga, user.ID, []string{"valid@mail.com", "invalid"}},
	}
	expected := []error{
		nil,
		errs.EmailInvalid,
		errs.EmailInvalid,
	}

	for i := range input {
		if err := InviteMember(input[i].Orga, input[i].UserId, InviteMemberData{input[i].Emails, false}, testMailSender); err != expected[i] {
			t.Fatalf("Expected error '%v', but was: %v", expected[i], err)
		}
	}
}

func TestInviteMemberSuccessNewUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	// their might be a second account for the same mail address, as you can sign up through SSO too!
	testutil.CreateUserWithoutOrganization(t, 322, "new@user.com")

	mails := []string{"new@user.com"}

	if err := InviteMember(orga, user.ID, InviteMemberData{mails, false}, testMailWasSend(t)); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	testInvitationWasCreated(t, orga)
}

func TestInviteMemberUserIsMemberAlready(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	existingUser := testutil.CreateUser(t, orga, 321, "new@user.com")
	mails := []string{existingUser.Email}

	if err := InviteMember(orga, user.ID, InviteMemberData{mails, false}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	if model.GetInvitationByOrganizationIdAndEmail(orga.ID, "new@user.com") == nil {
		t.Fatal("Invitation must have been created")
	}
}

func TestInviteMemberReadOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	mails := []string{"new@user.com"}

	if err := InviteMember(orga, user.ID, InviteMemberData{mails, true}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	invitation := model.GetInvitationByOrganizationIdAndEmail(orga.ID, "new@user.com")

	if invitation == nil {
		t.Fatalf("Invitation must have been created, but was: %v", invitation)
	}

	if !invitation.ReadOnly {
		t.Fatalf("Invitation must be read only, but was: %v", invitation.ReadOnly)
	}
}

func TestInviteMemberEntryReadOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	mails := []string{"new@user.com"}

	if err := InviteMember(orga, user.ID, InviteMemberData{mails, true}, testMailSender); err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	invitation := model.GetInvitationByOrganizationIdAndEmail(orga.ID, "new@user.com")

	if invitation == nil {
		t.Fatalf("Invitation must have been created, but was: %v", invitation)
	}

	if invitation.ReadOnly {
		t.Fatalf("Invitation must be not be read only, but was: %v", invitation.ReadOnly)
	}
}

func TestInviteMemberSSOUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateUserWithoutOrganization(t, 321, "existing@user.com")
	testutil.CreateUserWithoutOrganization(t, 322, "existing@user.com")
	emails := []string{"existing@user.com"}

	if err := InviteMember(orga, user.ID, InviteMemberData{emails, false}, testMailSender); err != nil {
		t.Fatalf("Expected invitation to be created, but was: %v", err)
	}

	invitation := model.GetInvitationByOrganizationIdAndEmail(orga.ID, "existing@user.com")

	if invitation == nil {
		t.Fatalf("Invitation must have been created, but was: %v", invitation)
	}
}

func TestInviteMemberSelf(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	emails := []string{user.Email}

	if err := InviteMember(orga, user.ID, InviteMemberData{emails, false}, testMailSender); err != nil {
		t.Fatalf("Expected no invitation to be created, but was: %v", err)
	}

	if model.GetInvitationByOrganizationIdAndEmail(orga.ID, user.Email) != nil {
		t.Fatal("Invitation must not have been created for the inviting member")
	}
}

func testMailWasSend(t *testing.T) mail.Sender {
	return func(subject, msgHTML, from string, to ...string) error {
		t.Log(msgHTML)

		if subject != i18n.GetMailTitle("en")["invite"] {
			t.Fatalf("Mail subject not as expected, was: %v", subject)
		}

		if len(to) != 0 && to[0] != "new@user.com" {
			t.Fatalf("Mail receiver not as expected, was: %v", to[0])
		}

		vars := inviteNewUserMailI18n["en"]

		if !strings.Contains(msgHTML, string(vars["title"])) {
			t.Fatalf("Mail body does not contain title, was: %v", msgHTML)
		}

		if !strings.Contains(msgHTML, string(vars["text"])) {
			t.Fatalf("Mail body does not contain text, was: %v", msgHTML)
		}

		if !strings.Contains(msgHTML, string(vars["link"])) {
			t.Fatalf("Mail body does not contain link, was: %v", msgHTML)
		}

		if !strings.Contains(msgHTML, string(vars["greeting"])) {
			t.Fatalf("Mail body does not contain greeting, was: %v", msgHTML)
		}

		if !strings.Contains(msgHTML, string(vars["goodbye"])) {
			t.Fatalf("Mail body does not contain goodbye, was: %v", msgHTML)
		}

		return nil
	}
}

func testInvitationWasCreated(t *testing.T, orga *model.Organization) {
	invitation := model.GetInvitationByOrganizationIdAndEmail(orga.ID, "new@user.com")

	if invitation == nil ||
		invitation.OrganizationId != orga.ID ||
		invitation.Email != "new@user.com" ||
		len(invitation.Code) != 32 {
		t.Fatalf("Invitation not as expected, was: %v", invitation)
	}
}

func testMailSender(subject, msgHTML, from string, to ...string) error {
	return nil
}
