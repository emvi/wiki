package server

import (
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	defaultStaticDir       = "public/static"
	defaultStaticDirPrefix = "/static/"
)

// ServeStaticFiles adds a static file handler for given static content directory and prefix.
// It adds the Cache-Control header (3600s) to each file returned.
// The staticDirPrefix and staticDir are optional, if left empty, defaults will be used.
// The version will be used as ETag header.
func ServeStaticFiles(router *mux.Router, staticDirPrefix, staticDir, version string) {
	if staticDirPrefix == "" {
		staticDirPrefix = defaultStaticDirPrefix
	}

	if staticDir == "" {
		staticDir = defaultStaticDir
	}

	fs := http.StripPrefix(staticDirPrefix, http.FileServer(http.Dir(staticDir)))
	router.PathPrefix(staticDirPrefix).Handler(gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AddContentHeaders(w, "", version)
		fs.ServeHTTP(w, r)
	})))
}

// AddContentHeaders adds security, Content-Type, Cache-Control and ETag header to the response.
// The ContentType header will only be set if contentType is non empty.
// The max age will be set to 31536000 (1 year) and the ETag to given version.
func AddContentHeaders(w http.ResponseWriter, contentType, version string) {
	if contentType != "" {
		w.Header().Add("Content-Type", contentType)
	}

	w.Header().Add("Cache-Control", "max-age=31536000")
	w.Header().Add("ETag", version)
	AddSecurityHeaders(w)
}
