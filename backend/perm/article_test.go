package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestFilterAccessList(t *testing.T) {
	testutil.CleanBackendDb(t)
	list := []SaveArticleAccess{{UserId: 1, Write: false},
		{UserId: 2, Write: false},
		{UserId: 1, Write: true},
		{UserGroupId: 3, Write: false},
		{UserGroupId: 3, Write: true},
		{UserGroupId: 4, Write: false}}
	newlist := FilterAccessList(list)

	if len(newlist) != 4 {
		t.Fatal("List must contain 4 items")
	}

	expected := []struct {
		UserId      hide.ID
		UserGroupId hide.ID
		Write       bool
	}{
		{1, 0, true},
		{2, 0, false},
		{0, 3, true},
		{0, 4, false},
	}

	for i, expect := range expected {
		if newlist[i].UserId != expect.UserId ||
			newlist[i].UserGroupId != expect.UserGroupId ||
			newlist[i].Write != expect.Write {
			t.Fatalf("Item %v not as expected, expected %v, but was %v", i, expect, newlist[i])
		}
	}
}

func TestCheckAccessList(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	list := []SaveArticleAccess{{UserId: 1}}

	if err := CheckAccessList(nil, orga, list); err != errs.UserOrUserGroupNotFound {
		t.Fatal("User must not be found")
	}

	list = []SaveArticleAccess{{UserGroupId: 1}}

	if err := CheckAccessList(nil, orga, list); err != errs.UserOrUserGroupNotFound {
		t.Fatal("User group must not be found")
	}

	list = []SaveArticleAccess{{UserId: user.ID}}

	if err := CheckAccessList(nil, orga, list); err != nil {
		t.Fatal("User must be found")
	}

	group := testutil.CreateUserGroup(t, orga, "usergroup")
	list = []SaveArticleAccess{{UserGroupId: group.ID}}

	if err := CheckAccessList(nil, orga, list); err != nil {
		t.Fatal("User group must be found")
	}
}

func TestGetUserIds(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user, true)
	access := []SaveArticleAccess{{UserId: 9999}, {UserGroupId: group.ID}}
	ids := GetUserIdsFromArticleAccess(access)

	if len(ids) != 2 {
		t.Fatal("IDs must contain two user ids")
	}

	if ids[0] != 9999 || ids[1] != user.ID {
		t.Fatal("User ids not ads expected")
	}
}
