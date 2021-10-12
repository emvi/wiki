package content

import (
	"emviwiki/shared/content"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
)

const (
	DefaultMaxFileSize = 52428800 // 50 MB
	filenameLength     = 20
	maxExtLen          = 10
	timeDirFormat      = "20060102"
	gbToBytes          = int64(1073741824)
	defaultMimeType    = "application/octet-stream"
)

type File struct {
	Request           *http.Request
	Organization      *model.Organization
	UserId            hide.ID
	ArticleId         hide.ID
	LangId            hide.ID
	RoomId            string
	Data              io.Reader
	ContentTypeHeader string
	Filename          string
	Path              string
	RequiresImage     bool
}

func UploadFile(file *File) (string, error) {
	cancelUpload, err := readMultipart(file)

	if err != nil {
		return "", err
	}

	if file.RequiresImage && !content.IsImage(getMimeType(file.ContentTypeHeader)) {
		return "", errs.FileType
	}

	if file.Organization == nil {
		file.Organization = &model.Organization{}
	}

	if file.ArticleId != 0 {
		if err := checkArticleAccess(file.Organization.ID, file.UserId, file.ArticleId, true); err != nil {
			logbuch.Debug("Access denied", logbuch.Fields{"orga_id": file.Organization.ID, "user_id": file.UserId, "article_id": file.ArticleId, "room_id": file.RoomId, "filename": file.Filename})
			return "", err
		}
	}

	// organization and user pictures don't count into storage usage
	if file.ArticleId != 0 || file.RoomId != "" {
		if err := checkUploadLimitReached(file.Organization); err != nil {
			logbuch.Debug("Upload limit reached", logbuch.Fields{"orga_id": file.Organization.ID, "user_id": file.UserId, "article_id": file.ArticleId, "room_id": file.RoomId, "filename": file.Filename})
			return "", err
		}
	}

	dir := getUploadDir(file.Organization, file.Path)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	uniqueName := generateUniqueFilename() + shortenExt(ext)
	md5Hash, size, err := saveFileInStore(file, dir, uniqueName, cancelUpload)

	if err != nil {
		logbuch.Error("Error saving file in store", logbuch.Fields{"err": err})
		go cleanupAttachment(dir, uniqueName)
		return "", err
	}

	fileUniqueName := uniqueName
	uniqueName, err = saveFileInDatabase(file, dir, uniqueName, md5Hash, size)

	if err != nil {
		logbuch.Error("Error saving file in database", logbuch.Fields{"err": err})
		go cleanupAttachment(dir, fileUniqueName)
		return "", err
	}

	return uniqueName, nil
}

func readMultipart(file *File) (<-chan bool, error) {
	pipeReader, pipeWriter := io.Pipe()
	cancelUpload := make(chan bool)
	filenameChan := make(chan string)
	contentTypeChan := make(chan string)
	reader, err := file.Request.MultipartReader()

	if err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if err := pipeWriter.Close(); err != nil {
				logbuch.Error("Error closing pipe writer on file upload", logbuch.Fields{"err": err})
			}
		}()
		readFilename := false

		for {
			part, err := reader.NextPart()

			if err == io.EOF {
				// file upload complete
				break
			} else if err != nil {
				logbuch.Warn("Error reading next part for multipart upload", logbuch.Fields{"err": err})
				cancelUpload <- true
				break
			}

			if !readFilename {
				readFilename = true
				filenameChan <- part.FileName()
				contentTypeChan <- part.Header.Get("Content-Type")
			}

			if _, err := io.Copy(pipeWriter, part); err != nil {
				logbuch.Warn("Error copying next part to pipe for multipart upload", logbuch.Fields{"err": err})
				cancelUpload <- true
				break
			}
		}
	}()

	file.Data = pipeReader

	select {
	case filename := <-filenameChan:
		file.Filename = filename
		file.ContentTypeHeader = <-contentTypeChan
	case <-cancelUpload:
		return nil, errs.IO
	}

	return cancelUpload, nil
}

func saveFileInStore(file *File, dir, uniqueName string, cancelUpload <-chan bool) (string, int64, error) {
	logbuch.Debug("Saving file to upload in store...", logbuch.Fields{"orga_id": file.Organization.ID, "user_id": file.UserId, "article_id": file.ArticleId, "room_id": file.RoomId, "filename": file.Filename})
	storeStartTime := time.Now()
	doneChan := make(chan bool)
	errChan := make(chan error, 2)

	go func() {
		if err := store.Save(dir, uniqueName, file.Data); err != nil {
			logbuch.Error("Error saving file on upload", logbuch.Fields{"orga_id": file.Organization.ID, "user_id": file.UserId, "article_id": file.ArticleId, "room_id": file.RoomId, "filename": file.Filename})
			errChan <- err
			return
		}

		doneChan <- true
	}()

	select {
	case <-doneChan: // we're good
	case err := <-errChan:
		return "", 0, err
	case <-cancelUpload:
		return "", 0, errs.IO
	}

	path := filepath.Join(dir, uniqueName)
	info, err := store.Info(path)

	if err != nil {
		logbuch.Error("Error reading info from uploaded file", logbuch.Fields{"orga_id": file.Organization.ID, "user_id": file.UserId, "article_id": file.ArticleId, "room_id": file.RoomId, "path": path})
		return "", 0, err
	}

	took := time.Now().Sub(storeStartTime)
	logbuch.Debug("Saved file in store", logbuch.Fields{"size": info.Size, "took_ms": took.Milliseconds()})
	return info.MD5, info.Size, nil
}

func saveFileInDatabase(file *File, dir, uniqueName, md5Hash string, size int64) (string, error) {
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction while uploading file", logbuch.Fields{"err": err})
		return "", errs.TxBegin
	}

	entity, err := createOrUpdateFile(tx, file.Organization, file.UserId, file.ArticleId, file.LangId, file.RoomId, dir, file.ContentTypeHeader, file.Filename, uniqueName, md5Hash, size)

	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when uploading file", logbuch.Fields{"err": err})
		return "", errs.TxCommit
	}

	return entity.UniqueName, nil
}

func createOrUpdateFile(tx *sqlx.Tx, organization *model.Organization, userId, articleId, langId hide.ID, roomId, dir, contentTypeHeader, filename, uniqueName, md5Hash string, size int64) (*model.File, error) {
	entity := findExistingFile(tx, organization, articleId, roomId, md5Hash)

	if entity == nil {
		// create a new file for non existing uploads and files that do not belong to articles
		entity = &model.File{OrganizationId: organization.ID,
			UserId:       userId,
			ArticleId:    articleId,
			RoomId:       null.NewString(roomId, roomId != ""),
			LanguageId:   langId,
			OriginalName: getFilename(filename),
			UniqueName:   uniqueName,
			Path:         dir,
			Type:         strings.ToLower(filepath.Ext(filename)),
			MimeType:     getMimeType(contentTypeHeader),
			Size:         size,
			MD5:          md5Hash}
	} else {
		// create a copy of the existing file and set references
		entity.ID = 0
		entity.UserId = userId
		entity.ArticleId = articleId
		entity.RoomId = null.NewString(roomId, roomId != "")
		entity.LanguageId = langId

		// delete the file we just uploaded
		go cleanupAttachment(dir, uniqueName)
	}

	if err := model.SaveFile(tx, entity); err != nil {
		logbuch.Error("Error saving file while uploading file", logbuch.Fields{"orga_id": organization.ID, "user_id": userId, "article_id": articleId, "room_id": roomId, "filename": filename})
		return nil, errs.Saving
	}

	return entity, nil
}

func cleanupAttachment(dir, uniqueName string) {
	path := filepath.Join(dir, uniqueName)
	logbuch.Debug("Cleaning up file that failed to be saved in database correctly", logbuch.Fields{"path": path})

	if err := store.Delete(path); err != nil {
		logbuch.Error("Error deleting uploaded file after failed transaction commit", logbuch.Fields{"err": err, "path": path})
	}
}

func checkUploadLimitReached(orga *model.Organization) error {
	if model.GetFileStorageUsageByOrganizationId(orga.ID) > orga.MaxStorageGB*gbToBytes {
		return errs.MaxStorageReached
	}

	return nil
}

func findExistingFile(tx *sqlx.Tx, orga *model.Organization, articleId hide.ID, roomId, md5Hash string) *model.File {
	var entity *model.File

	if orga != nil && (articleId != 0 || roomId != "") {
		entity = model.GetFileByOrganizationIdAndMD5AndArticleIdOrRoomIdNotNullTx(tx, orga.ID, md5Hash)
	}

	return entity
}

func getUploadDir(orga *model.Organization, path string) string {
	now := time.Now().Format(timeDirFormat)
	orgaName := ""

	if orga.NameNormalized != "" {
		orgaName = orga.NameNormalized + "_" + strconv.Itoa(int(orga.ID))
	}

	return filepath.Join(path, orgaName, now)
}

func getFilename(filename string) string {
	if len(filename) > 255 {
		return filename[:254]
	}

	return filename
}

func generateUniqueFilename() string {
	name := util.GenRandomString(filenameLength)

	for model.CountFileByUniqueNameTx(nil, name) != 0 {
		name = util.GenRandomString(filenameLength)
	}

	return name
}

func getMimeType(contentTypeHeader string) string {
	mimeType, _, err := mime.ParseMediaType(contentTypeHeader)

	if err != nil {
		logbuch.Warn("Error determining mime type from upload file header", logbuch.Fields{"err": err, "content_type": contentTypeHeader})
		return defaultMimeType
	}

	return mimeType
}

func shortenExt(ext string) string {
	if len(ext) > maxExtLen {
		return ext[:maxExtLen-1]
	}

	return ext
}
