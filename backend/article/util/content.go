package util

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
)

// RemoveNonPublicInformation removes information based on if the request was made by a client or by a user.
// In case the the request was made by a client, all non public information is removed depending on the client scopes.
func RemoveNonPublicInformation(ctx context.EmviContext, content []model.ArticleContent) {
	if ctx.IsClient() {
		for i := range content {
			RemoveNonPublicInformationFromContent(ctx, &content[i])
		}
	}
}

// RemoveNonPublicInformationFromContent see RemoveNonPublicInformation.
func RemoveNonPublicInformationFromContent(ctx context.EmviContext, content *model.ArticleContent) {
	if content != nil {
		content.User = nil

		if !ctx.HasScopes(client.Scopes["article_authors"]) {
			content.Authors = nil
		} else if !ctx.HasScopes(client.Scopes["article_authors_mails"]) {
			for i := range content.Authors {
				content.Authors[i].OrganizationMember = nil
				content.Authors[i].AcceptMarketing = false
				content.Authors[i].Email = ""
			}
		} else {
			for i := range content.Authors {
				content.Authors[i].OrganizationMember = nil
				content.Authors[i].AcceptMarketing = false
			}
		}
	}
}
