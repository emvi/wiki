package article

import (
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"strings"
	"testing"
)

// see recommend_test.go

func TestSendInviteArticleMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	to := testutil.CreateUser(t, orga, 321, "to@user.com")
	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}
	errs := sendRecommendInviteArticleMail(orga, article, "room-id", "message", user, []model.User{*to}, mailtpl.InviteArticleMailTemplate, inviteMailSubject, invitePath, inviteMailI18n, mailMock)

	if len(errs) != 0 {
		t.Fatalf("Mail must have been send, but was: %v", errs)
	}

	if subject != "You've got an invitation to edit an article on Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)
	idHash, _ := hide.ToString(article.ID)

	if !strings.Contains(body, article.LatestArticleContent.Title) ||
		!strings.Contains(body, user.Firstname) ||
		!strings.Contains(body, user.Lastname) ||
		!strings.Contains(body, idHash) ||
		!strings.Contains(body, "message") ||
		!strings.Contains(body, string(inviteMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["title-1"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["title-2"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["action"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["link"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestSendInviteArticleMailNewArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	to := testutil.CreateUser(t, orga, 321, "to@user.com")
	var subject, body string
	mailMock := func(mailSubject string, mailBody string, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}
	errs := sendRecommendInviteArticleMail(orga, nil, "room-id", "message", user, []model.User{*to}, mailtpl.InviteArticleMailTemplate, inviteMailSubject, invitePath, inviteMailI18n, mailMock)

	if len(errs) != 0 {
		t.Fatalf("Mail must have been send, but was: %v", errs)
	}

	if subject != "You've got an invitation to edit an article on Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, user.Firstname) ||
		!strings.Contains(body, user.Lastname) ||
		!strings.Contains(body, "room-id") ||
		!strings.Contains(body, string(inviteMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["title-1"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["title-2"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["action"])) ||
		!strings.Contains(body, string(inviteMailI18n["en"]["link"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
