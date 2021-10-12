package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestHasAccessToOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	input := []struct {
		UserId hide.ID
		OrgaId hide.ID
	}{
		{0, 0},
		{user.ID + 1, orga.ID + 1},
		{user.ID + 1, orga.ID},
		{user.ID, orga.ID + 1},
		{user.ID, orga.ID},
	}
	expected := []error{
		errs.PermissionDenied,
		errs.PermissionDenied,
		errs.PermissionDenied,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		if err := HasAccessToOrganization(in.UserId, in.OrgaId); err != expected[i] {
			t.Fatalf("Expected %v for user id %v and organization id %v, but was: %v", expected[i], in.UserId, in.OrgaId, err)
		}
	}
}
