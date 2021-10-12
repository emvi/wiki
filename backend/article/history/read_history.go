package history

import (
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

const (
	articleMaxHistory = 20
)

func ReadArticleHistory(ctx context.EmviContext, articleId, langId hide.ID, offset int) ([]model.ArticleContent, int, error) {
	article := model.GetArticleByOrganizationIdAndIdIgnoreArchived(ctx.Organization.ID, articleId)

	if article == nil {
		return nil, 0, errs.ArticleNotFound
	}

	if (ctx.IsClient() && !article.ClientAccess) || (!article.ReadEveryone && model.FindArticleAccessByArticleIdAndUserId(articleId, ctx.UserId) == nil) {
		return nil, 0, errs.PermissionDenied
	}

	maxHistory := articleMaxHistory

	if !ctx.Organization.Expert {
		offset = 0
		maxHistory = articleutil.ArticleMaxHistoryNonExpert
	}

	langId = util.DetermineLang(nil, ctx.Organization.ID, ctx.UserId, langId).ID
	results := model.FindArticleContentVersionCommitByOrganizationIdAndArticleIdAndLanguageIdAndNotWIPLimit(ctx.Organization.ID, articleId, langId, offset, maxHistory)
	count := model.CountArticleContentVersionByArticleIdAndLanguageIdAndNotWIP(articleId, langId)
	articleutil.RemoveNonPublicInformation(ctx, results)
	return results, count, nil
}
