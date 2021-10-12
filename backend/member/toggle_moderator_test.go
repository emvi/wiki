package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestToggleModerator(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@testutil.com")
	testDeleteAllUserGroupMember(t)
	member := model.GetOrganizationMemberByUsername("testuser2")
	member.IsAdmin = true

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}

	input := []struct {
		Orga         *model.Organization
		UserId       hide.ID
		MemberUserId hide.ID
	}{
		{orga, 0, 0},
		{orga, user.ID, 0},
		{orga, user.ID, member.UserId + 1},
		{orga, user.ID, user.ID},
		{orga, user2.ID, member.UserId},
		{orga, user.ID, member.UserId},
	}
	expected := []struct {
		Error error
	}{
		{errs.PermissionDenied},
		{errs.MemberNotFound},
		{errs.MemberNotFound},
		{errs.OrganizationOwner},
		{errs.ModeratorYourself},
		{nil},
	}

	for i := range input {
		if err := ToggleModerator(input[i].Orga, input[i].UserId, input[i].MemberUserId); err != expected[i].Error {
			t.Fatalf("Expected error '%v', but was: %v", expected[i].Error, err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "set_organization_moderator")

	if !testIsGroupMember(t, orga, member.UserId, constants.GroupModName) {
		t.Fatal("User must be member of default group moderator")
	}

	if err := ToggleModerator(orga, user.ID, member.UserId); err != nil {
		t.Fatalf("Expected moderator to be toggled, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "remove_organization_moderator")

	if testIsGroupMember(t, orga, member.UserId, constants.GroupModName) {
		t.Fatal("User must not be member of default group moderator")
	}
}

func testDeleteAllUserGroupMember(t *testing.T) {
	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "user_group_member" WHERE 1 = 1`); err != nil {
		t.Fatal(err)
	}
}
