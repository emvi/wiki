package api

import (
	"emviwiki/backend/context"
	"net/http"
	"time"

	"emviwiki/backend/search"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
)

func SearchUserHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")
	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchUserFilter{
		BaseSearch:    baseFilter,
		Username:      rest.GetParam(r, "username"),
		Email:         rest.GetParam(r, "email"),
		Firstname:     rest.GetParam(r, "firstname"),
		Lastname:      rest.GetParam(r, "lastname"),
		SortUsername:  rest.GetParam(r, "sort_username"),
		SortEmail:     rest.GetParam(r, "sort_email"),
		SortFirstname: rest.GetParam(r, "sort_firstname"),
		SortLastname:  rest.GetParam(r, "sort_lastname"),
	}

	user, count := search.SearchUser(ctx.Organization, query, filter)

	for i := range user {
		if user[i].Picture.Valid {
			user[i].Picture.SetValid(getResourceURL(user[i].Picture.String))
		}
	}

	rest.WriteResponse(w, struct {
		User  []model.User `json:"user"`
		Count int          `json:"count"`
	}{user, count})
	return nil
}

func SearchUsergroupHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")
	userIds, err := rest.GetIdParams(r, "user_ids")

	if err != nil {
		return []error{err}
	}

	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchUserGroupFilter{
		BaseSearch: baseFilter,
		Name:       rest.GetParam(r, "name"),
		Info:       rest.GetParam(r, "info"),
		UserIds:    userIds,
		SortName:   rest.GetParam(r, "sort_name"),
		SortInfo:   rest.GetParam(r, "sort_info"),
		FindGroups: rest.GetBoolParam(r, "find_groups"),
	}

	groups, count := search.SearchUsergroup(ctx.Organization, query, filter)
	rest.WriteResponse(w, struct {
		Groups []model.UserGroup `json:"groups"`
		Count  int               `json:"count"`
	}{groups, count})
	return nil
}

func SearchUserUsergroupHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")
	resp := search.SearchUserUsergroup(ctx.Organization, query)

	for i := range resp {
		if resp[i].User != nil && resp[i].User.Picture.Valid {
			resp[i].User.Picture.SetValid(getResourceURL(resp[i].User.Picture.String))
		}
	}

	rest.WriteResponse(w, resp)
	return nil
}

func SearchTagHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")
	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchTagFilter{
		BaseSearch: baseFilter,
		SortUsages: rest.GetParam(r, "sort_usages"),
		SortName:   rest.GetParam(r, "sort_name"),
	}

	tags, count := search.SearchTag(ctx, query, filter)
	rest.WriteResponse(w, struct {
		Tags  []model.Tag `json:"tags"`
		Count int         `json:"count"`
	}{tags, count})
	return nil
}

func SearchArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")
	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "language_id")

	if err != nil {
		return []error{err}
	}

	tagIds, err := rest.GetIdParams(r, "tag_ids")

	if err != nil {
		return []error{err}
	}

	authorUserIds, err := rest.GetIdParams(r, "authors")

	if err != nil {
		return []error{err}
	}

	userGroupIds, err := rest.GetIdParams(r, "user_group_ids")

	if err != nil {
		return []error{err}
	}

	publishedStart, err := rest.GetDateParam(r, "published_start")

	if err != nil {
		publishedStart = time.Time{}
	}

	publishedEnd, err := rest.GetDateParam(r, "published_end")

	if err != nil {
		publishedEnd = time.Time{}
	}

	filter := &model.SearchArticleFilter{
		BaseSearch:       baseFilter,
		LanguageId:       langId,
		Archived:         rest.GetBoolParam(r, "archived"),
		WIP:              rest.GetBoolParam(r, "wip"),
		ClientAccess:     rest.GetBoolParam(r, "client_access"),
		Preview:          rest.GetBoolParam(r, "preview"),
		PreviewParagraph: rest.GetBoolParam(r, "preview_paragraph"),
		PreviewImage:     rest.GetBoolParam(r, "preview_image"),
		Title:            rest.GetParam(r, "title"),
		Content:          rest.GetParam(r, "content"),
		Tags:             rest.GetParam(r, "tags"),
		TagIds:           tagIds,
		AuthorUserIds:    authorUserIds,
		UserGroupIds:     userGroupIds,
		Commits:          rest.GetParam(r, "commits"),
		PublishedStart:   publishedStart,
		PublishedEnd:     publishedEnd,
		SortTitle:        rest.GetParam(r, "sort_title"),
		SortPublished:    rest.GetParam(r, "sort_published"),
		SortRelevance:    rest.GetParam(r, "sort_relevance"),
	}

	articles, count := search.SearchArticle(ctx, query, filter)

	if ctx.IsClient() {
		for i := range articles {
			if articles[i].LatestArticleContent != nil {
				articles[i].LatestArticleContent.TitleTsvector = ""
				articles[i].LatestArticleContent.ContentTsvector = ""
			}
		}
	}

	rest.WriteResponse(w, struct {
		Articles []model.Article `json:"articles"`
		Count    int             `json:"count"`
	}{articles, count})
	return nil
}

func SearchArticleListHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")
	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	userIds, err := rest.GetIdParams(r, "user_ids")

	if err != nil {
		return []error{err}
	}

	userGroupIds, err := rest.GetIdParams(r, "user_group_ids")

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchArticleListFilter{
		BaseSearch:   baseFilter,
		ClientAccess: rest.GetBoolParam(r, "client_access"),
		Name:         rest.GetParam(r, "name"),
		Info:         rest.GetParam(r, "info"),
		UserIds:      userIds,
		UserGroupIds: userGroupIds,
		SortName:     rest.GetParam(r, "sort_name"),
		SortInfo:     rest.GetParam(r, "sort_info"),
	}

	lists, count := search.SearchArticleList(ctx, query, filter)
	rest.WriteResponse(w, struct {
		Lists []model.ArticleList `json:"lists"`
		Count int                 `json:"count"`
	}{lists, count})
	return nil
}

func SearchAllHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	query := rest.GetParam(r, "query")

	articlesLimit, err := rest.GetIntParam(r, "articles_limit")

	if err != nil {
		return []error{err}
	}

	listsLimit, err := rest.GetIntParam(r, "lists_limit")

	if err != nil {
		return []error{err}
	}

	userLimit, err := rest.GetIntParam(r, "user_limit")

	if err != nil {
		return []error{err}
	}

	groupsLimit, err := rest.GetIntParam(r, "groups_limit")

	if err != nil {
		return []error{err}
	}

	tagsLimit, err := rest.GetIntParam(r, "tags_limit")

	if err != nil {
		return []error{err}
	}

	filter := &search.SearchAllFilter{
		Articles:      rest.GetBoolParam(r, "articles"),
		Lists:         rest.GetBoolParam(r, "lists"),
		User:          rest.GetBoolParam(r, "user"),
		Groups:        rest.GetBoolParam(r, "groups"),
		Tags:          rest.GetBoolParam(r, "tags"),
		ArticlesLimit: articlesLimit,
		ListsLimit:    listsLimit,
		UserLimit:     userLimit,
		GroupsLimit:   groupsLimit,
		TagsLimit:     tagsLimit,
	}

	results := search.SearchAll(ctx, query, filter)

	for i := range results.User {
		if results.User[i].Picture.Valid {
			results.User[i].Picture.SetValid(getResourceURL(results.User[i].Picture.String))
		}
	}

	rest.WriteResponse(w, results)
	return nil
}

func getBaseFilter(w http.ResponseWriter, r *http.Request) (model.BaseSearch, error) {
	createdStart, err := rest.GetDateParam(r, "created_start")

	if err != nil {
		return model.BaseSearch{}, err
	}

	createdEnd, err := rest.GetDateParam(r, "created_end")

	if err != nil {
		return model.BaseSearch{}, err
	}

	updatedStart, err := rest.GetDateParam(r, "updated_start")

	if err != nil {
		return model.BaseSearch{}, err
	}

	updatedEnd, err := rest.GetDateParam(r, "updated_end")

	if err != nil {
		return model.BaseSearch{}, err
	}

	sortCreated := rest.GetParam(r, "sort_created")
	sortUpdated := rest.GetParam(r, "sort_updated")
	offset, err := rest.GetIntParam(r, "offset")

	if err != nil {
		return model.BaseSearch{}, err
	}

	limit, err := rest.GetIntParam(r, "limit")

	if err != nil {
		return model.BaseSearch{}, err
	}

	return model.BaseSearch{CreatedStart: createdStart,
		CreatedEnd:   createdEnd,
		UpdatedStart: updatedStart,
		UpdatedEnd:   updatedEnd,
		SortCreated:  sortCreated,
		SortUpdated:  sortUpdated,
		Offset:       offset,
		Limit:        limit}, nil
}
