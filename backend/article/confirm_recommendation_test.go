package article

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestConfirmRecommendation(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	recommendTo := testutil.CreateUser(t, orga, 321, "recommend@to.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	mailMock := func(subject, msgHTML, from string, to ...string) error {
		return nil
	}

	if err := RecommendArticle(orga, user.ID, article.ID, []hide.ID{recommendTo.ID}, []hide.ID{}, "Hello World!", true, mailMock); err != nil {
		t.Fatalf("Article must have been recommended, but was: %v", err)
	}

	if err := ConfirmRecommendation(orga, recommendTo.ID, article.ID, []hide.ID{user.ID}); err != nil {
		t.Fatalf("Recommendation must have been confirmed, but was: %v", err)
	}

	if len(model.FindArticleRecommendationByArticleIdAndRecommendedTo(article.ID, recommendTo.ID)) != 0 {
		t.Fatalf("Article recommendation must not exist anymore")
	}

	testutil.AssertFeedCreated(t, orga, "recommendation_confirmation")
}

func TestConfirmRecommendationNoNotification(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 222, "me@test.com")
	recommendTo := testutil.CreateUser(t, orga, 321, "recommend@to.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	mailMock := func(subject, msgHTML, from string, to ...string) error {
		return nil
	}

	// both want to receive a confirmation
	if err := RecommendArticle(orga, user.ID, article.ID, []hide.ID{recommendTo.ID}, []hide.ID{}, "", true, mailMock); err != nil {
		t.Fatalf("Article must have been recommended, but was: %v", err)
	}

	if err := RecommendArticle(orga, user2.ID, article.ID, []hide.ID{recommendTo.ID}, []hide.ID{}, "", true, mailMock); err != nil {
		t.Fatalf("Article must have been recommended, but was: %v", err)
	}

	// but the user doesn't want to send it to both
	if err := ConfirmRecommendation(orga, recommendTo.ID, article.ID, []hide.ID{user2.ID}); err != nil {
		t.Fatalf("Recommendation must have been confirmed, but was: %v", err)
	}

	if len(model.FindArticleRecommendationByArticleIdAndRecommendedTo(article.ID, recommendTo.ID)) != 0 {
		t.Fatalf("Article recommendation must not exist anymore")
	}

	feed := testutil.AssertFeedCreated(t, orga, "recommendation_confirmation")

	if feed[0].TriggeredByUserId != recommendTo.ID {
		t.Fatalf("Feed must have been created by user the article was recommended to, but was: %v", feed[0].TriggeredByUserId)
	}

	access := model.FindFeedAccessByFeedId(feed[0].ID)

	if len(access) != 1 || access[0].UserId != user2.ID {
		t.Fatalf("Only second user must have received a read confirmation notification, but was: %v", access)
	}
}
