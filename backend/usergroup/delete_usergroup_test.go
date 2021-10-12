package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDeleteUserGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if err := DeleteUserGroup(orga, user.ID, 1); err != errs.GroupNotFound {
		t.Fatalf("User group must not be found, but was: %v", err)
	}

	group := testutil.CreateUserGroup(t, orga, "groupname")

	if err := DeleteUserGroup(orga, user.ID, group.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	member := &model.UserGroupMember{UserGroupId: group.ID, UserId: user.ID, IsModerator: true}

	if err := model.SaveUserGroupMember(nil, member); err != nil {
		t.Fatal(err)
	}

	if err := DeleteUserGroup(orga, user.ID, group.ID); err != nil {
		t.Fatalf("User group must have been deleted, but was %v", err)
	}

	if model.GetUserGroupByOrganizationIdAndName(orga.ID, "groupname") != nil {
		t.Fatal("User group must not exist anymore")
	}

	feed := testutil.AssertFeedCreated(t, orga, "delete_usergroup")
	refs := model.FindFeedRefByFeedId(feed[0].ID)

	if len(refs) != 1 || refs[0].Key.String != "name" || refs[0].Value.String != "groupname" {
		t.Fatalf("Feed must have deleted object name, but was: %v", refs)
	}
}

func TestDeleteUserGroupReferencedObjects(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, "groupname")
	testutil.CreateUserGroupMember(t, group, user, true)
	testutil.CreateObservedObject(t, user, nil, nil, group)
	testutil.CreateFeedForObject(t, orga, user, nil, nil, group)

	if err := DeleteUserGroup(orga, user.ID, group.ID); err != nil {
		t.Fatalf("User group must have been deleted including all references, but was: %v", err)
	}
}
