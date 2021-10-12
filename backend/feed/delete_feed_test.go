package feed

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDeleteFeed(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	feed := testutil.CreateFeed(t, orga, user, lang, false)

	ref := &model.FeedRef{FeedId: feed.ID, ArticleID: article.ID}

	if err := model.SaveFeedRef(nil, ref); err != nil {
		t.Fatal(err)
	}

	access := &model.FeedAccess{UserId: user.ID, FeedId: feed.ID}

	if err := model.SaveFeedAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	if err := DeleteFeed(nil, &DeleteFeedData{ArticleId: article.ID}); err != nil {
		t.Fatalf("Feed must have been deleted, but was: %v", err)
	}
}
