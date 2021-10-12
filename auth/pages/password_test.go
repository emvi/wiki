package pages

import (
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendPasswordResetMail(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createUser(t)
	r := httptest.NewRequest("POST", "/password", nil)
	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if !sendPasswordResetMail(r, user, "new-password", mailMock) {
		t.Fatal("Password mail must have been send")
	}

	if subject != "Your password at Emvi has been reset" {
		t.Fatalf("Password not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, "new-password") ||
		!strings.Contains(body, string(passwordResetMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(passwordResetMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(passwordResetMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(passwordResetMailI18n["en"]["text"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func createUser(t *testing.T) *model.User {
	user := &model.User{Email: "test@user.com"}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	return user
}
