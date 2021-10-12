package search

import (
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSearchUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga := createTestUserWithMember(t)

	user, count := SearchUser(orga, "", nil)

	if len(user) != 2 || count != 2 {
		t.Fatalf("Two users must be found, but was %v, count was %v", len(user), count)
	}

	user, count = SearchUser(orga, "query", nil)

	if len(user) != 0 || count != 0 {
		t.Fatalf("No users must be found, but was %v", len(user))
	}

	user, count = SearchUser(orga, "firstname", nil)

	if len(user) != 2 || count != 2 {
		t.Fatalf("Two users must be found by firstname, but was %v", len(user))
	}

	user, count = SearchUser(orga, "lastname", nil)

	if len(user) != 2 || count != 2 {
		t.Fatalf("Two users must be found by lastname, but was %v", len(user))
	}

	user, count = SearchUser(orga, "username", nil)

	if len(user) != 2 || count != 2 {
		t.Fatalf("Two users must be found by username, but was %v", len(user))
	}

	if user[0].Email != "user1@test.com" ||
		user[0].Firstname != "firstname1" ||
		user[0].Lastname != "lastname1" ||
		user[0].OrganizationMember.ID == 0 ||
		user[0].OrganizationMember.Username != "username1" ||
		!user[0].OrganizationMember.IsModerator ||
		!user[0].OrganizationMember.IsAdmin {
		t.Fatal("First user not as expected")
	}

	if user[1].Email != "user2@test.com" ||
		user[1].Firstname != "firstname2" ||
		user[1].Lastname != "lastname2" ||
		user[1].OrganizationMember.ID == 0 ||
		user[1].OrganizationMember.Username != "username2" ||
		user[1].OrganizationMember.IsModerator ||
		user[1].OrganizationMember.IsAdmin {
		t.Fatal("Second user not as expected")
	}

	user, count = SearchUser(orga, "user@test.com", nil)

	if len(user) != 2 {
		t.Fatalf("Two users must be found by email, but was %v", len(user))
	}
}

func createTestUserWithMember(t *testing.T) *model.Organization {
	user1 := &model.User{BaseEntity: db.BaseEntity{ID: 1}, Email: "user1@test.com", Firstname: "firstname1", Lastname: "lastname1"}
	user2 := &model.User{BaseEntity: db.BaseEntity{ID: 2}, Email: "user2@test.com", Firstname: "firstname2", Lastname: "lastname2"}
	orga := &model.Organization{OwnerUserId: user1.ID, Name: "test orga", NameNormalized: "testorga", Expert: true}

	if err := model.SaveUser(nil, user1, true); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveUser(nil, user2, true); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	lang := &model.Language{OrganizationId: orga.ID, Name: "English", Code: "en", Default: true}

	if err := model.SaveLanguage(nil, lang); err != nil {
		t.Fatal(err)
	}

	member1 := &model.OrganizationMember{Username: "username1", IsModerator: true, IsAdmin: true, OrganizationId: orga.ID, UserId: user1.ID, LanguageId: lang.ID, Active: true}
	member2 := &model.OrganizationMember{Username: "username2", IsModerator: false, IsAdmin: false, OrganizationId: orga.ID, UserId: user2.ID, LanguageId: lang.ID, Active: true}

	if err := model.SaveOrganizationMember(nil, member1); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, member2); err != nil {
		t.Fatal(err)
	}

	return orga
}
