package search

import (
	"emviwiki/shared/model"
	"strings"
	"sync"
)

// Performs a fuzzy search on username, firstname and lastname for given organisaton.
func SearchUser(organization *model.Organization, query string, filter *model.SearchUserFilter) ([]model.User, int) {
	query = strings.TrimSpace(query)

	// search for all fields when no filter was passed
	if filter == nil {
		filter = new(model.SearchUserFilter)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var results []model.User
	var resultCount int

	go func() {
		results = model.FindUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmail(organization.ID, query, filter)
		wg.Done()
	}()

	go func() {
		resultCount = model.CountUserByOrganizationIdAndUsernameOrFirstnameOrLastnameOrEmail(organization.ID, query, filter)
		wg.Done()
	}()

	wg.Wait()
	return results, resultCount
}
