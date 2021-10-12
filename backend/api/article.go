package api

import (
	"emviwiki/backend/article"
	"emviwiki/backend/article/history"
	"emviwiki/backend/article/util"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"io"
	"net/http"
)

func SaveArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	// use pointer here, this can get quite large
	var req article.SaveArticleData

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	req.Organization = ctx.Organization
	id, err := article.SaveArticle(req)

	if err != nil {
		return err
	}

	rest.WriteResponse(w, struct {
		Id hide.ID `json:"id"`
	}{id})
	return nil
}

func ReadArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "lang") // default will be used if not set

	if err != nil {
		return []error{err}
	}

	version, err := rest.GetIntParam(r, "version")

	if err != nil {
		return []error{err}
	}

	userId := ctx.UserId

	if ctx.IsClient() {
		userId, err = rest.GetIdParam(r, "user_id")

		if err != nil {
			return []error{err}
		}

		ctx.UserId = userId
	}

	rawContent := rest.GetBoolParam(r, "raw_content")
	format := rest.GetParam(r, "format")
	result, err := article.ReadArticle(ctx, articleId, langId, version, !rawContent, format)

	if err != nil {
		return []error{err}
	}

	for i := range result.Article.Access {
		if result.Article.Access[i].User != nil && result.Article.Access[i].User.Picture.Valid {
			result.Article.Access[i].User.Picture.SetValid(getResourceURL(result.Article.Access[i].User.Picture.String))
		}
	}

	for i := range result.Authors {
		if result.Authors[i].Picture.Valid {
			result.Authors[i].Picture.SetValid(getResourceURL(result.Authors[i].Picture.String))
		}
	}

	for i := range result.Recommendations {
		if result.Recommendations[i].Member.User.Picture.Valid {
			result.Recommendations[i].Member.User.Picture.SetValid(getResourceURL(result.Recommendations[i].Member.User.Picture.String))
		}
	}

	if ctx.IsClient() {
		rest.WriteResponse(w, struct {
			Article *model.Article        `json:"article"`
			Content *model.ArticleContent `json:"content"`
			Authors []model.User          `json:"authors"`
		}{result.Article, result.Content, result.Authors})
	} else {
		rest.WriteResponse(w, struct {
			Article         *model.Article                `json:"article"`
			Content         *model.ArticleContent         `json:"content"`
			Authors         []model.User                  `json:"authors"`
			Write           bool                          `json:"write"`
			Observed        bool                          `json:"observed"`
			Bookmarked      bool                          `json:"bookmarked"`
			Recommendations []model.ArticleRecommendation `json:"recommendations"`
		}{
			result.Article,
			result.Content,
			result.Authors,
			result.WriteAccess,
			result.IsObserved,
			result.IsBookmarked,
			result.Recommendations,
		})
	}
	return nil
}

func ReadArticleHistoryHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "lang") // default will be used if not set

	if err != nil {
		return []error{err}
	}

	offset, err := rest.GetIntParam(r, "offset")

	if err != nil {
		return []error{err}
	}

	h, count, err := history.ReadArticleHistory(ctx, articleId, langId, offset)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		History []model.ArticleContent `json:"history"`
		Count   int                    `json:"count"`
	}{h, count})
	return nil
}

func DeleteArticleHistoryEntryHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	contentId, err := rest.GetIdParam(r, "content_id")

	if err != nil {
		return []error{err}
	}

	if err := history.DeleteHistoryEntry(ctx, contentId); err != nil {
		return []error{err}
	}

	return nil
}

func ReadArticlePreviewHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "lang") // default will be used if not set

	if err != nil {
		return []error{err}
	}

	content, err := article.GetArticlePreview(ctx, articleId, langId, rest.GetBoolParam(r, "preview_paragraph"))

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Content string `json:"content"`
	}{content})
	return nil
}

func RecommendArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		User                    []hide.ID `json:"user"`
		Groups                  []hide.ID `json:"groups"`
		Message                 string    `json:"message"`
		ReceiveReadConfirmation bool      `json:"receive_read_confirmation"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := article.RecommendArticle(ctx.Organization, ctx.UserId, articleId, req.User, req.Groups, req.Message, req.ReceiveReadConfirmation, mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func ConfirmRecommendationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		NotifyUserIds []hide.ID `json:"notify_user_ids"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := article.ConfirmRecommendation(ctx.Organization, ctx.UserId, articleId, req.NotifyUserIds); err != nil {
		return []error{err}
	}

	return nil
}

func InviteEditArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		LangId  hide.ID   `json:"language_id"`
		User    []hide.ID `json:"user"`
		Groups  []hide.ID `json:"groups"`
		RoomId  string    `json:"room_id"`
		Message string    `json:"message"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := article.InviteEditArticle(ctx.Organization, ctx.UserId, articleId, req.LangId, req.RoomId, req.Message, req.User, req.Groups, mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func ArchiveArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		Message string `json:"message"`
		Delete  bool   `json:"delete"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := article.ArchiveArticle(ctx.Organization, ctx.UserId, articleId, req.Message, req.Delete); err != nil {
		return []error{err}
	}

	return nil
}

func DeleteArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := article.DeleteArticle(ctx.Organization, ctx.UserId, id); err != nil {
		return []error{err}
	}

	return nil
}

func ResetArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		LangId  hide.ID `json:"language_id"`
		Version int     `json:"version"`
		Commit  string  `json:"commit"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := history.ResetArticle(ctx.Organization, ctx.UserId, articleId, req.LangId, req.Version, req.Commit); err != nil {
		return []error{err}
	}

	return nil
}

func CopyArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		LanguageId hide.ID `json:"language_id"` // optional
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	newArticleId, err := article.CopyArticle(ctx.Organization, ctx.UserId, articleId, req.LanguageId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		ArticleId hide.ID `json:"article_id"`
	}{newArticleId})
	return nil
}

func ReadPrivateArticlesHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	offset, err := rest.GetIntParam(r, "offset")

	if err != nil {
		return []error{err}
	}

	articles := article.ReadPrivateArticles(ctx.Organization, ctx.UserId, offset)
	rest.WriteResponse(w, articles)
	return nil
}

func ReadDraftsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	offset, err := rest.GetIntParam(r, "offset")

	if err != nil {
		return []error{err}
	}

	articles := article.ReadDrafts(ctx.Organization, ctx.UserId, offset)
	rest.WriteResponse(w, articles)
	return nil
}

func AddArticleToListHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		ArticleListIds []hide.ID `json:"article_list_ids"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := article.AddArticleToLists(ctx.Organization, ctx.UserId, articleId, req.ArticleListIds); err != nil {
		return []error{err}
	}

	return nil
}

func GetLinkMetaDataHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	url := rest.GetParam(r, "url")
	metaData := util.GetLinkMetaData(url)
	rest.WriteResponse(w, metaData)
	return nil
}

func ExportArticleHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "language_id")

	if err != nil {
		return []error{err}
	}

	format := rest.GetParam(r, "format")
	exportAttachments := rest.GetBoolParam(r, "export_attachments")
	reader, err := article.ExportArticle(ctx, articleId, langId, format, exportAttachments)

	if err != nil {
		return []error{err}
	}

	if _, err := io.Copy(w, reader); err != nil {
		logbuch.Error("Error sending zip to client", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}

	return nil
}
