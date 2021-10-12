package pinned

import (
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
)

const (
	maxArticles = 4
	maxLists    = 4
)

func ReadPinned(ctx context.EmviContext, readArticles, readLists bool, offsetArticles, offsetLists int) ([]model.Article, []model.ArticleList, int, int) {
	langId := util.DetermineLang(nil, ctx.Organization.ID, ctx.UserId, 0).ID
	var articles []model.Article
	var lists []model.ArticleList
	articlesCount := 0
	listsCount := 0

	if readArticles && ctx.HasScopes(client.Scopes["articles"]) {
		articles = model.FindArticleByOrganizationIdAndUserIdAndLanguageIdAndClientAccessAndPinnedLimit(ctx.Organization.ID, ctx.UserId, langId, ctx.IsClient(), offsetArticles, maxArticles)
		articlesCount = model.CountArticleByOrganizationIdAndUserIdAndClientAccessAndPinned(ctx.Organization.ID, ctx.UserId, ctx.IsClient())
		articleutil.RemoveAuthorsOrAuthorMails(ctx, articles)
	}

	if readLists && ctx.HasScopes(client.Scopes["lists"]) {
		lists = model.FindArticleListsByOrganizationIdAndUserIdAndLanguageIdAndPinnedLimit(ctx.Organization.ID, ctx.UserId, langId, offsetLists, maxLists)
		listsCount = model.CountArticleListsByOrganizationIdAndUserIdAndClientAccessAndPinned(ctx.Organization.ID, ctx.UserId, ctx.IsClient())
	}

	return articles, lists, articlesCount, listsCount
}
