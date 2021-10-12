package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestRegistrationPersonal(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true), RegistrationStep: StepPersonalData}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationPersonalData{"invalid", "Max", "Mustermann", "de"}

	if err := RegistrationPersonal(data); err[0] != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	data.Code = "code"

	if err := RegistrationPersonal(data); err != nil {
		t.Fatalf("Personal data must have been set, but was: %v", err)
	}

	user = model.GetUserByEmailIgnoreActive("test@user.com")

	if user.Firstname.String != "Max" || user.Lastname.String != "Mustermann" || user.Language.String != "de" {
		t.Fatal("Personal data must have been saved")
	}

	if user.RegistrationStep != StepCompletion {
		t.Fatal("Step must have been incremented")
	}
}

func TestRegistrationPersonalCheckStep(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true), AuthProvider: "emvi"}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationPersonalData{"code", "Max", "Mustermann", ""}

	if err := RegistrationPersonal(data); err[0] != errs.StepInvalid {
		t.Fatalf("Step must be invalid, but was: %v", err)
	}
}
