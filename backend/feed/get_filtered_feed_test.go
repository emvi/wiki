package feed

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestGetFilteredFeed(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateUser(t, orga, 321, "member@user.com")
	testutil.CreateFeed(t, orga, user, lang, false)
	testutil.CreateFeed(t, orga, user, lang, false)

	feed, count := GetFilteredFeed(orga, user.ID, nil)

	if len(feed) != 2 || count != 0 {
		t.Fatalf("Two feed entries must be found, but was: %v %v", len(feed), count)
	}

	for _, f := range feed {
		if f.Feed != "joined the organization." {
			t.Fatalf("Feed entry not as expected, was: '%v' '%v'", f.Feed, f.Notification)
		}

		testFeedHasRef(t, &f, user.ID, false, false)
		testFeedHasRef(t, &f, 0, true, false)
		testFeedHasRef(t, &f, 0, false, true)
	}
}

func TestGetFilteredFeedNotifications(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateUser(t, orga, 321, "member@user.com")
	testutil.CreateFeed(t, orga, user, lang, true)
	testutil.CreateFeed(t, orga, user, lang, true)
	testutil.CreateFeed(t, orga, user, lang, false)

	feed, count := GetFilteredFeed(orga, user.ID, &model.SearchFeedFilter{Notifications: true})

	if len(feed) != 2 || count != 2 {
		t.Fatalf("Two notifications must be found and two must be unread, but was: %v %v", len(feed), count)
	}

	for _, f := range feed {
		if f.Feed != "joined the organization." {
			t.Fatalf("Notification not as expected, was %v", f)
		}

		testFeedHasRef(t, &f, user.ID, false, false)
		testFeedHasRef(t, &f, 0, true, false)
		testFeedHasRef(t, &f, 0, false, true)
	}
}

func TestGetFilteredFeedFilterByUserId(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	testFeedPublic(t, testutil.CreateFeed(t, orga, user, lang, false))
	testFeedPublic(t, testutil.CreateFeed(t, orga, user2, lang, false))
	testFeedPublic(t, testutil.CreateFeed(t, orga, user2, lang, false))

	filter := &model.SearchFeedFilter{UserIds: []hide.ID{user.ID}}
	feed, count := GetFilteredFeed(orga, user.ID, filter)

	if len(feed) != 1 || count != 0 {
		t.Fatalf("Expected one entries to be found, but was: %v %v", len(feed), count)
	}

	filter = &model.SearchFeedFilter{UserIds: []hide.ID{user2.ID}}
	feed, count = GetFilteredFeed(orga, user.ID, filter)

	if len(feed) != 2 || count != 0 {
		t.Fatalf("Expected two entries to be found, but was: %v %v", len(feed), count)
	}

	filter = &model.SearchFeedFilter{UserIds: []hide.ID{user.ID, user2.ID}}
	feed, count = GetFilteredFeed(orga, user.ID, filter)

	if len(feed) != 3 || count != 0 {
		t.Fatalf("Expected three entries to be found, but was: %v %v", len(feed), count)
	}
}

func TestGetFilteredFeedUnread(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateFeed(t, orga, user, lang, true)
	testutil.CreateFeed(t, orga, user, lang, true)
	testFeedRead(t, testutil.CreateFeed(t, orga, user, lang, true), true)

	filter := &model.SearchFeedFilter{Notifications: true}
	feed, count := GetFilteredFeed(orga, user.ID, filter)

	if len(feed) != 3 && count != 2 {
		t.Fatalf("Three notifications and two unread notifications must be returned, but was: %v %v", len(feed), count)
	}

	filter.Unread = true
	feed, count = GetFilteredFeed(orga, user.ID, filter)

	if len(feed) != 2 && count != 2 {
		t.Fatalf("Two unread notifications must be returned, but was: %v %v", len(feed), count)
	}
}

func TestGetSupportedLanguage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := model.GetLanguageByOrganizationIdAndCode(orga.ID, "ja")
	lang.Default = true

	if err := model.SaveLanguage(nil, lang); err != nil {
		t.Fatal(err)
	}

	if getSupportedLanguage(orga.ID, user.ID) != "en" {
		t.Fatal("Must return user selected language")
	}

	user.Language.SetValid("ru")

	if err := model.SaveUser(nil, user, false); err != nil {
		t.Fatal(err)
	}

	if getSupportedLanguage(orga.ID, user.ID) != "en" {
		t.Fatal("Must return default language")
	}

	lang.Code = "de"

	if err := model.SaveLanguage(nil, lang); err != nil {
		t.Fatal(err)
	}

	if getSupportedLanguage(orga.ID, user.ID) != "de" {
		t.Fatal("Must return organization default language")
	}
}

func testFeedPublic(t *testing.T, feed *model.Feed) {
	feed.Public = true

	if err := model.SaveFeed(nil, feed); err != nil {
		t.Fatal(err)
	}
}

func testFeedRead(t *testing.T, feed *model.Feed, read bool) {
	access := model.FindFeedAccessByFeedId(feed.ID)

	for _, a := range access {
		a.Read = read

		if err := model.SaveFeedAccess(nil, &a); err != nil {
			t.Fatal(err)
		}
	}
}

func testFeedHasRef(t *testing.T, feed *model.Feed, userId hide.ID, groupId, articleId bool) {
	found := false

	if userId != 0 {
		for _, ref := range feed.FeedRefs {
			if ref.UserID == userId {
				found = true
				break
			}
		}
	} else if groupId {
		for _, ref := range feed.FeedRefs {
			if ref.UserGroupID != 0 {
				found = true
				break
			}
		}
	} else if articleId {
		for _, ref := range feed.FeedRefs {
			if ref.ArticleID != 0 {
				found = true
				break
			}
		}
	}

	if !found {
		t.Fatalf("Referenced object not found for feed: %v %v %v", userId, groupId, articleId)
	}
}
