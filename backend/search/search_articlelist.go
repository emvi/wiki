package search

import (
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"strings"
	"sync"
)

// Performs a fuzzy search on article lists names and infos.
func SearchArticleList(ctx context.EmviContext, query string, filter *model.SearchArticleListFilter) ([]model.ArticleList, int) {
	query = strings.TrimSpace(query)

	// search for all fields when no filter was passed
	if filter == nil {
		filter = new(model.SearchArticleListFilter)
	}

	filter.ClientAccess = ctx.IsClient()
	langId := util.DetermineLang(nil, ctx.Organization.ID, ctx.UserId, 0).ID
	var wg sync.WaitGroup
	wg.Add(2)
	var results []model.ArticleList
	var resultCount int

	go func() {
		results = model.FindArticleListsByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfo(ctx.Organization.ID, ctx.UserId, langId, query, filter)
		wg.Done()
	}()

	go func() {
		resultCount = model.CountArticleListsByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfo(ctx.Organization.ID, ctx.UserId, query, filter)
		wg.Done()
	}()

	wg.Wait()
	return results, resultCount
}
