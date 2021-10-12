package content

import (
	"net/http"
)

// SetContentDownloadHeader sets the relevant HTTP headers for file downloads
func SetContentDownloadHeader(w http.ResponseWriter, filename, mimeType, md5 string) {
	if mimeType != "" {
		w.Header().Add("Content-Type", mimeType)
		w.Header().Add("Content-Disposition", "inline; filename="+filename)
	} else {
		w.Header().Add("Content-Disposition", "attachment; filename="+filename)
	}

	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.Header().Add("X-Dns-Prefetch-Control", "off")
	w.Header().Add("X-Download-Options", "noopen")
	w.Header().Add("X-Xss-Protection", "1; mode=block")
	w.Header().Add("Cache-Control", "max-age=1200")
	w.Header().Add("ETag", md5)
}
