package util

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestRemoveAuthorsOrAuthorMails(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	articles := []model.Article{*testutil.CreateArticle(t, orga, user, lang, true, true)}
	ctx := context.NewEmviContext(orga, 0, nil, false)
	RemoveAuthorsOrAuthorMails(ctx, articles)

	if articles[0].LatestArticleContent.Authors != nil {
		t.Fatal("Authors must have been removed")
	}

	articles = []model.Article{*testutil.CreateArticle(t, orga, user, lang, true, true)}
	ctx = context.NewEmviContext(orga, 0, []string{client.Scopes["article_authors"].String()}, false)
	RemoveAuthorsOrAuthorMails(ctx, articles)

	if len(articles[0].LatestArticleContent.Authors) != 1 || articles[0].LatestArticleContent.Authors[0].Email != "" {
		t.Fatal("Author mails must have been removed")
	}

	articles = []model.Article{*testutil.CreateArticle(t, orga, user, lang, true, true)}
	ctx = context.NewEmviContext(orga, 0, []string{client.Scopes["article_authors"].String(), client.Scopes["article_authors_mails"].String()}, false)
	RemoveAuthorsOrAuthorMails(ctx, articles)

	if len(articles[0].LatestArticleContent.Authors) != 1 ||
		articles[0].LatestArticleContent.Authors[0].Email == "" ||
		articles[0].LatestArticleContent.Authors[0].OrganizationMember != nil ||
		articles[0].LatestArticleContent.User != nil {
		t.Fatal("Authors must have be complete")
	}
}
