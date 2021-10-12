package feed

import "testing"

func TestCheckReasonExists(t *testing.T) {
	if CheckReasonExists("delete_articlelistd") {
		t.Fatal("Feed reason must not exist")
	}

	if !CheckReasonExists("delete_articlelist") {
		t.Fatal("Feed reason must exist")
	}
}
