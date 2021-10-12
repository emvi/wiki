package search

import (
	"emviwiki/backend/context"
	"strings"
	"sync"

	"emviwiki/shared/model"
)

// Performs a fuzzy search for a tag.
func SearchTag(ctx context.EmviContext, query string, filter *model.SearchTagFilter) ([]model.Tag, int) {
	query = strings.TrimSpace(query)

	// search for all fields when no filter was passed
	if filter == nil {
		filter = new(model.SearchTagFilter)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var results []model.Tag
	var resultCount int

	go func() {
		results = model.FindTagByOrganizationIdAndTagLimit(ctx.Organization.ID, ctx.UserId, query, filter)
		wg.Done()
	}()

	go func() {
		resultCount = model.CountTagByOrganizationIdAndTagLimit(ctx.Organization.ID, ctx.UserId, query, filter)
		wg.Done()
	}()

	wg.Wait()
	return results, resultCount
}
