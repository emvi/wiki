package pinned

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadPinnedArticles(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true) // not pinned
	article2 := testutil.CreateArticle(t, orga, user2, lang, true, true)
	article3 := testutil.CreateArticle(t, orga, user2, lang, false, false) // no access
	article4 := testutil.CreateArticle(t, orga, user, lang, false, false)
	setArticlePinned(t, article2)
	setArticlePinned(t, article3)
	setArticlePinned(t, article4)

	if articles, _, count, _ := ReadPinned(context.NewEmviUserContext(orga, user.ID), true, false, 0, 0); len(articles) != 2 || count != 2 {
		t.Fatalf("Expected two articles to be returned, but was: %v %v", len(articles), count)
	}
}

func TestReadPinnedArticlesClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true) // not pinned
	article2 := testutil.CreateArticle(t, orga, user2, lang, true, true)
	article3 := testutil.CreateArticle(t, orga, user2, lang, false, false) // no access
	article4 := testutil.CreateArticle(t, orga, user, lang, false, false)
	setArticlePinned(t, article2)
	setArticlePinned(t, article3)
	setArticlePinned(t, article4)
	ctx := context.NewEmviContext(orga, 0, []string{client.Scopes["articles"].String()}, false)

	if articles, _, count, _ := ReadPinned(ctx, true, false, 0, 0); len(articles) != 0 || count != 0 {
		t.Fatalf("Expected no articles to be returned, but was: %v %v", len(articles), count)
	}

	testutil.SetArticleClientAccess(t, article2)

	if articles, _, count, _ := ReadPinned(ctx, true, false, 0, 0); len(articles) != 1 || count != 1 {
		t.Fatalf("Expected one article to be returned, but was: %v %v", len(articles), count)
	}

	ctx.Scopes = nil

	if articles, _, count, _ := ReadPinned(ctx, true, false, 0, 0); len(articles) != 0 || count != 0 {
		t.Fatalf("Expected no articles to be returned, but was: %v %v", len(articles), count)
	}
}

func TestReadPinnedArticleLists(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticleList(t, orga, user, lang, true) // not pinned
	list2, _ := testutil.CreateArticleList(t, orga, user2, lang, true)
	list3, _ := testutil.CreateArticleList(t, orga, user2, lang, false) // no access
	list4, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	setListPinned(t, list2)
	setListPinned(t, list3)
	setListPinned(t, list4)

	if _, lists, _, count := ReadPinned(context.NewEmviUserContext(orga, user.ID), false, true, 0, 0); len(lists) != 2 || count != 2 {
		t.Fatalf("Expected two article lists to be returned, but was: %v %v", len(lists), count)
	}
}

func TestReadPinnedArticleListsClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "tester@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticleList(t, orga, user, lang, true) // not pinned
	list2, _ := testutil.CreateArticleList(t, orga, user2, lang, true)
	list3, _ := testutil.CreateArticleList(t, orga, user2, lang, false) // no access
	list4, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	setListPinned(t, list2)
	setListPinned(t, list3)
	setListPinned(t, list4)
	ctx := context.NewEmviContext(orga, 0, []string{client.Scopes["lists"].String()}, false)

	if _, lists, _, count := ReadPinned(ctx, false, true, 0, 0); len(lists) != 0 || count != 0 {
		t.Fatalf("Expected no article lists to be returned, but was: %v %v", len(lists), count)
	}

	testutil.SetListClientAccess(t, list2)

	if _, lists, _, count := ReadPinned(ctx, false, true, 0, 0); len(lists) != 1 || count != 1 {
		t.Fatalf("Expected one article lists to be returned, but was: %v %v", len(lists), count)
	}

	ctx.Scopes = nil

	if _, lists, _, count := ReadPinned(ctx, false, true, 0, 0); len(lists) != 0 || count != 0 {
		t.Fatalf("Expected no article lists to be returned, but was: %v %v", len(lists), count)
	}
}

func setArticlePinned(t *testing.T, article *model.Article) {
	article.Pinned = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}
}

func setListPinned(t *testing.T, list *model.ArticleList) {
	list.Pinned = true

	if err := model.SaveArticleList(nil, list); err != nil {
		t.Fatal(err)
	}
}
