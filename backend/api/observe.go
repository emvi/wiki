package api

import (
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"net/http"

	"emviwiki/backend/observe"
	"emviwiki/shared/rest"
)

func ObserveObjectHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		ArticleId     hide.ID `json:"article_id"`
		ArticleListId hide.ID `json:"article_list_id"`
		UserGroupId   hide.ID `json:"user_group_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := observe.ObserveObject(ctx.Organization, ctx.UserId, req.ArticleId, req.ArticleListId, req.UserGroupId); err != nil {
		return []error{err}
	}

	return nil
}

func ReadObservedHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	readArticles := rest.GetBoolParam(r, "articles")
	readLists := rest.GetBoolParam(r, "lists")
	readGroups := rest.GetBoolParam(r, "groups")
	offsetArticles, err := rest.GetIntParam(r, "offset_articles")

	if err != nil {
		return []error{err}
	}

	offsetLists, err := rest.GetIntParam(r, "offset_lists")

	if err != nil {
		return []error{err}
	}

	offsetGroups, err := rest.GetIntParam(r, "offset_groups")

	if err != nil {
		return []error{err}
	}

	articles, lists, groups := observe.ReadObserved(ctx.Organization, ctx.UserId, readArticles, readLists, readGroups, offsetArticles, offsetLists, offsetGroups)
	rest.WriteResponse(w, struct {
		Articles []model.Article     `json:"articles"`
		Lists    []model.ArticleList `json:"lists"`
		Groups   []model.UserGroup   `json:"groups"`
	}{articles, lists, groups})
	return nil
}
