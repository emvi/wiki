package util

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestRemoveNonPublicInformation(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	input := [][]string{
		nil,
		{client.Scopes["article_authors"].String()},
		{client.Scopes["article_authors"].String(), client.Scopes["article_authors_mails"].String()},
	}
	expected := []struct {
		hasAuthors     bool
		hasAuthorMails bool
	}{
		{false, false},
		{true, false},
		{true, true},
	}

	for i, in := range input {
		article.LatestArticleContent.User = user
		article.LatestArticleContent.Authors[0].OrganizationMember = user.OrganizationMember
		article.LatestArticleContent.Authors[0].AcceptMarketing = true
		article.LatestArticleContent.Authors[0].Email = user.Email
		content := []model.ArticleContent{*article.LatestArticleContent}
		ctx := context.NewEmviContext(orga, 0, in, false)
		RemoveNonPublicInformation(ctx, content)

		if content[0].User != nil {
			t.Fatal("Must not have users")
		}

		if !expected[i].hasAuthors && content[0].Authors != nil {
			t.Fatal("Must not have authors")
		} else if expected[i].hasAuthors && content[0].Authors == nil {
			t.Fatal("Must have authors")
		} else if expected[i].hasAuthors && content[0].Authors != nil {
			if !expected[i].hasAuthorMails && content[0].Authors[0].Email != "" {
				t.Fatal("Must not have author mails")
			} else if expected[i].hasAuthorMails && content[0].Authors[0].Email == "" {
				t.Fatal("Must have author mails")
			}
		}
	}
}
