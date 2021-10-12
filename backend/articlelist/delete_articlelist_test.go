package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDeleteArticleList(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)

	if err := DeleteArticleList(orga, user.ID, 1); err != errs.ArticleListNotFound {
		t.Fatalf("Article list must not be found, but was: %v", err)
	}

	list, _ := testutil.CreateArticleList(t, orga, nil, lang, true)

	if err := DeleteArticleList(orga, user.ID, list.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	testutil.CreateArticleListMember(t, list, user.ID, 0, true)

	if err := DeleteArticleList(orga, user.ID, list.ID); err != nil {
		t.Fatalf("Article list must have been deleted, but was %v", err)
	}

	if model.GetArticleListByOrganizationIdAndId(orga.ID, list.ID) != nil {
		t.Fatal("Article list must not exist anymore")
	}

	feed := testutil.AssertFeedCreated(t, orga, "delete_articlelist")
	refs := model.FindFeedRefByFeedId(feed[0].ID)

	if len(refs) != 1 || refs[0].Key.String != "name" || refs[0].Value.String != "article list name" {
		t.Fatalf("Feed must have deleted object name, but was: %v", refs)
	}
}

func TestDeleteArticleListReferencedObjects(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleListEntry(t, list, article, 0)
	testutil.CreateObservedObject(t, user, nil, list, nil)
	testutil.CreateFeedForObject(t, orga, user, nil, list, nil)

	if err := DeleteArticleList(orga, user.ID, list.ID); err != nil {
		t.Fatalf("Article list must have been deleted including all references, but was: %v", err)
	}
}
