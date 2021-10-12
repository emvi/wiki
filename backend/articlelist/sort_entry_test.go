package articlelist

import (
	"github.com/emvi/hide"
	"testing"

	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
)

func TestSwapArticleListEntries(t *testing.T) {
	orga, user, _, list, articles, entries := createTestList(t)
	articleB := articles[1]
	articleC := articles[2]
	articleD := articles[3]
	entry3 := entries[2]
	entry4 := entries[3]

	in := []struct {
		userId   hide.ID
		listId   hide.ID
		articleA hide.ID
		ArticleB hide.ID
	}{
		{user.ID, list.ID, 0, 0},
		{0, list.ID, articleC.ID, articleD.ID},
		{user.ID, 0, articleC.ID, articleD.ID},
		{user.ID, list.ID, articleC.ID + 100, articleB.ID},
		{user.ID, list.ID, articleC.ID, articleD.ID + 100},
		{user.ID, list.ID, articleC.ID, articleD.ID},
	}
	out := []error{
		errs.InvalidPos,
		errs.PermissionDenied,
		errs.ArticleListNotFound,
		errs.PosNotFound,
		errs.PosNotFound,
		nil,
	}

	for i := range in {
		if err := SwapArticleListEntries(orga, in[i].userId, in[i].listId, in[i].articleA, in[i].ArticleB); err != out[i] {
			t.Fatalf("Error '%v' does not match expected error '%v'", err, out[i])
		}
	}

	newEntry3 := model.GetArticleListEntryByArticleListIdAndPosition(list.ID, 3)
	newEntry4 := model.GetArticleListEntryByArticleListIdAndPosition(list.ID, 4)

	if newEntry3.ID != entry4.ID {
		t.Fatal("Position 4 must have been swapped to 3")
	}

	if newEntry4.ID != entry3.ID {
		t.Fatal("Position 3 must have been swapped to 4")
	}
}

func TestSortArticleListEntry(t *testing.T) {
	orga, user, lang, list, articles, _ := createTestList(t)
	articleC := articles[2]

	in := []struct {
		userId    hide.ID
		listId    hide.ID
		article   hide.ID
		direction int
	}{
		{0, 0, 0, 1},
		{0, list.ID, 0, 1},
		{user.ID, list.ID, 0, 1},
		{user.ID, list.ID, articleC.ID, 0},
		{user.ID, list.ID, articleC.ID, 1},
	}
	out := []error{
		errs.ArticleListNotFound,
		errs.PermissionDenied,
		errs.PosNotFound,
		nil,
		nil,
	}

	for i := range in {
		if err := SortArticleListEntry(orga, in[i].userId, in[i].listId, in[i].article, in[i].direction); err != out[i] {
			t.Fatalf("Error '%v' does not match expected error '%v'", err, out[i])
		}
	}

	entries := model.FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orga.ID, user.ID, lang.ID, list.ID, &model.SearchArticleListEntryFilter{BaseSearch: model.BaseSearch{Limit: 100, Offset: 0}})

	if len(entries) != 5 {
		t.Fatalf("Five entries must have been returned, but was: %v", len(entries))
	}

	if entries[3].ID != articleC.ID {
		t.Fatal("Article must now be one position below original position")
	}
}

func TestSortArticleListEntryUp(t *testing.T) {
	orga, user, lang, list, articles, _ := createTestList(t)
	articleC := articles[2]

	if err := SortArticleListEntry(orga, user.ID, list.ID, articleC.ID, -1); err != nil {
		t.Fatalf("Entry position must have been updated, but was: %v", err)
	}

	entries := model.FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orga.ID, user.ID, lang.ID, list.ID, &model.SearchArticleListEntryFilter{BaseSearch: model.BaseSearch{Limit: 100, Offset: 0}})

	if len(entries) != 5 {
		t.Fatalf("Five entries must have been returned, but was: %v", len(entries))
	}

	if entries[1].ID != articleC.ID {
		t.Fatal("Article must now be one position above original position")
	}
}

func TestSortArticleListEntryFirst(t *testing.T) {
	orga, user, lang, list, articles, _ := createTestList(t)
	articleA := articles[0]

	if err := SortArticleListEntry(orga, user.ID, list.ID, articleA.ID, -1); err != nil {
		t.Fatalf("Entry position must have been updated, but was: %v", err)
	}

	entries := model.FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orga.ID, user.ID, lang.ID, list.ID, &model.SearchArticleListEntryFilter{BaseSearch: model.BaseSearch{Limit: 100, Offset: 0}})

	if len(entries) != 5 {
		t.Fatalf("Five entries must have been returned, but was: %v", len(entries))
	}

	if entries[4].ID != articleA.ID {
		t.Fatal("Article must now be on last position")
	}
}

func TestSortArticleListEntryLast(t *testing.T) {
	orga, user, lang, list, articles, _ := createTestList(t)
	articleE := articles[4]

	if err := SortArticleListEntry(orga, user.ID, list.ID, articleE.ID, 1); err != nil {
		t.Fatalf("Entry position must have been updated, but was: %v", err)
	}

	entries := model.FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orga.ID, user.ID, lang.ID, list.ID, &model.SearchArticleListEntryFilter{BaseSearch: model.BaseSearch{Limit: 100, Offset: 0}})

	if len(entries) != 5 {
		t.Fatalf("Five entries must have been returned, but was: %v", len(entries))
	}

	if entries[0].ID != articleE.ID {
		t.Fatal("Article must now be on first position")
	}
}

func createTestList(t *testing.T) (*model.Organization, *model.User, *model.Language, *model.ArticleList, []*model.Article, []*model.ArticleListEntry) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	articleA := testutil.CreateArticle(t, orga, user, lang, true, true)
	articleB := testutil.CreateArticle(t, orga, user, lang, true, true)
	articleC := testutil.CreateArticle(t, orga, user, lang, true, true)
	articleD := testutil.CreateArticle(t, orga, user, lang, true, true)
	articleE := testutil.CreateArticle(t, orga, user, lang, true, true)
	entry1 := testutil.CreateArticleListEntry(t, list, articleA, 1)
	entry2 := testutil.CreateArticleListEntry(t, list, articleB, 2)
	entry3 := testutil.CreateArticleListEntry(t, list, articleC, 3)
	entry4 := testutil.CreateArticleListEntry(t, list, articleD, 4)
	entry5 := testutil.CreateArticleListEntry(t, list, articleE, 5)
	return orga, user, lang, list, []*model.Article{articleA, articleB, articleC, articleD, articleE},
		[]*model.ArticleListEntry{entry1, entry2, entry3, entry4, entry5}
}
