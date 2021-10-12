package api

import (
	"emviwiki/backend/bookmark"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
)

func BookmarkHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		ArticleId     hide.ID `json:"article_id"`
		ArticleListId hide.ID `json:"article_list_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := bookmark.BookmarkObject(ctx.Organization, ctx.UserId, req.ArticleId, req.ArticleListId); err != nil {
		return []error{err}
	}

	return nil
}

func ReadBookmarksHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	offsetArticles, err := rest.GetIntParam(r, "offset_articles")

	if err != nil {
		return []error{err}
	}

	offsetLists, err := rest.GetIntParam(r, "offset_lists")

	if err != nil {
		return []error{err}
	}

	readArticles := rest.GetBoolParam(r, "articles")
	readLists := rest.GetBoolParam(r, "lists")
	articles, lists := bookmark.ReadBookmarks(ctx.Organization, ctx.UserId, readArticles, readLists, offsetArticles, offsetLists)
	rest.WriteResponse(w, struct {
		Articles []model.Bookmark `json:"articles"`
		Lists    []model.Bookmark `json:"lists"`
	}{articles, lists})
	return nil
}
