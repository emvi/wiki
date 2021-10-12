package util

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestGetArticleWithAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	ctx := context.NewEmviContext(orga, user.ID, nil, false)
	out, err := GetArticleWithAccess(nil, ctx, article.ID, false)

	if err != nil {
		t.Fatalf("Article must have been found, but was: %v", err)
	}

	if out.ID != article.ID || len(out.Access) != 1 {
		t.Fatalf("Article must have been returned together with access, but was: %v %v", out.ID, len(out.Access))
	}
}

func TestGetArticleWithAccessArchived(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.Archived.SetValid("archived article")

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	ctx := context.NewEmviContext(orga, user.ID, nil, false)
	out, err := GetArticleWithAccess(nil, ctx, article.ID, false)

	if err != errs.ArticleNotFound || out != nil {
		t.Fatal("Article must not be found")
	}

	out, err = GetArticleWithAccess(nil, ctx, article.ID, true)

	if err != nil {
		t.Fatalf("Article must have been found, but was: %v", err)
	}

	if out.ID != article.ID || len(out.Access) != 1 {
		t.Fatalf("Article must have been returned together with access, but was: %v %v", out.ID, len(out.Access))
	}
}

func TestGetArticleWithAccessClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	ctx := context.NewEmviContext(orga, 0, []string{client.Scopes["articles"].String()}, false)
	out, err := GetArticleWithAccess(nil, ctx, article.ID, false)

	if err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied for non public article, but was: %v", err)
	}

	article.ClientAccess = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	out, err = GetArticleWithAccess(nil, ctx, article.ID, false)

	if err != nil {
		t.Fatalf("Article must have been found, but was: %v", err)
	}

	if out.ID != article.ID || len(out.Access) != 0 {
		t.Fatalf("Article must have been returned without access, but was: %v %v", out.ID, len(out.Access))
	}
}

func TestCheckCommitMsg(t *testing.T) {
	if err := CheckCommitMsg("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891"); err != errs.CommitMsgLen {
		t.Fatalf("Expected commit message to be too long, but was: %v", err)
	}

	if err := CheckCommitMsg("okay"); err != nil {
		t.Fatalf("Expected commit message to be okay, but was: %v", err)
	}
}
