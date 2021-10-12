package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"testing"
)

func TestGetMember(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if _, err := GetMember(orga, user.ID+1); err != errs.MemberNotFound {
		t.Fatalf("Member must not be found, but was: %v", err)
	}

	if member, err := GetMember(orga, user.ID); err != nil || member == nil {
		t.Fatalf("Member must not found, but was: %v", err)
	}
}
