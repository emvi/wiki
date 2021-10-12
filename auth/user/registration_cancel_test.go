package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestCancelRegistration(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true)}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	if err := CancelRegistration("invalid"); err != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	if err := CancelRegistration("code"); err != nil {
		t.Fatalf("User must have been deleted, but was: %v", err)
	}

	if model.GetUserByEmailIgnoreActive("test@user.com") != nil {
		t.Fatal("User must not exist anymore")
	}
}
