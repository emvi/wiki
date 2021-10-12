package pinned

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestIsPinned(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	if IsPinned(orga.ID, article.ID, 0) {
		t.Fatal("Article must not be pinned")
	}

	if IsPinned(orga.ID, 0, list.ID) {
		t.Fatal("List must not be pinned")
	}

	article.Pinned = true
	list.Pinned = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveArticleList(nil, list); err != nil {
		t.Fatal(err)
	}

	if !IsPinned(orga.ID, article.ID, 0) {
		t.Fatal("Article must be pinned")
	}

	if !IsPinned(orga.ID, 0, list.ID) {
		t.Fatal("List must be pinned")
	}
}
