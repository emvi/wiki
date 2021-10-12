package util

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
)

// RemoveAuthorsOrAuthorMails removes the authors or author mails depending on
// if the context was created by a client call and if the client lacks permissions to read authors/author mails.
func RemoveAuthorsOrAuthorMails(ctx context.EmviContext, articles []model.Article) {
	if ctx.IsClient() {
		if !ctx.HasScopes(client.Scopes["article_authors"]) {
			for i := range articles {
				if articles[i].LatestArticleContent != nil {
					articles[i].LatestArticleContent.Authors = nil
				}
			}
		} else if !ctx.HasScopes(client.Scopes["article_authors_mails"]) {
			for i := range articles {
				if articles[i].LatestArticleContent != nil {
					for j := range articles[i].LatestArticleContent.Authors {
						articles[i].LatestArticleContent.Authors[j].Email = ""
					}
				}
			}
		}

		// remove non public information
		for i := range articles {
			if articles[i].LatestArticleContent != nil {
				RemoveNonPublicInformationFromContent(ctx, articles[i].LatestArticleContent)
			}
		}
	}
}
