package article

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadPrivateArticles(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	privateArticle := testutil.CreateArticle(t, orga, user, lang, false, false)
	privateArticle.Private = true

	if err := model.SaveArticle(nil, privateArticle); err != nil {
		t.Fatal(err)
	}

	articles := ReadPrivateArticles(orga, user.ID, 0)

	if len(articles) != 1 {
		t.Fatalf("Expected one articles, but was: %v", len(articles))
	}
}

func TestReadDrafts(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticle(t, orga, user, lang, true, false)
	testutil.CreateArticle(t, orga, user, lang, false, false)
	draft := testutil.CreateArticle(t, orga, user, lang, true, true)
	draft.Published.SetNil()

	if err := model.SaveArticle(nil, draft); err != nil {
		t.Fatal(err)
	}

	articles := ReadDrafts(orga, user.ID, 0)

	if len(articles) != 1 {
		t.Fatalf("Expected one articles, but was: %v", len(articles))
	}
}
