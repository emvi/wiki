package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestAddArticleListEntry(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	input := []struct {
		userId    hide.ID
		listId    hide.ID
		articleId hide.ID
	}{
		{user.ID, 0, 0},
		{0, list.ID, 0},
		{user.ID, list.ID, 0},
		{user.ID, list.ID, article.ID},
		{user.ID, list.ID, article.ID},
	}
	expected := []error{
		errs.ArticleListNotFound,
		errs.PermissionDenied,
		errs.ArticlePermissionDenied,
		nil,
		nil,
	}

	var entries []model.ArticleListEntry

	for i, in := range input {
		added, err := AddArticleListEntry(orga, in.userId, in.listId, []hide.ID{in.articleId})

		if err != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		}

		if len(added) != 0 {
			entries = added
		}
	}

	if entries == nil || len(entries) != 1 {
		t.Fatalf("Only one new entry must have been returned, but was: %v", len(entries))
	}

	entry := entries[0]

	if entry.ArticleListId != list.ID ||
		entry.ArticleId != article.ID ||
		entry.Position != 1 ||
		entry.Article == nil ||
		entry.Article.LatestArticleContent == nil ||
		entry.Article.LatestArticleContent.Authors == nil ||
		len(entry.Article.LatestArticleContent.Authors) != 1 ||
		entry.Article.LatestArticleContent.Authors[0].ID != user.ID {
		t.Fatalf("Entry not as expected: %v", entry)
	}

	testutil.AssertFeedCreatedN(t, orga, "add_article_list_entry", 2)
}

func TestAddArticleListEntryMultipleTranslations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	list, _ := testutil.CreateArticleList(t, orga, user, langEn, true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleContent(t, user, article, langDe, 0)

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID}); err != nil {
		t.Fatalf("Entry must have been added to list, but was: %v", err)
	}

	feeds := testutil.AssertFeedCreated(t, orga, "add_article_list_entry")
	feedId := feeds[0].ID
	refs := model.FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orga.ID, langEn.ID, feedId)

	if len(refs) != 2 {
		t.Fatalf("Expected list entry feed to have 2 references, but was: %v", len(refs))
	}
}

func TestAddArticleListEntryArticlesCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID}); err != nil {
		t.Fatalf("Entry must have been added to list, but was: %v", err)
	}

	list = model.GetArticleListByOrganizationIdAndUserIdAndId(orga.ID, user.ID, list.ID)

	if list.ArticleCount != 1 {
		t.Fatalf("Article list must have one entry, but was: %v", list.ArticleCount)
	}

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article2.ID}); err != nil {
		t.Fatalf("Entry must have been added to list, but was: %v", err)
	}

	list = model.GetArticleListByOrganizationIdAndUserIdAndId(orga.ID, user.ID, list.ID)

	if list.ArticleCount != 2 {
		t.Fatalf("Article list must have two entries, but was: %v", list.ArticleCount)
	}
}

func TestAddArticleListEntryProtectedArticles(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, false, false)

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID, article2.ID}); err != nil {
		t.Fatalf("Entries must have been added to list, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "add_protected_article_list_entry")
}

func TestAddArticleListEntryArchived(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.Archived = null.NewString("archived", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if _, err := AddArticleListEntry(orga, user.ID, list.ID, []hide.ID{article.ID}); err != nil {
		t.Fatalf("Entries must have been added to list, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "add_article_list_entry")
}
