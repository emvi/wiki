package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestLeaveOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")

	input := []struct {
		userId hide.ID
		name   string
	}{
		{user.ID, "wrong"},
		{user2.ID, "wrong"},
		{user2.ID, orga.Name},
	}
	expected := []error{
		errs.OrganizationOwnerLeave,
		errs.NameDoesNotMatch,
		nil,
	}

	for i, in := range input {
		if err := LeaveOrganization(orga, in.userId, in.name); err != expected[i] {
			t.Fatalf("Expectedd '%v', but was: %v", expected[i], err)
		}
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user2.ID)

	if member != nil {
		t.Fatal("User must have left organization")
	}

	testutil.AssertFeedCreated(t, orga, "left_organization")
}

func TestLeaveOrganizationUpdateGroupMemberCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user2, false)

	if err := LeaveOrganization(orga, user2.ID, orga.Name); err != nil {
		t.Fatalf("User must have left organization, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 0 {
		t.Fatalf("User group must not have members anymore")
	}
}
