package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/recaptcha"
	"emviwiki/shared/testutil"
	"testing"
)

type recaptchaValidatorMock struct{}

func (validator *recaptchaValidatorMock) Validate(token string) (*recaptcha.RecaptchaResponse, error) {
	return &recaptcha.RecaptchaResponse{Success: true}, nil
}

func TestRegistrationData(t *testing.T) {
	testutil.CleanAuthDb(t)
	blacklist := &model.EmailBlacklist{Domain: "allowed.com"}

	if err := model.SaveEmailBlacklist(nil, blacklist); err != nil {
		t.Fatal(err)
	}

	input := []string{
		"",
		"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
		"foobar",
		"not@allowed.com",
		"thisis@fine.com",
	}
	expected := []error{
		errs.EmailEmpty,
		errs.EmailTooLong,
		errs.EmailInvalid,
		errs.EmailForbidden,
		nil,
	}

	for i, in := range input {
		data := RegistrationData{in}

		if err := data.validate(); err != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		}
	}
}

func TestRegistrationPasswordData(t *testing.T) {
	testutil.CleanAuthDb(t)

	input := []struct {
		pwd1 string
		pwd2 string
	}{
		{"", ""},
		{"not the", "same"},
		{"okay", "okay"},
	}
	expected := []error{
		errs.PasswordInvalid,
		errs.PasswordMatch,
		nil,
	}

	for i, in := range input {
		data := RegistrationPasswordData{"", in.pwd1, in.pwd2}

		if err := data.validate(); err != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		}
	}
}

func TestRegistrationPersonalData(t *testing.T) {
	testutil.CleanAuthDb(t)

	input := []struct {
		firstname string
		lastname  string
		language  string
	}{
		{"", "valid", ""},
		{"valid", "", ""},
		{"valid", "valid", "asdf"},
		{"valid", "valid", "en"},
	}
	expected := []error{
		errs.FirstnameInvalid,
		errs.LastnameInvalid,
		errs.LanguageInvalid,
		nil,
	}

	for i, in := range input {
		data := RegistrationPersonalData{"", in.firstname, in.lastname, in.language}
		err := data.validate()

		if expected[i] != nil && err[0] != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		} else if expected[i] == nil && err != nil {
			t.Fatalf("Expected no error, but was: %v", err)
		}
	}
}

func TestRegistrationCompletionData(t *testing.T) {
	testutil.CleanAuthDb(t)
	recaptchaValidator = new(recaptchaValidatorMock)

	input := []struct {
		acceptTermsOfService bool
		acceptPrivacy        bool
	}{
		{false, true},
		{true, false},
		{true, true},
	}
	expected := []error{
		errs.TermsOfService,
		errs.Privacy,
		nil,
	}

	for i, in := range input {
		data := RegistrationCompletionData{"", in.acceptTermsOfService, in.acceptPrivacy, false, "valid"}
		err := data.validate()

		if expected[i] != nil && err[0] != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		} else if expected[i] == nil && err != nil {
			t.Fatalf("Expected no error, but was: %v", err)
		}
	}
}
