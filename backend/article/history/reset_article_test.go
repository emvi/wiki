package history

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestResetArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	userNoAccess := testutil.CreateUser(t, orga, 321, "noaccess@user.com")

	input := []struct {
		userId    hide.ID
		articleId hide.ID
		langId    hide.ID
		version   int
		commit    string
	}{
		{0, 0, 0, 0, ""},
		{user.ID, 0, 0, 0, ""},
		{userNoAccess.ID, article.ID, 0, 0, ""},
		{user.ID, article.ID, lang.ID, 0, ""},
		{user.ID, article.ID, lang.ID, 1, "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891"},
		{user.ID, article.ID, lang.ID, 2, "Reset article"},
		{user.ID, article.ID, lang.ID, 99, "Reset article"},
		{user.ID, article.ID, lang.ID, 1, "Reset article"},
	}
	expected := []error{
		errs.ArticleNotFound,
		errs.ArticleNotFound,
		errs.PermissionDenied,
		errs.ArticleContentVersionInvalid,
		errs.CommitMsgLen,
		errs.ArticleContentVersionInvalid,
		errs.ArticleContentVersionNotFound,
		nil,
	}

	for i, in := range input {
		if err := ResetArticle(orga, in.userId, in.articleId, in.langId, in.version, in.commit); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	latestContent := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, lang.ID, true)

	if latestContent == nil ||
		latestContent.Commit.String != "Reset article" ||
		latestContent.Version != 0 ||
		latestContent.UserId != user.ID ||
		latestContent.WIP {
		t.Fatalf("Unexpected latest content, was: %v", latestContent)
	}

	lastContent := model.GetArticleContentLastByArticleIdAndLanguageIdAndWIP(article.ID, lang.ID, false)

	if lastContent == nil ||
		lastContent.Commit.String != "Reset article" ||
		lastContent.Version != 3 ||
		lastContent.UserId != user.ID ||
		lastContent.WIP {
		t.Fatalf("Unexpected last content, was: %v", lastContent)
	}

	testutil.AssertFeedCreated(t, orga, "reset_article")
}

func TestResetArticleArchived(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	article.Archived = null.NewString("Archived", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := ResetArticle(orga, user.ID, article.ID, lang.ID, 1, "Reset article"); err != errs.ArticleNotFound {
		t.Fatalf("Archived article must not be reset, but was: %v", err)
	}
}

func TestResetArticleNonExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	testutil.CreateArticleContent(t, user, article, lang, 4)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := ResetArticle(orga, user.ID, article.ID, lang.ID, 1, "Reset article"); err != errs.RequiresExpertVersion {
		t.Fatalf("Resetting an article to a version outside of non expert range must return error, but was: %v", err)
	}
}
