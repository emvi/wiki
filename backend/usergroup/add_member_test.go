package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestAddUserGroupMember(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	emptyIds := []hide.ID{}

	if _, err := AddUserGroupMember(orga, user.ID, 0, emptyIds, emptyIds); err != errs.GroupNotFound {
		t.Fatal("Group must not be found")
	}

	group := testutil.CreateUserGroup(t, orga, "groupname")

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, emptyIds, emptyIds); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	testutil.CreateUserGroupMember(t, group, user, true)

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, []hide.ID{321}, emptyIds); err != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	newMember := testutil.CreateUser(t, orga, 2, "new@member.com")

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, []hide.ID{newMember.ID}, emptyIds); err != nil {
		t.Fatalf("New member must be added, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "add_user_group_member")
}

func TestAddUserGroupMemberFromGroupNotFound(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, "groupname")
	testutil.CreateUserGroupMember(t, group, user, true)
	emptyIds := []hide.ID{}

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, emptyIds, []hide.ID{567}); err != errs.GroupToAddNotFound {
		t.Fatalf("Group to add must not be found, but was: %v", err)
	}
}

func TestAddUserGroupMemberFromGroupSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	emptyIds := []hide.ID{}
	group := testutil.CreateUserGroup(t, orga, "groupname")
	member := &model.UserGroupMember{UserGroupId: group.ID, UserId: user.ID, IsModerator: true}

	if err := model.SaveUserGroupMember(nil, member); err != nil {
		t.Fatal(err)
	}

	hasMembers(t, group.ID, 1)
	groupToAdd := testutil.CreateUserGroup(t, orga, "othergroup")
	user2 := testutil.CreateUser(t, orga, 543, "user2@testutil.com")
	member1 := &model.UserGroupMember{UserGroupId: groupToAdd.ID, UserId: user.ID, IsModerator: true}
	member2 := &model.UserGroupMember{UserGroupId: groupToAdd.ID, UserId: user2.ID, IsModerator: false}

	if err := model.SaveUserGroupMember(nil, member1); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveUserGroupMember(nil, member2); err != nil {
		t.Fatal(err)
	}

	hasMembers(t, groupToAdd.ID, 2)

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, emptyIds, []hide.ID{groupToAdd.ID}); err != nil {
		t.Fatalf("Members from group must have been added, but was: %v", err)
	}

	hasMembers(t, group.ID, 2)
	testutil.AssertFeedCreated(t, orga, "add_user_group_member")
}

func TestAddUserGroupMemberMemberCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 543, "user2@testutil.com")
	user3 := testutil.CreateUser(t, orga, 544, "user3@testutil.com")
	user4 := testutil.CreateUser(t, orga, 545, "user4@testutil.com")
	group := testutil.CreateUserGroup(t, orga, "groupname")
	group2 := testutil.CreateUserGroup(t, orga, "group2")
	testutil.CreateUserGroupMember(t, group, user, true)
	testutil.CreateUserGroupMember(t, group2, user3, false)
	testutil.CreateUserGroupMember(t, group2, user4, false)

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, []hide.ID{user2.ID}, nil); err != nil {
		t.Fatalf("Member must have been added, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 2 {
		t.Fatalf("Group must have two members, but was: %v", group.MemberCount)
	}

	if _, err := AddUserGroupMember(orga, user.ID, group.ID, nil, []hide.ID{group2.ID}); err != nil {
		t.Fatalf("Member must have been added, but was: %v", err)
	}

	group = model.GetUserGroupByOrganizationIdAndId(orga.ID, group.ID)

	if group.MemberCount != 4 {
		t.Fatalf("Group must have four members, but was: %v", group.MemberCount)
	}
}

func hasMembers(t *testing.T, groupId hide.ID, n int) {
	is := model.CountUserGroupMemberByUserGroupId(groupId)

	if is != n {
		t.Fatalf("Group must have %v members, but was: %v", n, is)
	}
}
