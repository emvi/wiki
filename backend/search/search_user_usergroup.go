package search

import (
	"emviwiki/shared/model"
)

type UserUsergroupSearchResult struct {
	User      *model.User      `json:"user"`
	UserGroup *model.UserGroup `json:"usergroup"`
}

func SearchUserUsergroup(orga *model.Organization, query string) []UserUsergroupSearchResult {
	user, _ := SearchUser(orga, query, nil)
	groups, _ := SearchUsergroup(orga, query, &model.SearchUserGroupFilter{FindGroups: orga.Expert})
	results := make([]UserUsergroupSearchResult, len(user)+len(groups))

	for i := 0; i < len(user); i++ {
		results[i] = UserUsergroupSearchResult{User: &user[i]}
	}

	n := len(user)

	for i := 0; i < len(groups); i++ {
		results[n+i] = UserUsergroupSearchResult{UserGroup: &groups[i]}
	}

	return results
}
