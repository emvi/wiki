package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"strings"
	"testing"
)

func TestChangeUserEmail(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")

	if err := ChangeUserEmail(0, EmailData{}, "en", mailMock); err != errs.UserNotFound {
		t.Fatal("User must not be found")
	}

	if err := ChangeUserEmail(user.ID, EmailData{}, "en", mailMock); err != errs.EmailEmpty {
		t.Fatal("Email must be empty")
	}

	if err := ChangeUserEmail(user.ID, EmailData{"invalid"}, "en", mailMock); err != errs.EmailInvalid {
		t.Fatal("Email must be invalid")
	}

	createTestUser(t, "test2@test.com")

	if err := ChangeUserEmail(user.ID, EmailData{"test2@test.com"}, "en", mailMock); err != errs.EmailInUse {
		t.Fatal("Email must be in use")
	}

	if err := ChangeUserEmail(user.ID, EmailData{"test@test.com"}, "en", mailMock); err != errs.EmailNotChanged {
		t.Fatal("Email must not be changed")
	}

	if err := ChangeUserEmail(user.ID, EmailData{"test@test.de"}, "en", mailMock); err != nil {
		t.Fatal("User must be saved")
	}

	user = model.GetUserById(user.ID)

	if user.NewEmail.String != "test@test.de" || !user.NewEmailCode.Valid || len(user.NewEmailCode.String) != 40 {
		t.Fatal("New email and code must be set")
	}
}

func TestChangeUserEmailSSOUser(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")
	user.AuthProvider = "not_emvi"

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	if err := ChangeUserEmail(user.ID, EmailData{"test2@test.de"}, "en", mailMock); err != errs.OperationNotAllowed {
		t.Fatalf("Operation must not be allowed, but was: %v", err)
	}
}

func TestChangeUserEmailMail(t *testing.T) {
	testutil.CleanAuthDb(t)
	authHost = "auth-host"
	user := createTestUser(t, "test@test.com")
	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if err := ChangeUserEmail(user.ID, EmailData{"new@test.de"}, "jp", mailMock); err != nil {
		t.Fatalf("User email must have been changed, but was: %v", err)
	}

	if subject != "Your email address at Emvi was changed" {
		t.Fatalf("Title not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, "new@test.de") ||
		!strings.Contains(body, authHost) ||
		!strings.Contains(body, confirmMailEndpoint) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["action"])) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["link"])) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(changeEmailMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func mailMock(subject, msgHTML, from string, to ...string) error {
	return nil
}
