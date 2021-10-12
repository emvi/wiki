package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestToggleUserGroupModerator(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if err := ToggleUserGroupModerator(orga, user.ID, 0, 321); err != errs.MemberNotFound {
		t.Fatal("Member must not be found")
	}

	group := testutil.CreateUserGroup(t, orga, "groupname")
	moderator := &model.UserGroupMember{UserGroupId: group.ID, UserId: user.ID, IsModerator: true}

	if err := model.SaveUserGroupMember(nil, moderator); err != nil {
		t.Fatal(err)
	}

	if err := ToggleUserGroupModerator(orga, user.ID, group.ID, moderator.ID); err != errs.ModeratorYourself {
		t.Fatal("Must not toggle moderator for yourself")
	}

	memberUser := testutil.CreateUser(t, orga, 321, "member@group.com")
	member := testutil.CreateUserGroupMember(t, group, memberUser, false)

	if err := ToggleUserGroupModerator(orga, user.ID, group.ID, member.ID); err != nil {
		t.Fatal("Moderator must have been toggled")
	}

	member = model.GetUserGroupMemberById(member.ID)

	if member == nil || !member.IsModerator {
		t.Fatal("Moderator must have been toggled")
	}

	testutil.AssertFeedCreated(t, orga, "set_user_group_moderator")

	if err := ToggleUserGroupModerator(orga, user.ID, group.ID, member.ID); err != nil {
		t.Fatal("Moderator must have been toggled")
	}

	testutil.AssertFeedCreated(t, orga, "remove_user_group_moderator")
}
