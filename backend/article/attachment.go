package article

import (
	"emviwiki/backend/content"
	"emviwiki/backend/errs"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

const (
	attachmentPath   = "organization/attachments"
	fileNodeTypeImg  = "image"
	fileNodeTypeFile = "file"
	fileNodeTypePDF  = "pdf"
	fileNodeImgAttr  = "src"
	fileNodeFileAttr = "file"
	fileNodePDFAttr  = "src"
)

// UploadAttachment uploads a file for given article or room.
// The room is used to temporarily upload a file to an unsaved article.
func UploadAttachment(r *http.Request, orga *model.Organization, userId, articleId, langId hide.ID, roomId string) (string, error) {
	if articleId == 0 && roomId == "" {
		return "", errs.ArticleNotFound
	}

	uniqueName, err := content.UploadFile(&content.File{
		Request:      r,
		Organization: orga,
		UserId:       userId,
		ArticleId:    articleId,
		LangId:       langId,
		RoomId:       roomId,
		Path:         attachmentPath,
	})

	if err == errs.ArticleNotFound || err == errs.PermissionDenied {
		logbuch.Error("Error on article attachment upload", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId, "room_id": roomId})
		return "", err
	} else if err != nil {
		logbuch.Error("Unexpected error on article attachment upload", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId, "room_id": roomId})
		return "", errs.UploadingFile
	}

	return uniqueName, nil
}

// Updates the article ID of all files belonging to an article.
func updateAttachments(tx *sqlx.Tx, orgaId, articleId, langId hide.ID, roomId string) error {
	if roomId != "" {
		attachments := model.FindFileByOrganizationIdAndRoomIdTx(tx, orgaId, roomId)

		for _, file := range attachments {
			file.ArticleId = articleId
			file.RoomId.Valid = false
			file.LanguageId = langId

			if err := model.SaveFile(tx, &file); err != nil {
				return errs.Saving
			}
		}
	}

	return nil
}

// Deletes all files belonging to an article that are not used within the new document content anymore.
// This does not affect old content versions.
func cleanupAttachments(orgaId, userId, articleId hide.ID, lastContentDefTime time.Time, newContent *model.ArticleContent) error {
	doc, err := prosemirror.ParseDoc(newContent.Content)

	if err != nil {
		return err
	}

	attachments := model.FindFileByOrganizationIdAndArticleIdAndLanguageIdAndDefTimeAfter(orgaId, articleId, newContent.LanguageId, lastContentDefTime)
	files := prosemirror.FindNodes(doc, -1, fileNodeTypeImg, fileNodeTypeFile, fileNodeTypePDF)

	for _, file := range attachments {
		if !attachmentExistsInContent(file.UniqueName, files) {
			if len(model.FindFileByOrganizationIdAndUniqueNameAndNotId(orgaId, file.UniqueName, file.ID)) == 0 {
				path := filepath.Join(file.Path, file.UniqueName)
				logbuch.Debug("Deleting unused attachment in store", logbuch.Fields{"article_id": articleId, "path": path})
				content.DeleteFileInStore(orgaId, userId, path)
			}

			logbuch.Debug("Deleting unused file in database", logbuch.Fields{"article_id": articleId, "id": file.ID})

			if err := model.DeleteFileById(nil, file.ID); err != nil {
				logbuch.Error("Error deleting file when cleaning up article attachments", logbuch.Fields{"err": err, "article_id": articleId})
			}
		}
	}

	return nil
}

func attachmentExistsInContent(uniqueName string, fileNodes []prosemirror.Node) bool {
	for _, node := range fileNodes {
		if node.Type == fileNodeTypeImg {
			src, ok := node.Attrs[fileNodeImgAttr].(string)

			if ok && strings.Contains(src, uniqueName) {
				return true
			}
		} else if node.Type == fileNodeTypeFile {
			file, ok := node.Attrs[fileNodeFileAttr].(string)

			if ok && strings.Contains(file, uniqueName) {
				return true
			}
		} else {
			pdf, ok := node.Attrs[fileNodePDFAttr].(string)

			if ok && strings.Contains(pdf, uniqueName) {
				return true
			}
		}
	}

	return false
}
