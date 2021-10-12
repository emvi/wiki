package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestRegistrationPassword(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true), RegistrationStep: StepPassword}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationPasswordData{"invalid", "Valid", "Valid"}

	if err := RegistrationPassword(data); err != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	data.Code = "code"

	if err := RegistrationPassword(data); err != nil {
		t.Fatalf("Password must have been set, but was: %v", err)
	}

	user = model.GetUserByEmailIgnoreActive("test@user.com")

	if !user.Password.Valid || !user.PasswordSalt.Valid {
		t.Fatal("User password and salt must be valid")
	}

	if user.RegistrationStep != StepPersonalData {
		t.Fatal("Step must have been incremented")
	}
}

func TestRegistrationPasswordCheckStep(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true)}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationPasswordData{"code", "Valid", "Valid"}

	if err := RegistrationPassword(data); err != errs.StepInvalid {
		t.Fatalf("Step must be invalid, but was: %v", err)
	}
}

func TestGetUserByRegistrationCodeAndCheckStep(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true), RegistrationStep: 3}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	if _, err := getUserByRegistrationCodeAndCheckStep("code", 5); err != errs.StepInvalid {
		t.Fatalf("Step must be invalid, but was: %v", err)
	}

	if _, err := getUserByRegistrationCodeAndCheckStep("code", 3); err != nil {
		t.Fatalf("Step must be valid, but was: %v", err)
	}

	if _, err := getUserByRegistrationCodeAndCheckStep("code", 1); err != nil {
		t.Fatalf("Step must be valid, but was: %v", err)
	}
}
