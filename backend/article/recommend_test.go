package article

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"strings"
	"testing"
)

func TestCheckUserReadAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	articleId := article.ID

	input := []struct {
		OrgaId    hide.ID
		UserId    hide.ID
		ArticleId hide.ID
	}{
		{0, 0, 0},
		{orga.ID, 0, 0},
		{orga.ID, user.ID, 0},
		{orga.ID, user.ID + 1, article.ID},
		{orga.ID, user.ID, article.ID},
	}
	expected := []error{
		errs.ArticleNotFound,
		errs.ArticleNotFound,
		errs.ArticleNotFound,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		var err error
		article, err = checkUserReadAccess(in.OrgaId, in.UserId, in.ArticleId)

		if err != expected[i] {
			t.Fatalf("Expected %v but was: %v", expected[i], err)
		}
	}

	if article == nil || article.ID != articleId {
		t.Fatal("Article must be returned")
	}
}

func TestJoinLatestArticleContent(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)

	if err := joinLatestArticleContent(orga.ID, user.ID, article); err != nil {
		t.Fatal("Latest content must have been joined")
	}
}

func TestCheckAllUserExist(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	_, err := checkAllUserExist(orga.ID, []hide.ID{user.ID, user.ID + 1})

	if err != errs.UserNotFound {
		t.Fatal("User must not have been found")
	}

	list, err := checkAllUserExist(orga.ID, []hide.ID{user.ID})

	if err != nil || len(list) != 1 {
		t.Fatal("User must have been found")
	}
}

func TestAppendGroupMembers(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	user1 := testutil.CreateUser(t, orga, 567, "user1@testutil.com")
	user2 := testutil.CreateUser(t, orga, 678, "user2@testutil.com")
	group1 := testutil.CreateUserGroup(t, orga, "group1")
	group2 := testutil.CreateUserGroup(t, orga, "group2")
	testutil.CreateUserGroupMember(t, group1, user1, false)
	testutil.CreateUserGroupMember(t, group1, user2, false)
	testutil.CreateUserGroupMember(t, group2, user2, false)

	user := make([]model.User, 0)
	user = appendGroupMembers(orga.ID, user, []hide.ID{group1.ID, group2.ID, group1.ID + 100})

	if len(user) != 2 {
		t.Fatalf("Two user must be in list, but was: %v", len(user))
	}
}

func TestCreateRecommendInviteArticleFeed(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	notify := []model.User{*user}

	createRecommendInviteArticleFeed(orga, user.ID, notify, article, "", "Hello World!", recommendArticleFeed)
	testutil.AssertFeedCreated(t, orga, recommendArticleFeed)
}

func TestSendRecommendInviteArticleMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	frontendHost = "http://domain.com/"
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	notify := []model.User{*user}

	// for existing article
	sendRecommendInviteArticleMail(orga, article, "", "", user, notify, mailtpl.RecommendMailTemplate, recommendMailSubject, recommendPath, recommendMailI18n, func(subject, body, from string, to ...string) error {
		t.Log(body)

		if subject != "You've got an article recommendation on Emvi" {
			t.Fatalf("Subject not as expected: %v", subject)
		}

		if !strings.Contains(body, user.Firstname) ||
			!strings.Contains(body, user.Lastname) ||
			!strings.Contains(body, "http://testutil.domain.com/read/") ||
			!strings.Contains(body, string(recommendMailI18n["en"]["title-1"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["title-2"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["text"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["action"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["link"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["greeting"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["goodbye"])) {
			t.Fatalf("Body not as expected: %v", body)
		}

		return nil
	})

	// for new article
	sendRecommendInviteArticleMail(orga, nil, "roomId123", "", user, notify, mailtpl.RecommendMailTemplate, recommendMailSubject, recommendPath, recommendMailI18n, func(subject, body, from string, to ...string) error {
		t.Log(body)

		if subject != "You've got an article recommendation on Emvi" {
			t.Fatalf("Subject not as expected: %v", subject)
		}

		if !strings.Contains(body, user.Firstname) ||
			!strings.Contains(body, user.Lastname) ||
			!strings.Contains(body, "http://testutil.domain.com/read/?room=roomId123") ||
			!strings.Contains(body, string(recommendMailI18n["en"]["title-1"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["title-2"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["text"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["action"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["link"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["greeting"])) ||
			!strings.Contains(body, string(recommendMailI18n["en"]["goodbye"])) {
			t.Fatalf("Body not as expected: %v", body)
		}

		return nil
	})
}

func TestRecommendArticleMailDisabled(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	user3 := testutil.CreateUser(t, orga, 322, "user3@test.com")
	user4 := testutil.CreateUser(t, orga, 323, "user4@test.com")
	user5 := testutil.CreateUser(t, orga, 324, "user5@test.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user4, false)
	testutil.CreateUserGroupMember(t, group, user5, false)
	user3.OrganizationMember.RecommendationMail = false
	user5.OrganizationMember.RecommendationMail = false

	if err := model.SaveOrganizationMember(nil, user3.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, user5.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	mailCount := 0
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		mailCount++
		return nil
	}

	if err := RecommendArticle(orga, user.ID, article.ID, []hide.ID{user2.ID, user3.ID}, []hide.ID{group.ID}, "", false, mailMock); err != nil {
		t.Fatalf("Article must be recommended, but was: %v", err)
	}

	if mailCount != 2 {
		t.Fatalf("Two mails must have been send, but was: %v", mailCount)
	}
}

func TestRecommendArticleReadConfirmation(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")

	mailCount := 0
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		mailCount++
		return nil
	}

	if err := RecommendArticle(orga, user.ID, article.ID, []hide.ID{user2.ID}, []hide.ID{}, "", true, mailMock); err != nil {
		t.Fatalf("Article must be recommended, but was: %v", err)
	}

	if mailCount != 1 {
		t.Fatalf("One mails must have been send, but was: %v", mailCount)
	}

	if model.GetArticleRecommendationByArticleIdAndUserIdAndRecommendedTo(article.ID, user.ID, user2.ID) == nil {
		t.Fatal("Article recommendation must exist")
	}

	if err := RecommendArticle(orga, user.ID, article.ID, []hide.ID{user2.ID}, []hide.ID{}, "", true, mailMock); err != nil {
		t.Fatalf("Article must be recommended, but was: %v", err)
	}

	if model.GetArticleRecommendationByArticleIdAndUserIdAndRecommendedTo(article.ID, user.ID, user2.ID) == nil {
		t.Fatal("Only one article recommendation must exist")
	}
}

func TestRecommendArticleLanguageNotAvailable(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")

	mailCount := 0
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		mailCount++
		return nil
	}

	if err := RecommendArticle(orga, user.ID, article.ID, []hide.ID{user2.ID}, []hide.ID{}, "", true, mailMock); err != nil {
		t.Fatalf("Article must be recommended, but was: %v", err)
	}

	if mailCount != 1 {
		t.Fatalf("One mails must have been send, but was: %v", mailCount)
	}

	if model.GetArticleRecommendationByArticleIdAndUserIdAndRecommendedTo(article.ID, user.ID, user2.ID) == nil {
		t.Fatal("Article recommendation must exist")
	}
}

func TestCheckMessageLen(t *testing.T) {
	input := []string{"", "this is fine", util.GenRandomString(maxMessageLen), util.GenRandomString(maxMessageLen + 1)}
	expected := []error{nil, nil, nil, errs.MessageTooLong}

	for i, in := range input {
		if err := checkMessageLen(in); err != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		}
	}
}
