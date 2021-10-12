package pages

import (
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginPageHandlerGET(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	LoginPageHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Login page must return status OK, but was %v", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "<form method=\"post\">") {
		t.Fatalf("Login page must be returned, but was %v", rec.Body.String())
	}
}

func TestLoginPageHandlerPOSTInvalidInput(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	LoginPageHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Login page must return status OK, but was %v", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "<p class=\"login-error\">") {
		t.Fatalf("Login page must be returned, but was %v", rec.Body.String())
	}
}

func TestSaveUserLogin(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := &model.User{
		Email:  "test@user.com",
		Active: true,
	}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	saveUserLogin(user.ID)
	login := model.FindLoginByUserId(user.ID)

	if len(login) != 1 {
		t.Fatalf("One login entity must have been created")
	}
}
