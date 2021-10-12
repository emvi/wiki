package article

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"strings"
	"testing"
)

func TestCopyArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	article.Pinned = true
	article.Views = 123

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	input := []struct {
		orga      *model.Organization
		userId    hide.ID
		articleId hide.ID
	}{
		{orga, 0, 0},
		{orga, user.ID, 0},
		{orga, user.ID, article.ID + 1},
		{orga, user.ID + 1, article.ID},
		{orga, user.ID, article.ID},
	}
	expected := []error{
		errs.ArticleNotFound,
		errs.ArticleNotFound,
		errs.ArticleNotFound,
		errs.PermissionDenied,
		nil,
	}

	var id hide.ID

	for i, in := range input {
		var err error
		id, err = CopyArticle(in.orga, in.userId, in.articleId, 0)

		if err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "copy_article")

	if id == 0 {
		t.Fatal("ID must be returned")
	}

	newArticle := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if newArticle == nil {
		t.Fatal("Article must have been copied")
	}

	if newArticle.Views != 0 || newArticle.Pinned {
		t.Fatal("Views and pinned must not be copied")
	}

	content := model.FindArticleContentByArticleId(id)

	if len(content) != 3 {
		t.Fatal("Article content must have been copied")
	}

	for _, c := range content {
		author := model.FindArticleContentAuthorByArticleContentId(c.ID)

		if len(author) != 1 {
			t.Fatal("Article content author must have been copied")
		}

		if c.Version == 0 && !strings.Contains(c.Title, "(copy)") {
			t.Fatalf("Title must contain '(copy)', but was: %v", c.Title)
		} else if c.Version != 0 && strings.Contains(c.Title, "(copy") {
			t.Fatalf("Title must not contain '(copy)', but was: %v", c.Title)
		}
	}

	tags := model.FindTagByOrganizationIdAndUserIdAndArticleId(orga.ID, user.ID, id)

	if len(tags) != 4 {
		t.Fatal("Tags mustt have been copied")
	}
}

func TestCopyArticleContentForLanguageNotFound(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "ru", "Russian", true)
	lang := testutil.CreateLang(t, orga, "en", "English", false)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)

	if _, err := CopyArticle(orga, user.ID, article.ID, 0); err != nil {
		t.Fatalf("Article must have been copied even if language for content does not exist, but was: %v", err)
	}

	feed := testutil.AssertFeedCreated(t, orga, "copy_article")
	refs := model.FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orga.ID, lang.ID, feed[0].ID)

	if len(refs) != 2 {
		t.Fatalf("Two feed references must have been created, but was: %v", len(refs))
	}
}
