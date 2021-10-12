package notification

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"strings"
	"sync"
	"testing"
	"time"
)

type testMailSend struct {
	subject string
	body    string
	to      string
}

func TestSendNotificationForMember(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateLang(t, orga, "de", "Deutsch", false)
	setNotificationInterval(t, user.OrganizationMember, 3)
	createTestFeed(t, orga, user, lang, true, 2, "create_article_list")

	var mailsSend []testMailSend
	var m sync.Mutex
	mailProvider = func(subject, msgHTML, from string, to ...string) error {
		m.Lock()
		defer m.Unlock()
		// from is the receiver in this case
		mailsSend = append(mailsSend, testMailSend{subject, msgHTML, from})
		return nil
	}

	if err := sendNotificationForMember(user.OrganizationMember); err != nil {
		t.Fatalf("Must send notification mail for user, but was: %v", err)
	}

	if len(mailsSend) != 1 {
		t.Fatalf("Must have send 1 mail, but was: %v", len(mailsSend))
	}

	if mailsSend[0].subject != "Your unread notifications on Emvi" {
		t.Fatalf("Unexpected subject, was: %v", mailsSend[0].subject)
	}

	if mailsSend[0].to != user.Email {
		t.Fatalf("Mail must have been send to user, but was send to: %v", mailsSend[0].to)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user.ID)

	if member.NextNotificationMail.Before(time.Now().Add(time.Hour*24*3 - time.Minute)) {
		t.Fatalf("Next notification must have been set, but was: %v", member.NextNotificationMail)
	}
}

func TestSendNotificationMails(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga1, user1 := testutil.CreateOrgaAndUser(t)
	orga2, member2 := testutil.CreateOrga(t, user1, "onemoreorga")
	lang := testutil.CreateLang(t, orga1, "en", "English", true)
	testutil.CreateLang(t, orga1, "de", "Deutsch", false)
	testutil.CreateLang(t, orga2, "ru", "Russian", true)
	testutil.CreateLang(t, orga2, "de", "Deutsch", false)
	user2 := testutil.CreateUser(t, orga1, 321, "user2@test.com")

	setNotificationInterval(t, user1.OrganizationMember, 3)
	setNotificationInterval(t, member2, 4)
	setNotificationInterval(t, user2.OrganizationMember, 2)

	// these two in one mail
	createTestFeed(t, orga1, user1, lang, true, 1, "create_article_list")
	createTestFeed(t, orga1, user1, lang, true, 2, "create_article_list")

	// another mail for different organization
	createTestFeed(t, orga2, user1, lang, true, 2, "create_article_list")

	// irrelevant notifications for user1
	createTestFeed(t, orga1, user1, lang, true, 4, "create_article_list")
	createTestFeed(t, orga1, user1, lang, false, 2, "create_article_list")
	createTestFeed(t, orga2, user1, lang, true, 3, "create_article_list")
	createTestFeed(t, orga2, user1, lang, false, 1, "create_article_list")

	// mail for user2
	createTestFeed(t, orga1, user2, lang, true, 1, "create_article_list") // send this one only
	createTestFeed(t, orga1, user2, lang, true, 2, "create_article_list")
	createTestFeed(t, orga1, user2, lang, true, 3, "create_article_list")
	createTestFeed(t, orga1, user2, lang, false, 1, "create_article_list")

	var mailsSend []testMailSend
	var m sync.Mutex
	mailProvider = func(subject, msgHTML, from string, to ...string) error {
		m.Lock()
		defer m.Unlock()
		// from is the receiver in this case
		mailsSend = append(mailsSend, testMailSend{subject, msgHTML, from})
		return nil
	}

	SendNotificationMails()

	if len(mailsSend) != 3 {
		t.Fatalf("Must have send 3 mails, but was: %v", len(mailsSend))
	}

	member1 := model.GetOrganizationMemberByOrganizationIdAndUserId(orga1.ID, user1.ID)

	if member1.NextNotificationMail.Before(time.Now().Add(time.Hour*24*3 - time.Minute)) {
		t.Fatalf("Next notification must have been set, but was: %v", member1.NextNotificationMail)
	}

	member2 = model.GetOrganizationMemberByOrganizationIdAndUserId(orga1.ID, user2.ID)

	if member1.NextNotificationMail.Before(time.Now().Add(time.Hour*24*2 - time.Minute)) {
		t.Fatalf("Next notification must have been set, but was: %v", member1.NextNotificationMail)
	}

	member1Orga2 := model.GetOrganizationMemberByOrganizationIdAndUserId(orga2.ID, user1.ID)

	if member1Orga2.NextNotificationMail.Before(time.Now().Add(time.Hour*24*4 - time.Minute)) {
		t.Fatalf("Next notification must have been set, but was: %v", member1.NextNotificationMail)
	}
}

func TestSendNotificationForMemberTwoNotifications(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateLang(t, orga, "de", "Deutsch", false)
	setNotificationInterval(t, user.OrganizationMember, 3)
	feedId1 := createTestFeed(t, orga, user, lang, true, 2, "create_article_list")
	feedId2 := createTestFeed(t, orga, user, lang, true, 2, "create_article_list")

	var mailsSend []testMailSend
	var m sync.Mutex
	mailProvider = func(subject, msgHTML, from string, to ...string) error {
		m.Lock()
		defer m.Unlock()
		// from is the receiver in this case
		mailsSend = append(mailsSend, testMailSend{subject, msgHTML, from})
		return nil
	}

	if err := sendNotificationForMember(user.OrganizationMember); err != nil {
		t.Fatalf("Must send notification mail for user, but was: %v", err)
	}

	if len(mailsSend) != 1 {
		t.Fatalf("Must have send 1 mail, but was: %v", len(mailsSend))
	}

	id, _ := hide.ToString(getFeedArticleListId(orga.ID, lang.ID, feedId1))
	if !strings.Contains(mailsSend[0].body, id) {
		t.Fatalf("Mail must contain first feed with ID %v, but was: %v", id, mailsSend[0].body)
	}

	id, _ = hide.ToString(getFeedArticleListId(orga.ID, lang.ID, feedId2))
	if !strings.Contains(mailsSend[0].body, id) {
		t.Fatalf("Mail must contain second feed with ID %v, but was: %v", id, mailsSend[0].body)
	}
}

func TestSendNotificationMailsIgnoreTooOld(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	setNotificationInterval(t, user.OrganizationMember, 3)
	feedId1 := createTestFeed(t, orga, user, lang, true, 2, "create_article_list")
	feedId2 := createTestFeed(t, orga, user, lang, true, 32, "create_article_list")

	var mailsSend []testMailSend
	var m sync.Mutex
	mailProvider = func(subject, msgHTML, from string, to ...string) error {
		m.Lock()
		defer m.Unlock()
		// from is the receiver in this case
		mailsSend = append(mailsSend, testMailSend{subject, msgHTML, from})
		return nil
	}

	if err := sendNotificationForMember(user.OrganizationMember); err != nil {
		t.Fatalf("Must send notification mail for user, but was: %v", err)
	}

	if len(mailsSend) != 1 {
		t.Fatalf("Must have send 1 mail, but was: %v", len(mailsSend))
	}

	id, _ := hide.ToString(getFeedArticleListId(orga.ID, lang.ID, feedId1))
	if !strings.Contains(mailsSend[0].body, id) {
		t.Fatalf("Mail must contain first feed with ID %v, but was: %v", id, mailsSend[0].body)
	}

	id, _ = hide.ToString(getFeedArticleListId(orga.ID, lang.ID, feedId2))
	if strings.Contains(mailsSend[0].body, id) {
		t.Fatalf("Mail must not contain second feed with ID %v, but was: %v", id, mailsSend[0].body)
	}
}

func TestRenderAndSendMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateLang(t, orga, "de", "Deutsch", false)
	setNotificationInterval(t, user.OrganizationMember, 3)
	createTestFeed(t, orga, user, lang, true, 2, "create_article_list")
	createTestFeed(t, orga, user, lang, true, 2, "joined_organization")

	var mailsSend []testMailSend
	var m sync.Mutex
	mailProvider = func(subject, msgHTML, from string, to ...string) error {
		m.Lock()
		defer m.Unlock()
		// from is the receiver in this case
		mailsSend = append(mailsSend, testMailSend{subject, msgHTML, from})
		return nil
	}

	if err := sendNotificationForMember(user.OrganizationMember); err != nil {
		t.Fatalf("Must send notification mail for user, but was: %v", err)
	}

	if len(mailsSend) != 1 {
		t.Fatalf("Must have send 1 mail, but was: %v", len(mailsSend))
	}

	if len(mailsSend[0].subject) == 0 {
		t.Fatal("Title must be set")
	}

	body := mailsSend[0].body
	t.Log(body)

	if !strings.Contains(body, string(notificationMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(notificationMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(notificationMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(notificationMailI18n["en"]["text-3"])) ||
		!strings.Contains(body, string(notificationMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(notificationMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, "3 days") ||
		!strings.Contains(body, "created a new list") ||
		!strings.Contains(body, "joined the organization") {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func setNotificationInterval(t *testing.T, member *model.OrganizationMember, days uint) {
	member.SendNotificationsInterval = days

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}
}

func createTestFeed(t *testing.T, orga *model.Organization, user *model.User, lang *model.Language, notification bool, ageDays int, reason string) hide.ID {
	feed := testutil.CreateFeed(t, orga, user, lang, notification)
	defTime := time.Now().Add(time.Hour * 24 * time.Duration(-ageDays))

	if _, err := model.GetConnection().Exec(nil, `UPDATE "feed" SET def_time = $1, reason = $2 WHERE id = $3`, defTime, reason, feed.ID); err != nil {
		t.Fatal(err)
	}

	return feed.ID
}

func getFeedArticleListId(orgaId, langId, feedId hide.ID) hide.ID {
	refs := model.FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orgaId, langId, feedId)

	for _, ref := range refs {
		if ref.ArticleListID != 0 {
			return ref.ArticleListID
		}
	}

	return 0
}
