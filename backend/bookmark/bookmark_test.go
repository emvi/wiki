package bookmark

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestBookmarkObject(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	input := []struct {
		orga      *model.Organization
		userId    hide.ID
		articleId hide.ID
		listId    hide.ID
	}{
		{orga, user.ID, 0, 0},
		{orga, user.ID, article.ID + 1, 0},
		{orga, user.ID, article.ID, 0},
		{orga, user.ID, 0, list.ID + 1},
		{orga, user.ID, 0, list.ID},
	}
	expected := []error{
		errs.NoObjectToBookmark,
		errs.ArticleNotFound,
		nil,
		errs.ArticleListNotFound,
		nil,
	}

	for i, in := range input {
		if err := BookmarkObject(in.orga, in.userId, in.articleId, in.listId); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}

func TestBookmarkObjectArticleToggle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := BookmarkObject(orga, user.ID, article.ID, 0); err != nil {
		t.Fatal(err)
	}

	if model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, user.ID, article.ID, 0) == nil {
		t.Fatal("Article bookmark must exist")
	}

	if err := BookmarkObject(orga, user.ID, article.ID, 0); err != nil {
		t.Fatal(err)
	}

	if model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, user.ID, article.ID, 0) != nil {
		t.Fatal("Article bookmark must not exist")
	}
}

func TestBookmarkObjectListToggle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	if err := BookmarkObject(orga, user.ID, 0, list.ID); err != nil {
		t.Fatal(err)
	}

	if model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, user.ID, 0, list.ID) == nil {
		t.Fatal("List bookmark must exist")
	}

	if err := BookmarkObject(orga, user.ID, 0, list.ID); err != nil {
		t.Fatal(err)
	}

	if model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, user.ID, 0, list.ID) != nil {
		t.Fatal("List bookmark must not exist")
	}
}

func TestBookmarkObjectArchivedArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.Archived = null.NewString("Test", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := BookmarkObject(orga, user.ID, article.ID, 0); err != nil {
		t.Fatalf("Article must have been bookmarked, but was: %v", err)
	}

	if model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, user.ID, article.ID, 0) == nil {
		t.Fatal("Article bookmark must exist")
	}

	if err := BookmarkObject(orga, user.ID, article.ID, 0); err != nil {
		t.Fatalf("Article must have been bookmarked, but was: %v", err)
	}

	if model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, user.ID, article.ID, 0) != nil {
		t.Fatal("Article bookmark must not exist")
	}
}
