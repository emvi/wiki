package article

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDeleteArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := ArchiveArticle(orga, user.ID, article.ID, "test", false); err != nil {
		t.Fatal(err)
	}

	if err := DeleteArticle(orga, user.ID, article.ID); err != nil {
		t.Fatalf("Article must have been deleted, but was: %v", err)
	}

	feed := testutil.AssertFeedCreated(t, orga, "delete_article")
	refs := model.FindFeedRefByFeedId(feed[0].ID)

	if len(refs) != 1 || refs[0].Key.String != "name" || refs[0].Value.String != "title 2" {
		t.Fatalf("Feed must have deleted object name, but was: %v", refs)
	}
}

func TestDeleteArticleReferencedObjects(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateObservedObject(t, user, article, nil, nil)
	testutil.CreateFeedForObject(t, orga, user, article, nil, nil)
	testutil.CreateArticleVisit(t, article, user)

	if err := DeleteArticle(orga, user.ID, article.ID); err != nil {
		t.Fatalf("Article must have been deleted including all references, but was: %v", err)
	}
}

func TestDeleteArticleUpdateArticleLists(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateArticleListEntry(t, list, article, 1)

	if err := DeleteArticle(orga, user.ID, article.ID); err != nil {
		t.Fatalf("Article must have been deleted, but was: %v", err)
	}

	list = model.GetArticleListByOrganizationIdAndUserIdAndId(orga.ID, user.ID, list.ID)

	if list.ArticleCount != 0 {
		t.Fatal("List article count must be 0")
	}
}

func TestDeleteArticleRecommended(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleRecommendation(t, article, user, user2)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := DeleteArticle(orga, user.ID, article.ID); err != nil {
		t.Fatalf("Recommended article must have been deleted, but was: %v", err)
	}
}
