package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestToggleReadOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateUser(t, orga, 321, "member@testutil.com")
	testDeleteAllUserGroupMember(t)
	member := model.GetOrganizationMemberByUsername("testuser2")
	member.IsAdmin = true
	member.IsModerator = true

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
		{orga, member.UserId, member.UserId},
		{orga, user.ID, member.UserId},
	}
	expected := []struct {
		Error error
	}{
		{errs.PermissionDenied},
		{errs.MemberNotFound},
		{errs.MemberNotFound},
		{errs.OrganizationOwner},
		{errs.ChangeReadOnlyYourself},
		{nil},
	}

	for i := range input {
		if err := ToggleReadOnly(input[i].Orga, input[i].UserId, input[i].MemberUserId); err != expected[i].Error {
			t.Fatalf("Expected error '%v', but was: %v", expected[i].Error, err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "set_member_read_only")

	if !testIsGroupMember(t, orga, member.UserId, constants.GroupReadOnlyName) {
		t.Fatal("User must be member of default group read only")
	}

	if err := ToggleReadOnly(orga, user.ID, member.UserId); err != nil {
		t.Fatalf("Expected read only to be toggled, but was: %v", err)
	}

	if testIsGroupMember(t, orga, member.UserId, constants.GroupReadOnlyName) {
		t.Fatal("User must not be member of default group read only")
	}

	testutil.AssertFeedCreated(t, orga, "remove_member_read_only")
}
