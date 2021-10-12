package user

import (
	"emviwiki/auth/errs"
	"emviwiki/shared/testutil"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")

	if err := ValidatePassword(0, ""); err != errs.UserNotFound {
		t.Fatal("User must not be found")
	}

	if err := ValidatePassword(user.ID, ""); err != errs.PasswordWrong {
		t.Fatal("Password must be wrong")
	}

	if err := ValidatePassword(user.ID, "password"); err != nil {
		t.Fatal("Password must be correct")
	}
}
