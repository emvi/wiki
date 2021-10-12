package content

import (
	"emviwiki/dashboard/model"
	sharedcontent "emviwiki/shared/content"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"syscall"
)

func ReadContentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]
	file := model.GetFileByFilename(filename)

	if file == nil {
		logbuch.Debug("Reading file by filename not found", logbuch.Fields{"filename": filename})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reader, err := store.Read(file.Filename)

	if err != nil {
		logbuch.Debug("Error reading file from store", logbuch.Fields{"err": err, "filename": file.Filename})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		if err := reader.Close(); err != nil {
			logbuch.Error("Error closing reader for content", logbuch.Fields{"err": err, "filename": filename})
		}
	}()

	sharedcontent.SetContentDownloadHeader(w, file.Filename, file.MimeType, file.MD5)

	// EPIPE ignore broken pipe errors
	if _, err := io.Copy(w, reader); err != nil && err != syscall.EPIPE {
		logbuch.Error("Error sending content to client", logbuch.Fields{"err": err, "filename": filename})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
