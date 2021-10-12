package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadUserGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, "group")

	result, _, _, err := ReadUserGroup(orga, user.ID, 0)

	if err != errs.GroupNotFound {
		t.Fatal("Group must not be found")
	}

	result, _, _, err = ReadUserGroup(orga, user.ID, group.ID)

	if err != nil {
		t.Fatal("Group must be found")
	}

	if result.Name != "group" {
		t.Fatalf("Group must have propper name, but was: %v", result.Name)
	}
}

func TestReadUserGroupMember(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester2@test.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user, false)
	testutil.CreateUserGroupMember(t, group, user2, true)

	_, _, err := ReadUserGroupMember(orga, group.ID+1, nil)

	if err != errs.GroupNotFound {
		t.Fatal("Group must not be found")
	}

	member, count, err := ReadUserGroupMember(orga, group.ID, nil)

	if err != nil {
		t.Fatalf("Group member must be found, but was: %v", err)
	}

	if len(member) != 2 || count != 2 {
		t.Fatalf("Group must have 2 members, but was: %v %v", len(member), count)
	}

	if member[0].User.ID != user.ID {
		t.Fatalf("First user must be as expected, but was: %v", member[0].User.ID)
	}

	if member[1].User.ID != user2.ID {
		t.Fatalf("Second user must be as expected, but was: %v", member[0].User.ID)
	}
}
