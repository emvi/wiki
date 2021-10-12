package search

import (
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSearchTag(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	tag1 := testutil.CreateTag(t, orga, "Test")
	tag2 := testutil.CreateTag(t, orga, "Teest")
	tag3 := testutil.CreateTag(t, orga, "A test")
	testutil.CreateArticleTag(t, article, tag1)
	testutil.CreateArticleTag(t, article, tag2)
	testutil.CreateArticleTag(t, article, tag3)
	ctx := context.NewEmviUserContext(orga, user.ID)
	results, count := SearchTag(ctx, "  ", nil)

	if len(results) != 7 || count != 7 {
		t.Fatalf("All results must be returned, but was: %v %v", len(results), count)
	}

	results, count = SearchTag(ctx, "query", nil)

	if len(results) != 0 || count != 0 {
		t.Fatalf("No results must be returned, but was: %v %v", len(results), count)
	}

	results, count = SearchTag(ctx, "test", nil)

	if len(results) != 4 || count != 4 {
		t.Fatalf("Four results must be returned, but was %v %v", len(results), count)
	}
}

func TestSearchTagFilter(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	tag1 := testutil.CreateTag(t, orga, "Test")
	tag2 := testutil.CreateTag(t, orga, "Teest")
	tag3 := testutil.CreateTag(t, orga, "A test")
	testutil.CreateArticleTag(t, article, tag1)
	testutil.CreateArticleTag(t, article, tag2)
	testutil.CreateArticleTag(t, article, tag3)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchTagFilter{
		SortUsages: "desc",
		SortName:   "asc",
	}
	results, count := SearchTag(ctx, "test", filter)

	if len(results) != 4 || count != 4 {
		t.Fatalf("Four results must be returned, but was %v %v", len(results), count)
	}
}

func TestSearchTagAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user2, lang, false, false)
	tag := testutil.CreateTag(t, orga, "uniquename")
	testutil.CreateArticleTag(t, article, tag)
	ctx := context.NewEmviUserContext(orga, user.ID)
	result, count := SearchTag(ctx, "uniquename", nil)

	if len(result) != 0 || count != 0 {
		t.Fatalf("Expected no tag to be returned, but was: %v %v", len(result), count)
	}

	testutil.CreateArticleAccess(t, article, user, nil, false)
	result, count = SearchTag(ctx, "uniquename", nil)

	if len(result) != 1 || count != 1 {
		t.Fatalf("Expected tag to be returned, but was: %v %v", len(result), count)
	}
}

func TestSearchTagClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	article2 := testutil.CreateArticle(t, orga, user, lang, false, false)
	tag := testutil.CreateTag(t, orga, "uniquename")
	testutil.CreateArticleTag(t, article, tag)
	testutil.CreateArticleTag(t, article2, tag)
	ctx := context.NewEmviContext(orga, 0, nil, false) // don't need to set scopes in test here
	results, count := SearchTag(ctx, "uniquename", nil)

	if len(results) != 0 || count != 0 {
		t.Fatalf("Expected no tag to be returned, but was: %v %v", len(results), count)
	}

	testutil.SetArticleClientAccess(t, article)
	results, count = SearchTag(ctx, "uniquename", nil)

	if len(results) != 1 || count != 1 {
		t.Fatalf("Expected tag to be returned, but was: %v %v", len(results), count)
	}

	// just one because the second article using this tag does not allow client access
	if results[0].Usages != 1 {
		t.Fatalf("Expected one usage to be returned, but was: %v", results[0].Usages)
	}
}
