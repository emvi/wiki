package search

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestSearchUsergroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga := createTestUserWithMember(t)
	createTestGroups(t, orga)

	groups, count := SearchUsergroup(orga, "query", nil)

	if len(groups) != 0 || count != 0 {
		t.Fatal("No user groups must be found")
	}

	groups, count = SearchUsergroup(orga, "group", nil)

	if len(groups) != 2 || count != 2 {
		t.Fatalf("Two user groups must be found, but was: %v %v", len(groups), count)
	}

	if groups[0].Name != "group" ||
		groups[0].Info.String != "A user group" {
		t.Fatal("First user group not as expected")
	}

	if groups[1].Name != "test" ||
		groups[1].Info.String != "Another user group" {
		t.Fatal("Second user group not as expected")
	}

	groups, count = SearchUsergroup(orga, "test", nil)

	if len(groups) != 1 || count != 1 {
		t.Fatalf("One user group must be found, but was: %v %v", len(groups), count)
	}

	if groups[0].Name != "test" ||
		groups[0].Info.String != "Another user group" {
		t.Fatal("User group not as expected")
	}
}

func TestSearchUsergroupNoExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateUserGroup(t, orga, "test")

	if groups, n := SearchUsergroup(orga, "test", nil); len(groups) != 1 || n != 1 {
		t.Fatalf("One group must have been found, but was: %v %v", len(groups), n)
	}

	filter := model.SearchUserGroupFilter{FindGroups: false}

	if groups, n := SearchUsergroup(orga, "test", &filter); len(groups) != 0 || n != 0 {
		t.Fatalf("No groups must have been found, but was: %v %v", len(groups), n)
	}
}

func TestSearchUsergroupUserIds(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, "test")
	testutil.CreateUserGroup(t, orga, "test2")
	testutil.CreateUserGroupMember(t, group, user, false)
	filter := model.SearchUserGroupFilter{
		UserIds:    []hide.ID{user.ID},
		FindGroups: true,
	}
	groups, count := SearchUsergroup(orga, "", &filter)

	// including default groups!
	if len(groups) != 5 || count != 5 {
		t.Fatalf("One user groups must be found, but was: %v %v", len(groups), count)
	}
}

func createTestGroups(t *testing.T, orga *model.Organization) {
	group1 := &model.UserGroup{Name: "group", Info: null.NewString("A user group", true), OrganizationId: orga.ID}
	group2 := &model.UserGroup{Name: "test", Info: null.NewString("Another user group", true), OrganizationId: orga.ID}

	if err := model.SaveUserGroup(nil, group1); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveUserGroup(nil, group2); err != nil {
		t.Fatal(err)
	}
}
