package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestSaveUserGroupNewUserGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if _, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{0, "", ""}); len(err) != 1 || err[0] != errs.NameTooShort {
		t.Fatal("Name must not be set")
	}

	if _, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{0, "12345678901234567890123456789012345678901", "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"}); len(err) != 2 || err[0] != errs.NameTooLong || err[1] != errs.InfoTooLong {
		t.Fatal("Name and Info must be too long")
	}

	if id, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{0, "groupname", "groupinfo"}); id == 0 || err != nil {
		t.Fatal("New user group must be created")
	}

	group := model.GetUserGroupByOrganizationIdAndName(orga.ID, "GROUPname")

	if group == nil {
		t.Fatal("Group must have been created")
	}

	if group.Name != "groupname" || group.Info.String != "groupinfo" || group.MemberCount != 1 {
		t.Fatal("Group fields must be as expected")
	}

	testutil.AssertFeedCreated(t, orga, "create_user_group")
}

func TestSaveUserGroupNameInUse(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateUserGroup(t, orga, "groupname")

	if _, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{0, "groupname", "groupinfo"}); len(err) != 1 || err[0] != errs.NameInUse {
		t.Fatal("Name must be in use already")
	}
}

func TestSaveUserGroupUpdateUserGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	group := &model.UserGroup{OrganizationId: orga.ID, Name: "groupname", Info: null.NewString("groupinfo", true)}

	if err := model.SaveUserGroup(nil, group); err != nil {
		t.Fatal(err)
	}

	if _, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{group.ID, "name", "info"}); len(err) != 1 || err[0] != errs.PermissionDenied {
		t.Fatal("User must not have access to group")
	}

	member := &model.UserGroupMember{UserGroupId: group.ID, UserId: user.ID}

	if err := model.SaveUserGroupMember(nil, member); err != nil {
		t.Fatal(err)
	}

	if _, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{group.ID, "name", "info"}); len(err) != 1 || err[0] != errs.PermissionDenied {
		t.Fatal("User must not have access to group (moderator)")
	}

	member.IsModerator = true

	if err := model.SaveUserGroupMember(nil, member); err != nil {
		t.Fatal(err)
	}

	if _, err := SaveUserGroup(orga, user.ID, SaveUserGroupData{group.ID, "name", "info"}); err != nil {
		t.Fatal("Group must have been updated")
	}

	group = model.GetUserGroupByOrganizationIdAndName(orga.ID, "name")

	if group == nil {
		t.Fatal("Group must have been updated")
	}

	if group.Name != "name" || group.Info.String != "info" {
		t.Fatal("Group fields must be as expected")
	}

	testutil.AssertFeedCreated(t, orga, "update_user_group")
}

func TestCheckUserAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	mod := testutil.CreateUser(t, orga, 222, "mod@test.com")
	mod.OrganizationMember.IsModerator = true

	if err := model.SaveOrganizationMember(nil, mod.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	user := testutil.CreateUser(t, orga, 333, "user@test.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, admin, true)
	testutil.CreateUserGroupMember(t, group, mod, true)
	testutil.CreateUserGroupMember(t, group, user, true)
	orga.CreateGroupAdmin = true
	orga.CreateGroupMod = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if checkUserAccess(orga, admin.ID, group.ID) != nil {
		t.Fatal("Admin must have access to group")
	}

	if checkUserAccess(orga, mod.ID, group.ID) != nil {
		t.Fatal("Mod must have access to group")
	}

	if checkUserAccess(orga, user.ID, group.ID) != errs.PermissionDenied {
		t.Fatal("User must not have access to group")
	}

	orga.CreateGroupMod = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if checkUserAccess(orga, admin.ID, group.ID) != nil {
		t.Fatal("Admin must have access to group")
	}

	if checkUserAccess(orga, mod.ID, group.ID) != errs.PermissionDenied {
		t.Fatal("Mod must not have access to group")
	}

	if checkUserAccess(orga, user.ID, group.ID) != errs.PermissionDenied {
		t.Fatal("User must not have access to group")
	}

	orga.CreateGroupAdmin = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if checkUserAccess(orga, admin.ID, group.ID) != nil {
		t.Fatal("Admin must have access to group")
	}

	if checkUserAccess(orga, mod.ID, group.ID) != nil {
		t.Fatal("Mod must have access to group")
	}

	if checkUserAccess(orga, user.ID, group.ID) != nil {
		t.Fatal("User must have access to group")
	}
}
