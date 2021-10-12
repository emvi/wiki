package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestUpdateIntroduction(t *testing.T) {
	testutil.CleanBackendDb(t)
	_, user := testutil.CreateOrgaAndUser(t)

	if err := UpdateIntroduction(0, false); err != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	if err := UpdateIntroduction(user.ID, false); err != nil {
		t.Fatalf("Introduction must have been updated, but was: %v", err)
	}

	user = model.GetUserById(user.ID)

	if user.Introduction {
		t.Fatal("Introduction must be false")
	}

	if err := UpdateIntroduction(user.ID, true); err != nil {
		t.Fatalf("Introduction must have been updated, but was: %v", err)
	}

	user = model.GetUserById(user.ID)

	if !user.Introduction {
		t.Fatal("Introduction must be true")
	}
}
