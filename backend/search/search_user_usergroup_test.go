package search

import (
	"emviwiki/shared/testutil"
	"testing"
)

func TestSearchUserUsergroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga := createTestUserWithMember(t)
	createTestGroups(t, orga)

	results := SearchUserUsergroup(orga, "")

	if len(results) != 4 {
		t.Fatalf("Four results must be found, but was: %v", len(results))
	}

	results = SearchUserUsergroup(orga, "test")

	if len(results) != 3 {
		t.Fatal("Three results must be found")
	}

	if results[0].User == nil || results[1].User == nil || results[2].UserGroup == nil {
		t.Fatal("Results must contain: user, user, group")
	}

	results = SearchUserUsergroup(orga, "group")

	if len(results) != 2 {
		t.Fatal("Two results (groups) must be found")
	}

	if results[0].UserGroup == nil || results[1].UserGroup == nil {
		t.Fatal("Results must contain: group, group")
	}

	results = SearchUserUsergroup(orga, "name")

	if len(results) != 2 {
		t.Fatalf("Two results (user) must be found, but was: %v", len(results))
	}

	if results[0].User == nil || results[1].User == nil {
		t.Fatal("Results must contain: user, user")
	}
}
