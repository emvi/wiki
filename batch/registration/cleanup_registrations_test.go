package registration

import (
	"emviwiki/auth/model"
	testutil2 "emviwiki/shared/testutil"
	"testing"
	"time"
)

func TestCleanupRegistrations(t *testing.T) {
	testutil2.CleanAuthDb(t)
	user1 := createTestUser(t, "test1@user.com", true, time.Now(), 0)
	user2 := createTestUser(t, "test2@user.com", false, time.Now(), 0)
	user3 := createTestUser(t, "test3@user.com", true, time.Now().Add(-time.Hour*24*33), 0)
	user4 := createTestUser(t, "test4@user.com", false, time.Now().Add(-time.Hour*24*33), 4)
	user5 := createTestUser(t, "test5@user.com", false, time.Now().Add(-time.Hour*24*33), 3)
	CleanupRegistrations()

	if model.GetUserByEmailIgnoreActive(user1.Email) == nil {
		t.Fatal("User 1 must exist")
	}

	if model.GetUserByEmailIgnoreActive(user2.Email) == nil {
		t.Fatal("User 2 must exist")
	}

	if model.GetUserByEmailIgnoreActive(user3.Email) == nil {
		t.Fatal("User 3 must exist")
	}

	if model.GetUserByEmailIgnoreActive(user4.Email) == nil {
		t.Fatal("User 4 must exist")
	}

	if model.GetUserByEmailIgnoreActive(user5.Email) != nil {
		t.Fatal("User 5 must not exist")
	}
}

func createTestUser(t *testing.T, email string, active bool, defTime time.Time, step int) *model.User {
	u := &model.User{Email: email, Active: active, RegistrationStep: step, AuthProvider: "emvi"}

	if err := model.SaveUser(nil, u); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `UPDATE "user" SET def_time = $1 WHERE id = $2`, defTime, u.ID); err != nil {
		t.Fatal(err)
	}

	return u
}
