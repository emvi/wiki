package api

import (
	"emviwiki/backend/content"
	sharedcontent "emviwiki/shared/content"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strings"
	"syscall"
)

// TODO signed URL
func GetContentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := strings.TrimSpace(vars["filename"])
	file, reader, err := content.ReadFile(filename)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		if err := reader.Close(); err != nil {
			logbuch.Error("Error closing reader for content", logbuch.Fields{"err": err, "filename": filename})
		}
	}()

	sharedcontent.SetContentDownloadHeader(w, file.OriginalName, file.MimeType, file.MD5)

	// EPIPE ignore broken pipe errors
	if _, err := io.Copy(w, reader); err != nil && err != syscall.EPIPE {
		logbuch.Error("Error sending content to client", logbuch.Fields{"err": err, "filename": filename})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
