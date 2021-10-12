package user

import (
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
	"time"
)

func TestCreateOrUpdateUser(t *testing.T) {
	testutil.CleanBackendDb(t)

	if model.GetUserById(99) != nil {
		t.Fatal("User must not exist yet")
	}

	data := AuthUser{auth.UserResponse{99,
		"email",
		"firstname",
		"lastname",
		"en",
		"",
		"info",
		true,
		time.Now(),
		time.Now(),
		false}}
	user, err := CreateOrUpdateUser(data)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil ||
		user.Email != "email" ||
		user.Firstname != "firstname" ||
		user.Lastname != "lastname" ||
		user.Language.String != "en" ||
		user.Info.String != "info" ||
		!user.AcceptMarketing {
		t.Fatalf("User must have been created: %v", user)
	}

	data.Email = "new email"
	user, err = CreateOrUpdateUser(data)

	if err != nil {
		t.Fatal(err)
	}

	if user == nil ||
		user.Email != "new email" ||
		user.Firstname != "firstname" ||
		user.Lastname != "lastname" ||
		user.Language.String != "en" ||
		user.Info.String != "info" ||
		!user.AcceptMarketing ||
		user.Picture.Valid {
		t.Fatalf("User must have been updated: %v", user)
	}
}
