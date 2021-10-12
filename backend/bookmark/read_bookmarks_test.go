package bookmark

import (
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadBookmarks(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article1 := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)
	article3 := testutil.CreateArticle(t, orga, user, lang, true, true)
	list1, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	list2, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	if err := BookmarkObject(orga, user.ID, article1.ID, 0); err != nil {
		t.Fatal(err)
	}

	if err := BookmarkObject(orga, user.ID, article2.ID, 0); err != nil {
		t.Fatal(err)
	}

	if err := BookmarkObject(orga, user.ID, article3.ID, 0); err != nil {
		t.Fatal(err)
	}

	if err := BookmarkObject(orga, user.ID, 0, list1.ID); err != nil {
		t.Fatal(err)
	}

	if err := BookmarkObject(orga, user.ID, 0, list2.ID); err != nil {
		t.Fatal(err)
	}

	articles, lists := ReadBookmarks(orga, user.ID, true, true, 0, 0)

	if len(articles) != 3 {
		t.Fatalf("Bookmarked articles must have been found, but was: %v", len(articles))
	}

	if len(lists) != 2 {
		t.Fatalf("Bookmarked lists must have been found, but was: %v", len(lists))
	}
}
