package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/config"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"strings"
	"testing"
)

func TestCompleteRegistration(t *testing.T) {
	testutil.CleanAuthDb(t)
	recaptchaValidator = new(recaptchaValidatorMock)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true), RegistrationStep: StepCompletion}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationCompletionData{"", true, true, true, ""}

	if _, _, err := CompleteRegistration(data, "", mailMock); err[0] != errs.UserNotFound {
		t.Fatalf("User must not be found, but was: %v", err)
	}

	data.Code = "code"
	token, expires, err := CompleteRegistration(data, "", mailMock)

	if err != nil {
		t.Fatalf("Registration must have been completed, but was: %v", err)
	}

	if token == "" {
		t.Fatalf("Token must be returned, but was: %v", token)
	}

	if expires.IsZero() {
		t.Fatalf("Expires in must have been returned, but was: %v", expires)
	}

	user = model.GetUserByEmailIgnoreActive("test@user.com")

	if !user.AcceptMarketing || user.RegistrationCode.Valid || !user.Active {
		t.Fatal("Data must have been saved")
	}

	if user.RegistrationStep != StepDone {
		t.Fatal("Step must have been incremented")
	}
}

func TestCompleteRegistrationCheckStep(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true)}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	data := RegistrationCompletionData{"code", true, true, true, ""}

	if _, _, err := CompleteRegistration(data, "", mailMock); err[0] != errs.StepInvalid {
		t.Fatalf("Step must be invalid, but was: %v", err)
	}
}

func TestSendRegistrationCompletedMail(t *testing.T) {
	testutil.CleanAuthDb(t)
	config.Get().Registration.CompletedNewOrgaURI = "new-orga"
	config.Get().Registration.CompletedJoinOrgaURI = "join-orga"
	LoadConfig()
	user := &model.User{Email: "test@user.com", RegistrationCode: null.NewString("code", true), RegistrationStep: StepCompletion}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}
	sendRegistrationCompletedMail(user.Email, "en", mailMock)

	if subject != "Thank you for signing up at Emvi!" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, registrationCompletedNewOrgaURI) ||
		!strings.Contains(body, registrationCompletedJoinOrgaURI) ||
		!strings.Contains(body, string(registrationCompleteMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(registrationCompleteMailI18n["en"]["text"])) ||
		!strings.Contains(body, string(registrationCompleteMailI18n["en"]["or"])) ||
		!strings.Contains(body, string(registrationCompleteMailI18n["en"]["action_create"])) ||
		!strings.Contains(body, string(registrationCompleteMailI18n["en"]["action_join"])) ||
		!strings.Contains(body, string(registrationCompleteMailI18n["en"]["link"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
