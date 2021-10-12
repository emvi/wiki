package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestGetRegistrationStep(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "user@test.com", RegistrationStep: StepInitial, RegistrationCode: null.NewString("code", true)}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	if _, err := ConfirmRegistration("invalid"); err != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	if step, err := ConfirmRegistration("code"); err != nil || step != StepPassword {
		t.Fatalf("Step must be returned, but was: %v %v", err, step)
	}
}
