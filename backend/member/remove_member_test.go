package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestRemoveMember(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateUser(t, orga, 321, "member@testutil.com")
	adminMember := model.GetOrganizationMemberByUsername("testuser1")
	member := model.GetOrganizationMemberByUsername("testuser2")

	input := []struct {
		UserId       hide.ID
		MemberUserId hide.ID
	}{
		{0, 0},
		{user.ID, 0},
		{user.ID, adminMember.UserId},
		{user.ID, member.UserId},
	}
	expected := []struct {
		Error error
	}{
		{errs.PermissionDenied},
		{errs.MemberNotFound},
		{errs.RemoveYourself},
		{nil},
	}

	for i := range input {
		if err := RemoveMember(orga, input[i].UserId, input[i].MemberUserId, false); err != expected[i].Error {
			t.Fatalf("Expected error '%v', but was: %v", expected[i].Error, err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "left_organization")
}

func TestRemoveMemberUpdateGroupMemberCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user2, true)

	if err := RemoveMember(orga, user.ID, user2.ID, false); err != nil {
		t.Fatalf("User must have been removed, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 1 {
		t.Fatalf("User group must have one member, but was: %v", group.MemberCount)
	}
}
