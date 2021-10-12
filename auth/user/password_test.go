package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"strings"
	"testing"
)

func TestChangePassword(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")

	input := []struct {
		UserId  hide.ID
		OldPwd  string
		NewPwd1 string
		NewPwd2 string
	}{
		{0, "", "", ""},
		{user.ID, "test", "foo", "bar"},
		{user.ID, "password", "foo", "bar"},
		{user.ID, "password", "newpwd", "newpwd"},
	}
	expected := []error{
		errs.UserNotFound,
		errs.PasswordWrong,
		errs.PasswordMatch,
		nil,
	}

	for i, in := range input {
		if err := ChangePassword(in.UserId, PasswordData{in.OldPwd, in.NewPwd1, in.NewPwd2}, "en", mailMock); err != expected[i] {
			t.Fatalf("Expected % but was: %v", expected[i], err)
		}
	}
}

func TestChangePasswordSSOUser(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")
	user.AuthProvider = "not_emvi"

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	if err := ChangePassword(user.ID, PasswordData{"123", "321", "321"}, "en", mailMock); err != errs.OperationNotAllowed {
		t.Fatalf("Operation must not be allowed, but was: %v", err)
	}
}

func TestSendChangedPasswordMail(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")
	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}
	sendChangedPasswordMail(user, "jp", mailMock)

	if subject != "Your password at Emvi was changed" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, string(changePasswordMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(changePasswordMailI18n["en"]["text"])) ||
		!strings.Contains(body, string(changePasswordMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(changePasswordMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
