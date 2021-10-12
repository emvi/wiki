package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestChangeUserData(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")

	input := []struct {
		UserId    hide.ID
		Language  string
		Marketing bool
	}{
		{0, "", false},
		{user.ID, "a", false},
	}
	expected := []error{
		errs.UserNotFound,
		errs.LanguageInvalid,
	}

	for i, in := range input {
		if err := ChangeUserData(in.UserId, UserData{user.Firstname.String, user.Lastname.String, in.Language, in.Marketing}); err[0] != expected[i] {
			t.Fatalf("Expected %v but was %v", expected[i], err[0])
		}
	}

	if err := ChangeUserData(user.ID, UserData{user.Firstname.String, user.Lastname.String, "", true}); err != nil {
		t.Fatal("User must have been saved")
	}

	user = model.GetUserById(user.ID)

	if user.Language.Valid {
		t.Fatal("User attributes must not be set")
	}

	if err := ChangeUserData(user.ID, UserData{user.Firstname.String, user.Lastname.String, "en", true}); err != nil {
		t.Fatalf("User must be saved, but was: %v", err)
	}

	user = model.GetUserById(user.ID)

	if !user.Language.Valid || !user.AcceptMarketing {
		t.Fatal("User attributes must be set")
	}
}

func TestChangeUserDataName(t *testing.T) {
	testutil.CleanAuthDb(t)
	user := createTestUser(t, "test@test.com")

	input := []struct {
		UserId    hide.ID
		Firstname string
		Lastname  string
	}{
		{0, "", ""},
		{user.ID, "", ""},
		{user.ID, "Valid", ""},
		{user.ID, "", "Valid"},
		{user.ID, "Владимир", "Путин"},
	}
	expected := [][]error{
		{errs.UserNotFound},
		{errs.FirstnameInvalid, errs.LastnameInvalid},
		{errs.LastnameInvalid},
		{errs.FirstnameInvalid},
		nil,
	}

	for i, in := range input {
		if err := ChangeUserData(in.UserId, UserData{in.Firstname, in.Lastname, user.Language.String, user.AcceptMarketing}); len(err) != len(expected[i]) {
			for j, e := range err {
				if e != expected[i][j] {
					t.Fatalf("Expected %v but was %v", expected[i][j], e)
				}
			}
		}
	}

	user = model.GetUserById(user.ID)

	if user.Firstname.String != "Владимир" || user.Lastname.String != "Путин" {
		t.Fatal("User firstname and lastname must have been changed")
	}
}

func createTestUser(t *testing.T, email string) *model.User {
	user := &model.User{Email: email,
		Firstname:    null.NewString("Paul", true),
		Lastname:     null.NewString("Johnson", true),
		Active:       true,
		Password:     null.NewString(util.Sha256Base64("password"), true),
		AuthProvider: "emvi"}

	if err := model.SaveUser(nil, user); err != nil {
		t.Fatal(err)
	}

	return user
}
