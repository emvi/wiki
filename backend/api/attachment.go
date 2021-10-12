package api

import (
	"emviwiki/backend/article"
	"emviwiki/backend/content"
	"emviwiki/backend/context"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"net/http"
	"time"
)

func UploadArticleAttachmentHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	startTime := time.Now()
	r.Body = http.MaxBytesReader(w, r.Body, content.DefaultMaxFileSize)
	articleId, langId, roomId, err := getArticleAttachmentMetaData(r)

	if err != nil {
		return []error{err}
	}

	uniqueName, err := article.UploadAttachment(r, ctx.Organization, ctx.UserId, articleId, langId, roomId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		UniqueName string `json:"unique_name"`
	}{uniqueName})
	took := time.Now().Sub(startTime)
	logbuch.Debug("Finished attachment upload", logbuch.Fields{"unique_name": uniqueName, "took_ms": took.Milliseconds()})
	return nil
}

func getArticleAttachmentMetaData(r *http.Request) (hide.ID, hide.ID, string, error) {
	articleId, err := rest.GetIdParam(r, "article")

	if err != nil {
		return 0, 0, "", err
	}

	langId, err := rest.GetIdParam(r, "language")

	if err != nil {
		return 0, 0, "", err
	}

	roomId := rest.GetParam(r, "room")
	return articleId, langId, roomId, nil
}
