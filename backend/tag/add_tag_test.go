package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"strconv"
	"testing"
)

func TestAddTag(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	tag := testutil.CreateTag(t, orga, "test")
	testutil.CreateArticleTag(t, article, tag)

	input := []struct {
		ArticleId hide.ID
		Tag       string
	}{
		{0, "tag"},
		{0, "  "},
		{article.ID, "test"},
		{article.ID, "$invalid"},
		{article.ID, "0123456789012345678901234567890123456789012345678901234567891"},
		{article.ID, "0123.T-a_g.o k"},
	}
	expected := []error{
		errs.ArticleNotFound,
		errs.TagEmpty,
		errs.TagExistsAlready,
		errs.TagInvalid,
		errs.TagLen,
		nil,
	}

	for i, in := range input {
		if err := AddTag(orga, AddTagData{in.ArticleId, in.Tag}); err != expected[i] {
			t.Fatalf("Expected %v but was: %v", expected[i], err)
		}
	}

	if !checkArticleHasTag(t, orga.ID, article.ID, "0123.T-a_g.o k") {
		t.Fatal("Article must have tag '0123.T-a_g.o k'")
	}
}

func TestAddTagMaxTags(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)

	for i := 0; i < MaxTagsPerArticle; i++ {
		tag := testutil.CreateTag(t, orga, "test"+strconv.Itoa(i))
		testutil.CreateArticleTag(t, article, tag)
	}

	if err := AddTag(orga, AddTagData{article.ID, "testMax"}); err != errs.MaxTagsReached {
		t.Fatalf("Expected %v but was: %v", errs.MaxTagsReached, err)
	}
}

func TestAddTagArchivedArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.Archived = null.NewString("archived", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := AddTag(orga, AddTagData{article.ID, "newtag"}); err != nil {
		t.Fatalf("Expected tag to be added, but was: %v", err)
	}
}
