package feed

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestToggleNotificationReadSingle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateUser(t, orga, 321, "member@user.com")
	feed := testutil.CreateFeed(t, orga, user, lang, true)
	testutil.CreateFeed(t, orga, user, lang, true)
	access := model.FindFeedAccessByFeedId(feed.ID)

	if len(access) != 1 || access[0].ID == 0 {
		t.Fatal("Feed access must be found")
	}

	if err := ToggleNotificationRead(nil, orga, user.ID, feed.ID+2); err != errs.FeedAccessNotFound {
		t.Fatalf("Feed access must not be found, but was %v", err)
	}

	if err := ToggleNotificationRead(nil, orga, user.ID, feed.ID); err != nil {
		t.Fatalf("Feed access must have been marked as read, but was %v", err)
	}

	assertNotificationsExists(t, orga, user.ID, 1)

	if err := ToggleNotificationRead(nil, orga, user.ID, feed.ID); err != nil {
		t.Fatalf("Feed access must have been marked as read, but was %v", err)
	}

	assertNotificationsExists(t, orga, user.ID, 2)
}

func TestToggleNotificationReadAll(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateUser(t, orga, 321, "member@user.com")
	testutil.CreateFeed(t, orga, user, lang, true)
	testutil.CreateFeed(t, orga, user, lang, true)

	if err := ToggleNotificationRead(nil, orga, user.ID, 0); err != nil {
		t.Fatalf("All feed access must have been marked as read, but was %v", err)
	}

	assertNotificationsExists(t, orga, user.ID, 0)
}

func assertNotificationsExists(t *testing.T, orga *model.Organization, userId hide.ID, n int) {
	notifications, _ := GetFilteredFeed(orga, userId, &model.SearchFeedFilter{Notifications: true})
	notificationsCount := 0

	for _, notification := range notifications {
		if !notification.Read {
			notificationsCount++
		}
	}

	if notificationsCount != n {
		t.Fatalf("Expected %v notifications to be found, but was %v", n, len(notifications))
	}
}
