package pages

import (
	"emviwiki/dashboard/auth"
	"emviwiki/dashboard/model"
	"emviwiki/shared/util"
	"errors"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"path/filepath"
)

const (
	maxFileMem  = 1024 * 1024 * 100 // 100 MB
	formFile    = "file"
	filenameLen = 10
)

func FilesPageHandler(claims *auth.UserTokenClaims, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := uploadFile(r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/files", http.StatusFound)
	} else if r.Method == http.MethodDelete {
		if err := deleteFile(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	data := struct {
		Files    []model.File
		MailHost string
	}{
		model.FindFile(),
		mailHost,
	}

	RenderPage(w, filesPageTemplate, claims, &data)
}

func uploadFile(r *http.Request) error {
	if err := r.ParseMultipartForm(maxFileMem); err != nil {
		logbuch.Error("Error parsing form on file upload", logbuch.Fields{"err": err})
		return err
	}

	file, head, err := r.FormFile(formFile)

	if err != nil {
		logbuch.Error("Error reading form file", logbuch.Fields{"err": err})
		return err
	}

	mimeType, _, err := mime.ParseMediaType(head.Header.Get("Content-Type"))

	if err != nil {
		logbuch.Error("Error reading mime type from file", logbuch.Fields{"err": err})
		return err
	}

	filename := getUniqueFilename() + filepath.Ext(head.Filename)

	if err := store.Save("", filename, file); err != nil {
		logbuch.Error("Error saving file in store", logbuch.Fields{"err": err, "filename": filename})
		return err
	}

	info, err := store.Info(filename)

	if err != nil {
		logbuch.Error("Error reading file info from store", logbuch.Fields{"err": err, "filename": filename})
		return err
	}

	f := &model.File{Filename: filename,
		OriginalFilename: head.Filename,
		MimeType:         mimeType,
		Size:             info.Size,
		MD5:              info.MD5}

	if err := model.SaveFile(nil, f); err != nil {
		logbuch.Error("Error saving file on upload", logbuch.Fields{"err": err, "filename": filename})
		return err
	}

	return nil
}

func getUniqueFilename() string {
	filename := util.GenRandomString(filenameLen)

	for model.GetFileByFilename(filename) != nil {
		filename = util.GenRandomString(filenameLen)
	}

	return filename
}

func deleteFile(r *http.Request) error {
	params := mux.Vars(r)
	id, err := hide.FromString(params["id"])

	if err != nil {
		return err
	}

	file := model.GetFileById(id)

	if file == nil {
		return errors.New("file not found")
	}

	if err := store.Delete(file.Filename); err != nil {
		logbuch.Error("Error deleting file in store", logbuch.Fields{"err": err, "id": file.ID, "filename": file.Filename})
	}

	if err := model.DeleteFileById(nil, id); err != nil {
		return err
	}

	return nil
}
