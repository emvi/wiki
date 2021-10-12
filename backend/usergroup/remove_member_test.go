package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestRemoveUserGroupMember(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	var emptyIds []hide.ID

	if err := RemoveUserGroupMember(orga, user.ID, 0, emptyIds); err != errs.GroupNotFound {
		t.Fatalf("Group must not be found, but was: %v", err)
	}

	group := testutil.CreateUserGroup(t, orga, "groupname")
	moderator := &model.UserGroupMember{UserGroupId: group.ID, UserId: user.ID, IsModerator: true}

	if err := model.SaveUserGroupMember(nil, moderator); err != nil {
		t.Fatal(err)
	}

	memberUser := testutil.CreateUser(t, orga, 321, "member@group.com")
	member := testutil.CreateUserGroupMember(t, group, memberUser, false)

	if err := RemoveUserGroupMember(orga, user.ID, group.ID, []hide.ID{moderator.ID}); err != errs.RemoveYourself {
		t.Fatalf("Must not be able to remove yourself: %v", err)
	}

	if err := RemoveUserGroupMember(orga, user.ID, group.ID, []hide.ID{member.ID}); err != nil {
		t.Fatalf("Member must have been removed, but was: %v", err)
	}

	member = model.GetUserGroupMemberById(member.ID)

	if member != nil {
		t.Fatal("Member must not exist")
	}

	testutil.AssertFeedCreated(t, orga, "remove_user_group_member")
}

func TestRemoveUserGroupMemberMemberCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 543, "user2@testutil.com")
	user3 := testutil.CreateUser(t, orga, 544, "user3@testutil.com")
	user4 := testutil.CreateUser(t, orga, 545, "user4@testutil.com")
	group := testutil.CreateUserGroup(t, orga, "groupname")
	testutil.CreateUserGroupMember(t, group, user, true)
	member1 := testutil.CreateUserGroupMember(t, group, user2, false)
	member2 := testutil.CreateUserGroupMember(t, group, user3, false)
	member3 := testutil.CreateUserGroupMember(t, group, user4, false)

	if err := RemoveUserGroupMember(orga, user.ID, group.ID, []hide.ID{member1.ID, member2.ID}); err != nil {
		t.Fatalf("Member must have been removed, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 2 {
		t.Fatalf("Group must have two members, but was: %v", group.MemberCount)
	}

	if err := RemoveUserGroupMember(orga, user.ID, group.ID, []hide.ID{member3.ID}); err != nil {
		t.Fatalf("Member must have been removed, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 1 {
		t.Fatalf("Group must have one members, but was: %v", group.MemberCount)
	}
}
