package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
	"time"
)

func TestRemoveTag(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)

	input := []struct {
		articleId hide.ID
		userId    hide.ID
		tag       string
	}{
		{0, user.ID, "tag"},
		{article.ID, user.ID + 1, "tag"},
		{article.ID, user.ID, "tag"},
	}
	expected := []error{
		errs.ArticleNotFound,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		if err := RemoveTag(orga, in.userId, in.articleId, in.tag); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	if err := AddTag(orga, AddTagData{article.ID, "tag"}); err != nil {
		t.Fatal(err)
	}

	if !checkArticleHasTag(t, orga.ID, article.ID, "tag") {
		t.Fatal("Article must have tag 'tag'")
	}

	if err := RemoveTag(orga, user.ID, article.ID, "tag"); err != nil {
		t.Fatalf("Existing tag must have been removed, but was %v", err)
	}

	if checkArticleHasTag(t, orga.ID, article.ID, "tag") {
		t.Fatal("Article must not have tag 'tag'")
	}

	time.Sleep(time.Millisecond * 10) // wait for goroutine to finish
}

func TestRemoveTagArchivedArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	tag := testutil.CreateTag(t, orga, "test")
	testutil.CreateArticleTag(t, article, tag)
	article.Archived = null.NewString("archived", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := RemoveTag(orga, user.ID, article.ID, "test"); err != nil {
		t.Fatalf("Expected tag to be removed, but was: %v", err)
	}
}

func checkArticleHasTag(t *testing.T, orgaId, articleId hide.ID, tag string) bool {
	return model.GetArticleTagByOrganizationIdAndArticleIdAndName(orgaId, articleId, tag) != nil
}
