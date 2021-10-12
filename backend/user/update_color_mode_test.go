package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestUpdateColorMode(t *testing.T) {
	testutil.CleanBackendDb(t)
	_, user := testutil.CreateOrgaAndUser(t)

	if err := UpdateColorMode(0, -1); err != errs.ColorModeInvalid {
		t.Fatalf("Color mode must be invalid, but was: %v", err)
	}

	if err := UpdateColorMode(0, 3); err != errs.ColorModeInvalid {
		t.Fatalf("Color mode must be invalid, but was: %v", err)
	}

	if err := UpdateColorMode(0, 2); err != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	if err := UpdateColorMode(user.ID, 2); err != nil {
		t.Fatalf("User must have been updated, but was: %v", err)
	}

	user = model.GetUserById(user.ID)

	if user.ColorMode != 2 {
		t.Fatalf("Color mode must have been updated, but was: %v", user.ColorMode)
	}
}
