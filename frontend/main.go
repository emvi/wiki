package main

import (
	"bytes"
	"emviwiki/shared/config"
	"emviwiki/shared/server"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	envJsPrefix     = "/dist/env.js"
	buildJsFile     = "public/dist/build.js"
	buildJsGzipFile = "public/dist/build.js.gz"
	buildJsPrefix   = "/dist/build.js"
	indexFile       = "public/index.html"
	rootDirPrefix   = "/"
	envVars         = `var EMVI_WIKI_AUTH_HOST = '%s';
var EMVI_WIKI_WEBSITE_HOST = '%s';
var EMVI_WIKI_BACKEND_HOST = '%s';
var EMVI_WIKI_COLLAB_HOST = '%s';
var EMVI_WIKI_AUTH_CLIENT_ID = '%s';
var EMVI_WIKI_INTEGRATION = %t;
var EMVI_WIKI_STRIPE_PUBLIC_KEY = '%s';`
)

var (
	envJs, buildJs, buildJsGzip, indexHtml []byte
	watchBuildJs, watchIndexHtml           bool
	version                                string
)

func loadConfig() {
	c := config.Get()
	watchBuildJs = c.Dev.WatchBuildJs
	watchIndexHtml = c.Dev.WatchIndexHtml
	version = c.Version
}

func loadEnvJs() {
	logbuch.Info("Loading env.js...")
	c := config.Get()
	envJs = []byte(fmt.Sprintf(envVars,
		c.Hosts.Auth,
		c.Hosts.Website,
		c.Hosts.Backend,
		c.Hosts.Collab,
		c.AuthClient.ID,
		c.IsIntegration,
		c.Stripe.PublicKey))
}

func loadBuildJs() {
	logbuch.Info("Loading build.js...")
	var err error
	buildJs, err = ioutil.ReadFile(buildJsFile)

	if err != nil {
		logbuch.Fatal("build.js not found", logbuch.Fields{"err": err})
	}

	buildJsGzip, err = ioutil.ReadFile(buildJsGzipFile)

	if err != nil {
		logbuch.Fatal("build.js.gz not found", logbuch.Fields{"err": err})
	}
}

func loadIndexHtml() {
	logbuch.Info("Loading index.html...")
	tpl, err := template.ParseFiles(indexFile)

	if err != nil {
		logbuch.Fatal("index.html not found", logbuch.Fields{"err": err})
		return
	}

	data := struct {
		Version string
	}{version}
	var buffer bytes.Buffer

	if err := tpl.Execute(&buffer, data); err != nil {
		logbuch.Fatal("Error executing index.html template", logbuch.Fields{"err": err})
	}

	indexHtml = buffer.Bytes()
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	server.ServeStaticFiles(router, "", "", version)
	server.ResourceHandler(router, envJsPrefix, func(w http.ResponseWriter, r *http.Request) {
		if watchBuildJs {
			loadEnvJs()
		}

		server.AddContentHeaders(w, "text/javascript", version)

		if _, err := w.Write(envJs); err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	server.ResourceHandler(router, buildJsPrefix, func(w http.ResponseWriter, r *http.Request) {
		if watchBuildJs {
			loadBuildJs()
		}

		server.AddContentHeaders(w, "text/javascript", version)

		if strings.Contains(strings.ToLower(r.Header.Get("Accept-Encoding")), "gzip") {
			w.Header().Add("Content-Encoding", "gzip")

			if _, err := w.Write(buildJsGzip); err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			if _, err := w.Write(buildJs); err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	})
	router.HandleFunc("/health", server.SecurityHeadersMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	server.ResourceHandler(router, rootDirPrefix, func(w http.ResponseWriter, r *http.Request) {
		if watchIndexHtml {
			loadIndexHtml()
		}

		server.AddContentHeaders(w, "text/html", version)

		if _, err := w.Write(indexHtml); err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	return router
}

func main() {
	config.Load()
	loadConfig()
	stdout, stderr := server.ConfigureLogging()
	defer server.CloseLogger(stdout, stderr)
	loadEnvJs()
	loadBuildJs()
	loadIndexHtml()
	router := setupRouter()
	cors := server.ConfigureCors(router)
	server.Start(cors, nil)
}
