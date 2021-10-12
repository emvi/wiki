package articlelist

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"

	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
)

func TestRemoveArticleListEntry(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleListEntry(t, list, article, 1)
	emptyIds := []hide.ID{}

	input := []struct {
		UserId     hide.ID
		ListId     hide.ID
		ArticleIds []hide.ID
	}{
		{0, 0, emptyIds},
		{0, list.ID, emptyIds},
		{user.ID, list.ID, []hide.ID{article.ID}},
	}
	expected := []error{
		errs.ArticleListNotFound,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		if err := RemoveArticleListEntry(orga, in.UserId, in.ListId, in.ArticleIds); err != expected[i] {
			t.Fatalf("Expected %v but was: %v", expected[i], err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "remove_article_list_entry")
}

func TestRemoveArticleListEntryArticlesCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID, article2.ID}); err != nil {
		t.Fatal(err)
	}

	list = model.GetArticleListByOrganizationIdAndUserIdAndId(orga.ID, user.ID, list.ID)

	if list.ArticleCount != 2 {
		t.Fatalf("Article list must have two entries, but was: %v", list.ArticleCount)
	}

	if err := RemoveArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID, article2.ID}); err != nil {
		t.Fatalf("Articles must have been removed, but was: %v", err)
	}

	list = model.GetArticleListByOrganizationIdAndUserIdAndId(orga.ID, user.ID, list.ID)

	if list.ArticleCount != 0 {
		t.Fatalf("Article list must have no entries, but was: %v", list.ArticleCount)
	}
}

func TestRemoveArticleListEntryProtectedArticles(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, false, false)

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID, article2.ID}); err != nil {
		t.Fatal(err)
	}

	if err := RemoveArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID, article2.ID}); err != nil {
		t.Fatalf("Articles must have been removed, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "remove_article_list_entry")
	testutil.AssertFeedCreated(t, orga, "remove_protected_article_list_entry")
}

func TestRemoveArticleListEntryArchived(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleListEntry(t, list, article, 1)
	article.Archived = null.NewString("archived", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := RemoveArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID}); err != nil {
		t.Fatalf("Articles must have been removed, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "remove_article_list_entry")
}
