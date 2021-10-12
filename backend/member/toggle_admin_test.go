package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestToggleAdmin(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@testutil.com")
	adminMember := model.GetOrganizationMemberByUsername("testuser1")
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
		{orga, user2.ID, 0},
		{orga, user2.ID, adminMember.UserId},
		{orga, user2.ID, user2.ID},
		{orga, user.ID, member.UserId},
	}
	expected := []struct {
		Error error
	}{
		{errs.PermissionDenied},
		{errs.MemberNotFound},
		{errs.OrganizationOwner},
		{errs.ChangeAdminYourself},
		{nil},
	}

	for i := range input {
		if err := ToggleAdmin(input[i].Orga, input[i].UserId, input[i].MemberUserId); err != expected[i].Error {
			t.Fatalf("Expected error '%v', but was: %v", expected[i].Error, err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "remove_organization_admin")

	if testIsGroupMember(t, orga, member.UserId, constants.GroupAdminName) {
		t.Fatal("User must not be member of default group admin")
	}

	if err := ToggleAdmin(orga, user.ID, member.UserId); err != nil {
		t.Fatalf("Expected admin to be toggled, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "set_organization_admin")

	if !testIsGroupMember(t, orga, member.UserId, constants.GroupAdminName) {
		t.Fatal("User must be member of default group admin")
	}
}

func testIsGroupMember(t *testing.T, orga *model.Organization, memberUserId hide.ID, name string) bool {
	group := model.GetUserGroupByOrganizationIdAndName(orga.ID, name)

	if group == nil {
		t.Fatal("Group not found")
	}

	return model.GetUserGroupMemberByGroupIdAndUserId(group.ID, memberUserId) != nil
}
