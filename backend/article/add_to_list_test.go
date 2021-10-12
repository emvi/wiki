package article

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestAddArticleToLists(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)

	if err := AddArticleToLists(orga, user.ID, article.ID, []hide.ID{}); err != nil {
		t.Fatalf("Must not return error, but was: %v", err)
	}

	if err := AddArticleToLists(orga, user.ID, article.ID, []hide.ID{list.ID}); err != nil {
		t.Fatalf("Article must be added to list, but was: %v", err)
	}

	entries := model.FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orga.ID, user.ID, lang.ID, list.ID, new(model.SearchArticleListEntryFilter))

	if len(entries) != 1 {
		t.Fatalf("List must have one entry, but was: %v", len(entries))
	}
}
