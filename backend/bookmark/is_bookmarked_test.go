package bookmark

import (
	"emviwiki/shared/testutil"
	"testing"
)

func TestIsBookmarkedArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if IsBookmarked(user.ID, article.ID, 0) {
		t.Fatal("Article must not be bookmarked")
	}

	if err := BookmarkObject(orga, user.ID, article.ID, 0); err != nil {
		t.Fatal(err)
	}

	if !IsBookmarked(user.ID, article.ID, 0) {
		t.Fatal("Article must be bookmarked")
	}
}

func TestIsBookmarkedList(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	if IsBookmarked(user.ID, 0, list.ID) {
		t.Fatal("Article must not be bookmarked")
	}

	if err := BookmarkObject(orga, user.ID, 0, list.ID); err != nil {
		t.Fatal(err)
	}

	if !IsBookmarked(user.ID, 0, list.ID) {
		t.Fatal("Article must be bookmarked")
	}
}
