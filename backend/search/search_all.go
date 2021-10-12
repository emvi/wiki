package search

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"strings"
	"sync"
)

const (
	searchAllLimitArticles = 5
	searchAllLimitLists    = 3
	searchAllLimitUser     = 2
	searchAllLimitGroups   = 3
	searchAllLimitTags     = 2
)

type SearchAllFilter struct {
	Articles      bool `json:"articles"`
	Lists         bool `json:"lists"`
	User          bool `json:"user"`
	Groups        bool `json:"groups"`
	Tags          bool `json:"tags"`
	ArticlesLimit int  `json:"articles_limit"`
	ListsLimit    int  `json:"lists_limit"`
	UserLimit     int  `json:"user_limit"`
	GroupsLimit   int  `json:"groups_limit"`
	TagsLimit     int  `json:"tags_limit"`
}

type SearchAllResults struct {
	Articles     []model.Article     `json:"articles"`
	Lists        []model.ArticleList `json:"lists"`
	User         []model.User        `json:"user"`
	Groups       []model.UserGroup   `json:"groups"`
	Tags         []model.Tag         `json:"tags"`
	ArticleCount int                 `json:"articles_count"`
	ListCount    int                 `json:"lists_count"`
	UserCount    int                 `json:"user_count"`
	GroupCount   int                 `json:"groups_count"`
	TagCount     int                 `json:"tags_count"`
}

func SearchAll(ctx context.EmviContext, query string, filter *SearchAllFilter) *SearchAllResults {
	query = strings.TrimSpace(query)

	if query == "" {
		return &SearchAllResults{}
	}

	// search for all fields when no filter was passed
	if filter == nil {
		filter = &SearchAllFilter{Articles: true,
			Lists:  true,
			User:   true,
			Groups: true,
			Tags:   true}
	}

	// modify filter if call was made by client instead of a user
	if ctx.IsClient() {
		filter.Articles = filter.Articles && ctx.HasScopes(client.Scopes["articles"], client.Scopes["search_articles"])
		filter.Lists = filter.Lists && ctx.HasScopes(client.Scopes["lists"], client.Scopes["search_lists"])
		filter.Tags = filter.Tags && ctx.HasScopes(client.Scopes["tags"], client.Scopes["search_tags"])
		filter.User = false
		filter.Groups = false
	}

	// start concurrent search
	var articles []model.Article
	var lists []model.ArticleList
	var user []model.User
	var groups []model.UserGroup
	var tags []model.Tag
	var articleCount, listCount, userCount, groupCount, tagCount int
	var wg sync.WaitGroup

	if filter.Articles {
		wg.Add(1)
		go func() {
			articles, articleCount = SearchArticle(ctx, query, nil)
			wg.Done()
		}()
	}

	if filter.Lists {
		wg.Add(1)
		go func() {
			lists, listCount = SearchArticleList(ctx, query, nil)
			wg.Done()
		}()
	}

	if filter.User {
		wg.Add(1)
		go func() {
			user, userCount = SearchUser(ctx.Organization, query, nil)
			wg.Done()
		}()
	}

	if filter.Groups {
		wg.Add(1)
		go func() {
			groups, groupCount = SearchUsergroup(ctx.Organization, query, nil)
			wg.Done()
		}()
	}

	if filter.Tags {
		wg.Add(1)
		go func() {
			tags, tagCount = SearchTag(ctx, query, nil)
			wg.Done()
		}()
	}

	wg.Wait()
	return &SearchAllResults{
		articles[:searchAllLimit(searchAllLimitArticles, filter.ArticlesLimit, len(articles))],
		lists[:searchAllLimit(searchAllLimitLists, filter.ListsLimit, len(lists))],
		user[:searchAllLimit(searchAllLimitUser, filter.UserLimit, len(user))],
		groups[:searchAllLimit(searchAllLimitGroups, filter.GroupsLimit, len(groups))],
		tags[:searchAllLimit(searchAllLimitTags, filter.TagsLimit, len(tags))],
		articleCount,
		listCount,
		userCount,
		groupCount,
		tagCount,
	}
}

func searchAllLimit(limit, filterLimit, n int) int {
	if filterLimit > limit {
		limit = filterLimit
	}

	if limit < 0 {
		return 0
	} else if limit > n {
		return n
	}

	return limit
}
