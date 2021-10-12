package testutil

import (
	"emviwiki/shared/model"
	"encoding/json"
	"reflect"
	"testing"
)

func AssertFeedCreated(t *testing.T, orga *model.Organization, reason string) []model.Feed {
	feed := model.FindFeedByOrganizationIdAndReason(orga.ID, reason)

	if len(feed) != 1 {
		t.Fatalf("Expected one feed with reason '%v' to be created, but was %v", reason, len(feed))
	}

	return feed
}

func AssertFeedCreatedN(t *testing.T, orga *model.Organization, reason string, n int) []model.Feed {
	feed := model.FindFeedByOrganizationIdAndReason(orga.ID, reason)

	if len(feed) != n {
		t.Fatalf("Expected %v feed entries with reason '%v' to be created, but was %v", n, reason, len(feed))
	}

	return feed
}

func AssertJSONEquals(t *testing.T, s1, s2 string) {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(s1), &o1)

	if err != nil {
		t.Fatalf("Error mashalling string 1: %s", err)
	}

	err = json.Unmarshal([]byte(s2), &o2)

	if err != nil {
		t.Fatalf("Error mashalling string 2: %s", err)
	}

	if !reflect.DeepEqual(o1, o2) {
		t.Fatal("JSON not equal")
	}
}
