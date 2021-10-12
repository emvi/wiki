package feed

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"testing"
)

func TestCreateFeedNoAccessNoNotifications(t *testing.T) {
	testutil.CleanBackendDb(t)
	data := &CreateFeedData{}

	if err := CreateFeed(data); err != errs.NonPublicFeedWithoutAccess {
		t.Fatalf("No access and notifications must be provided, but was %v", err)
	}
}

func TestCreateFeedUnknownReason(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	data := &CreateFeedData{Organization: orga,
		UserId: user.ID,
		Reason: "unknown",
		Public: true}

	if err := CreateFeed(data); err != errs.ReasonNotFound {
		t.Fatalf("Reason must not be found, but was %v", err)
	}
}

func TestCreateFeedSuccess(t *testing.T) {
	orga := testCreateFeedSuccess(t, false)
	testFeedAccess(t, orga, 555, 0, false)
	testFeedAccess(t, orga, 123, 1, false)
	testFeedAccess(t, orga, 9998, 1, false)
	testFeedAccess(t, orga, 9999, 1, true)
	testFeedAccess(t, orga, 123, 0, true)
	testFeedAccess(t, orga, 9998, 0, true)
}

func TestCreateFeedSuccessPublic(t *testing.T) {
	orga := testCreateFeedSuccess(t, true)
	testFeedAccess(t, orga, 555, 1, false) // feed is public, so must have access
	testFeedAccess(t, orga, 123, 1, false)
	testFeedAccess(t, orga, 9998, 1, false)
	testFeedAccess(t, orga, 9999, 1, false)
	testFeedAccess(t, orga, 9999, 1, true)
	testFeedAccess(t, orga, 123, 0, true)
	testFeedAccess(t, orga, 9998, 0, true)
}

func TestCreateFeedReferences(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 999, "second@user.com")

	refs := make([]interface{}, 0)
	refs = append(refs, user)
	refs = append(refs, user2)
	data := &CreateFeedData{Organization: orga,
		UserId: user.ID,
		Reason: "joined_organization",
		Public: true,
		Refs:   refs}

	if err := CreateFeed(data); err != nil {
		t.Fatalf("Feed must have been created, but was %v", err)
	}

	feed := model.FindFeedByOrganizationIdAndReason(orga.ID, "joined_organization")

	if len(feed) != 1 {
		t.Fatalf("One feed must have been created, but was: %v", len(feed))
	}

	createdRefs := model.FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orga.ID, util.DetermineLang(nil, orga.ID, user.ID, 0).ID, feed[0].ID)

	if len(createdRefs) != 2 {
		t.Fatalf("Two feed refs must have been created, but was: %v", len(createdRefs))
	}

	foundUser1, foundUser2 := false, false

	for _, ref := range createdRefs {
		if ref.UserID == user.ID {
			foundUser1 = true
		} else if ref.UserID == user2.ID {
			foundUser2 = true
		}
	}

	if !foundUser1 || !foundUser2 {
		t.Fatal("Both users must have been references in feed ref")
	}
}

func TestCreateFeedNotificationNotForYourself(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@testutil.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	data := &CreateFeedData{Organization: orga,
		UserId: user.ID,
		Reason: "joined_organization",
		Public: false,
		Access: []hide.ID{user.ID, user2.ID, user.ID, user2.ID},
		Notify: []hide.ID{user.ID, user2.ID, user.ID, user2.ID},
		Refs:   []interface{}{article}}

	if err := CreateFeed(data); err != nil {
		t.Fatalf("Feed must be created, but was %v", err)
	}

	feed := model.FindFeedByOrganizationIdAndReason(orga.ID, "joined_organization")

	if len(feed) != 1 {
		t.Fatalf("One feed must have been created, but was: %v", len(feed))
	}

	access := model.FindFeedAccessByFeedId(feed[0].ID)

	if len(access) != 2 {
		t.Fatalf("Two feed access must have been created, but was: %v", len(feed))
	}

	foundUser := false
	foundUser2 := false

	for _, a := range access {
		if a.Notification && a.UserId == user2.ID && !a.Read {
			foundUser2 = true
		} else if !a.Notification && a.UserId == user.ID && a.Read {
			foundUser = true
		}
	}

	if !foundUser || !foundUser2 {
		t.Fatalf("Propper access for users not found, was: %v", access)
	}
}

func TestCreateFeedDontNotifyIsRead(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test2@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	data := &CreateFeedData{Organization: orga,
		UserId: user.ID,
		Reason: "joined_organization",
		Public: false,
		Access: []hide.ID{user2.ID},
		Refs:   []interface{}{article}}

	if err := CreateFeed(data); err != nil {
		t.Fatalf("Feed must be created, but was %v", err)
	}

	feed, _ := GetFilteredFeed(orga, user2.ID, nil)

	if len(feed) != 1 {
		t.Fatalf("One feed must have been created, but was: %v", len(feed))
	}

	access := model.GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotification(orga.ID, 321, feed[0].ID, false)

	if access == nil || !access.Read {
		t.Fatalf("Feed must be read, but was: %v", access)
	}
}

func TestCreateFeedKeyValue(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	kv1 := KeyValue{"foo", "bar"}
	kv2 := KeyValue{"hello", "world"}
	data := &CreateFeedData{Organization: orga,
		UserId: user.ID,
		Reason: "joined_organization",
		Public: true,
		Refs:   []interface{}{kv1, kv2}}

	if err := CreateFeed(data); err != nil {
		t.Fatalf("Feed must be created, but was %v", err)
	}

	feed, _ := GetFilteredFeed(orga, user.ID, nil)

	if len(feed) != 1 {
		t.Fatalf("One feed must have been created, but was: %v", len(feed))
	}

	refs := model.FindFeedRefByFeedId(feed[0].ID)

	if len(refs) != 2 {
		t.Fatalf("Two feed references must exist, but was: %v", len(refs))
	}

	if refs[0].Key.String != "foo" && refs[0].Key.String != "hello" ||
		refs[1].Key.String != "foo" && refs[1].Key.String != "hello" ||
		refs[0].Value.String != "bar" && refs[0].Value.String != "world" ||
		refs[1].Value.String != "bar" && refs[1].Value.String != "world" {
		t.Fatal("Key/values not as expected")
	}
}

func testCreateFeedSuccess(t *testing.T, public bool) *model.Organization {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 9998, "user2@testutil.com")
	user3 := testutil.CreateUser(t, orga, 9999, "user3@testutil.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	data := &CreateFeedData{Organization: orga,
		UserId: user.ID,
		Reason: "joined_organization",
		Public: public,
		Access: []hide.ID{user2.ID},
		Notify: []hide.ID{user3.ID},
		Refs:   []interface{}{article}}

	if err := CreateFeed(data); err != nil {
		t.Fatalf("Feed must be created, but was %v", err)
	}

	return orga
}

func testFeedAccess(t *testing.T, orga *model.Organization, userId hide.ID, n int, notification bool) []model.Feed {
	var filter *model.SearchFeedFilter

	if notification {
		filter = &model.SearchFeedFilter{Notifications: true}
	}

	feed, _ := GetFilteredFeed(orga, userId, filter)

	if len(feed) != n {
		t.Fatalf("Expected user '%v' to have access to '%v' feed entries, but was '%v'", userId, n, len(feed))
	}

	return feed
}
