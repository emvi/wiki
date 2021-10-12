package history

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestDeleteHistoryEntryFailure(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := model.FindArticleContentByArticleId(article.ID) // version: 0 (latest), 1, 2 (last)
	input := []struct {
		userId    hide.ID
		contentId hide.ID
	}{
		{0, 0},
		{0, content[1].ID},
		{user.ID, content[0].ID},
	}
	expected := []error{
		errs.ArticleContentVersionNotFound,
		errs.PermissionDenied,
		errs.ArticleContentVersionInvalid,
	}

	for i, in := range input {
		ctx := context.NewEmviUserContext(orga, in.userId)

		if err := DeleteHistoryEntry(ctx, in.contentId); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	testutil.AssertFeedCreatedN(t, orga, "delete_article_history_entry", 0)
}

func TestDeleteHistoryEntrySuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := model.FindArticleContentByArticleId(article.ID) // version: 0 (latest), 1, 2 (last)
	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := DeleteHistoryEntry(ctx, content[1].ID); err != nil {
		t.Fatalf("Expected content version to be deleted, but was: %v", err)
	}

	if model.GetArticleContentById(content[1].ID) != nil {
		t.Fatal("Content version must have been deleted")
	}

	testutil.AssertFeedCreated(t, orga, "delete_article_history_entry")
}

func TestDeleteHistoryEntrySuccessUpdateLatest(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := model.FindArticleContentByArticleId(article.ID) // version: 0 (latest), 1, 2 (last)
	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := DeleteHistoryEntry(ctx, content[2].ID); err != nil {
		t.Fatalf("Expected content version to be deleted, but was: %v", err)
	}

	if model.GetArticleContentById(content[2].ID) != nil {
		t.Fatal("Content version must have been deleted")
	}

	latestContent := model.GetArticleContentById(content[0].ID)

	if latestContent == nil ||
		latestContent.Title != content[1].Title ||
		latestContent.Content != content[1].Content ||
		latestContent.Commit != content[1].Commit {
		t.Fatalf("Latest content must have been updated, but was: %v", latestContent)
	}

	testutil.AssertFeedCreated(t, orga, "delete_article_history_entry")
}

func TestDeleteHistoryEntryLastEntry(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := model.FindArticleContentByArticleId(article.ID) // version: 0 (latest), 1, 2 (last)
	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := DeleteHistoryEntry(ctx, content[1].ID); err != nil {
		t.Fatalf("Expected content version to be deleted, but was: %v", err)
	}

	if err := DeleteHistoryEntry(ctx, content[2].ID); err != errs.ArticleContentRemainingVersion {
		t.Fatalf("Remaining content version must not be deleted, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "delete_article_history_entry")
}

func TestDeleteHistoryEntryFeedRefs(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := testutil.CreateArticleContent(t, user, article, lang, 3)
	feed := &model.Feed{Public: false,
		Reason:            "reason",
		OrganizationId:    orga.ID,
		TriggeredByUserId: user.ID}

	if err := model.SaveFeed(nil, feed); err != nil {
		t.Fatal(err)
	}

	ref := &model.FeedRef{FeedId: feed.ID,
		ArticleContentID: content.ID}

	if err := model.SaveFeedRef(nil, ref); err != nil {
		t.Fatal(err)
	}

	access := &model.FeedAccess{FeedId: feed.ID,
		UserId: user.ID}

	if err := model.SaveFeedAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := DeleteHistoryEntry(ctx, content.ID); err != nil {
		t.Fatalf("Article history entry must have been deleted, but was: %v", err)
	}
}
