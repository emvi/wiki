package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/pinned"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
)

func PinHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		ArticleId     hide.ID `json:"article_id"`
		ArticleListId hide.ID `json:"article_list_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := pinned.PinObject(ctx.Organization, ctx.UserId, req.ArticleId, req.ArticleListId); err != nil {
		return []error{err}
	}

	return nil
}

func ReadPinnedHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
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
	articles, lists, articlesCount, listsCount := pinned.ReadPinned(ctx, readArticles, readLists, offsetArticles, offsetLists)
	rest.WriteResponse(w, struct {
		Articles      []model.Article     `json:"articles"`
		Lists         []model.ArticleList `json:"lists"`
		ArticlesCount int                 `json:"articles_count"`
		ListsCount    int                 `json:"lists_count"`
	}{articles, lists, articlesCount, listsCount})
	return nil
}
