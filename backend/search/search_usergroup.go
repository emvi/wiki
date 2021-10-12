package search

import (
	"emviwiki/shared/model"
	"strings"
	"sync"
)

// Performs a fuzzy search on name and info for given organization.
func SearchUsergroup(organization *model.Organization, query string, filter *model.SearchUserGroupFilter) ([]model.UserGroup, int) {
	if filter != nil && !filter.FindGroups {
		return nil, 0
	}

	query = strings.TrimSpace(query)

	// search for all fields when no filter was passed
	if filter == nil {
		filter = new(model.SearchUserGroupFilter)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var results []model.UserGroup
	var resultCount int

	go func() {
		results = model.FindUserGroupsByOrganizationIdAndNameOrInfo(organization.ID, query, filter)
		wg.Done()
	}()

	go func() {
		resultCount = model.CountUserGroupsByOrganizationIdAndNameOrInfo(organization.ID, query, filter)
		wg.Done()
	}()

	wg.Wait()
	return results, resultCount
}
