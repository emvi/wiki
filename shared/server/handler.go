package server

import (
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"net/http"
)

// ResourceHandler adds a handler for given prefix to the router. The prefix will be stripped and the response gzipped.
func ResourceHandler(router *mux.Router, prefix string, handler func(http.ResponseWriter, *http.Request)) {
	router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, gziphandler.GzipHandler(http.HandlerFunc(handler))))
}
