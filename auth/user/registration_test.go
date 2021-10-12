package user

import (
	"emviwiki/auth/model"
	"emviwiki/shared/config"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"strings"
	"testing"
)

func TestRegistrationNewUser(t *testing.T) {
	testutil.CleanAuthDb(t)
	data := RegistrationData{"test@user.com"}

	if err := Registration(data, "en", mailMock); err != nil {
		t.Fatalf("Registration must be successful, but was: %v", err)
	}

	user := model.GetUserByEmailIgnoreActive("test@user.com")

	if user == nil {
		t.Fatal("New user must exist")
	}

	if user.Language.String != "en" {
		t.Fatalf("User language must be en, but was: %v", user.Language)
	}

	if user.RegistrationStep != StepPassword {
		t.Fatalf("User registration step must be 1, but was: %v", user.RegistrationStep)
	}
}

func TestRegistrationExistingUser(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com",
		Language:         null.NewString("de", true),
		RegistrationCode: null.NewString("code", true),
		RegistrationStep: 42,
		AuthProvider:     emviAuthProviderName}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationData{"test@user.com"}

	if err := Registration(data, "en", mailMock); err != nil {
		t.Fatalf("Registration must be successful, but was: %v", err)
	}

	user = model.GetUserByEmailIgnoreActive("test@user.com")

	if user == nil {
		t.Fatal("New user must exist")
	}

	if user.Language.String != "de" {
		t.Fatalf("User language must be de, but was: %v", user.Language)
	}

	if user.RegistrationMailsSend != 1 {
		t.Fatalf("One additional registration mail must have been send, but was: %v", user.RegistrationMailsSend)
	}

	if user.RegistrationStep != 42 {
		t.Fatalf("User registration step must be 42, but was: %v", user.RegistrationStep)
	}
}

func TestSendRegistrationConfirmationMail(t *testing.T) {
	testutil.CleanAuthDb(t)
	config.Get().Registration.ConfirmationURI = "registration-url"
	LoadConfig()
	data := RegistrationData{"test@user.com"}
	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if err := Registration(data, "en", mailMock); err != nil {
		t.Fatalf("Registration must be successful, but was: %v", err)
	}

	if subject != "Your registration at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, "?code=") ||
		!strings.Contains(body, string(registrationMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(registrationMailI18n["en"]["text"])) ||
		!strings.Contains(body, string(registrationMailI18n["en"]["action"])) ||
		!strings.Contains(body, string(registrationMailI18n["en"]["link"])) ||
		!strings.Contains(body, string(registrationMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(registrationMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, registrationConfirmationURI) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
