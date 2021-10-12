package main

import (
	"bytes"
	"emviwiki/shared/config"
	"emviwiki/shared/server"
	"emviwiki/shared/util"
	"emviwiki/website/legal"
	"emviwiki/website/pages"
	"emviwiki/website/sitemap"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	envJsPrefix            = "/dist/env.js"
	buildJsFile            = "public/dist/build.js"
	buildJsGzipFile        = "public/dist/build.js.gz"
	buildJsPrefix          = "/dist/build.js"
	robotsTxtFile          = "public/robots.txt"
	robotsTxtPrefix        = "/robots.txt"
	sitemapXmlPrefix       = "/sitemap.xml"
	microsoftAppJsonFile   = "public/microsoft-identity-association.json"
	microsoftAppJsonPrefix = "/.well-known/microsoft-identity-association.json"
	indexFile              = "public/index.html"
	rootDirPrefix          = "/"
	envVars                = `var EMVI_WIKI_WEBSITE_HOST = '%s';
var EMVI_WIKI_FRONTEND_HOST = '%s';
var EMVI_WIKI_AUTH_HOST = '%s';
var EMVI_WIKI_BACKEND_HOST = '%s';
var EMVI_WIKI_AUTH_CLIENT_ID = '%s';
var EMVI_WIKI_RECAPTCHA_CLIENT_SECRET = '%s';
var EMVI_WIKI_INTEGRATION = %t;
var EMVI_WIKI_GITHUB_CLIENT_ID = '%s';
var EMVI_WIKI_SLACK_CLIENT_ID = '%s';
var EMVI_WIKI_GOOGLE_CLIENT_ID = '%s';
var EMVI_WIKI_MICROSOFT_CLIENT_ID = '%s';`
)

var (
	envJs, buildJs, buildJsGzip, indexHtml, robotsTxt, sitemapXml, microsoftAppJson []byte
	watchBuildJs, watchIndexHtml, isIntegration                                     bool
	version, websiteHost                                                            string
	sitemapURLs                                                                     []sitemap.URL
)

func loadConfig() {
	c := config.Get()
	watchBuildJs = c.Dev.WatchBuildJs
	watchIndexHtml = c.Dev.WatchIndexHtml
	isIntegration = c.IsIntegration
	version = c.Version
	websiteHost = c.Hosts.Website
	sitemapURLs = []sitemap.URL{
		{Loc: websiteHost + "/", Priority: "1.0"},
		{Loc: websiteHost + "/pricing", Priority: "0.8"},
		{Loc: websiteHost + "/blog", Priority: "0.8"},
		{Loc: websiteHost + "/blog/", Priority: "0.8"},
		{Loc: websiteHost + "/registration", Priority: "0.5"},
		{Loc: websiteHost + "/terms", Priority: "0.1"},
		{Loc: websiteHost + "/privacy", Priority: "0.1"},
		{Loc: websiteHost + "/cookies", Priority: "0.1"},
		{Loc: websiteHost + "/legal", Priority: "0.1"},
	}
}

func loadEnvJs() {
	logbuch.Info("Loading env.js...")
	c := config.Get()
	envJs = []byte(fmt.Sprintf(envVars,
		c.Hosts.Website,
		c.Hosts.Frontend,
		c.Hosts.Auth,
		c.Hosts.Backend,
		c.AuthClient.ID,
		c.RecaptchaClientSecret,
		c.IsIntegration,
		c.SSO.GitHub.ID,
		c.SSO.Slack.ID,
		c.SSO.Google.ID,
		c.SSO.Microsoft.ID,
	))
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
		Version       string
		IsIntegration bool
	}{version, isIntegration}
	var buffer bytes.Buffer

	if err := tpl.Execute(&buffer, data); err != nil {
		logbuch.Fatal("Error executing index.html template", logbuch.Fields{"err": err})
	}

	indexHtml = buffer.Bytes()
}

func loadRobotsTxt() {
	logbuch.Info("Loading robots.txt...")
	content, err := ioutil.ReadFile(robotsTxtFile)

	if err != nil {
		logbuch.Fatal("robots.txt not found", logbuch.Fields{"err": err})
	}

	robotsTxt = content
}

func loadSitemapXml() {
	logbuch.Info("Loading sitemap.xml...")
	content, err := sitemap.GenerateSitemap(sitemapURLs)

	if err != nil {
		logbuch.Fatal("Error generating sitemap.xml", logbuch.Fields{"err": err})
	}

	sitemapXml = content
}

func loadMicrosoftApp() {
	logbuch.Info("Loading microsoft-identity-association.json...")
	content, err := ioutil.ReadFile(microsoftAppJsonFile)

	if err != nil {
		logbuch.Fatal("microsoft-identity-association.json not found", logbuch.Fields{"err": err})
	}

	microsoftAppJson = content
}

func billingRedirect(w http.ResponseWriter, r *http.Request, success bool) {
	organization := strings.TrimSpace(r.URL.Query().Get("organization"))

	if organization == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := "success"

	if !success {
		result = "cancel"
	}

	target := util.InjectSubdomain(config.Get().Hosts.Frontend, organization)
	target = fmt.Sprintf("%s/billing?%s=true", target, result)
	http.Redirect(w, r, target, http.StatusFound)
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
	server.ResourceHandler(router, robotsTxtPrefix, func(w http.ResponseWriter, r *http.Request) {
		server.AddContentHeaders(w, "text/plain", version)

		if isIntegration {
			if _, err := w.Write(robotsTxt); err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	server.ResourceHandler(router, sitemapXmlPrefix, func(w http.ResponseWriter, r *http.Request) {
		server.AddContentHeaders(w, "text/xml", version)

		if _, err := w.Write(sitemapXml); err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	server.ResourceHandler(router, microsoftAppJsonPrefix, func(w http.ResponseWriter, r *http.Request) {
		server.AddContentHeaders(w, "application/json", version)

		if _, err := w.Write(microsoftAppJson); err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	})
	router.HandleFunc("/billing/success", server.SecurityHeadersMiddleware(func(w http.ResponseWriter, r *http.Request) {
		billingRedirect(w, r, true)
	}))
	router.HandleFunc("/billing/cancel", server.SecurityHeadersMiddleware(func(w http.ResponseWriter, r *http.Request) {
		billingRedirect(w, r, false)
	}))
	router.HandleFunc("/health", server.SecurityHeadersMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	router.HandleFunc("/pricing", server.SecurityHeadersMiddleware(pages.PricingPageHandler))
	router.HandleFunc("/terms", server.SecurityHeadersMiddleware(pages.TermsPageHandler))
	router.HandleFunc("/privacy", server.SecurityHeadersMiddleware(pages.PrivacyPageHandler))
	router.HandleFunc("/cookies", server.SecurityHeadersMiddleware(pages.CookiePageHandler))
	router.HandleFunc("/legal", server.SecurityHeadersMiddleware(pages.LegalPageHandler))
	router.HandleFunc("/notfound", server.SecurityHeadersMiddleware(pages.NotFoundPageHandler))
	router.HandleFunc("/blog", server.SecurityHeadersMiddleware(pages.BlogPageHandler))
	router.HandleFunc("/blog/{slug}", server.SecurityHeadersMiddleware(pages.ArticlePageHandler))
	router.HandleFunc("/", server.SecurityHeadersMiddleware(pages.LandingPageHandler))
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
	legal.LoadConfig()
	loadEnvJs()
	loadBuildJs()
	loadIndexHtml()
	loadRobotsTxt()
	loadSitemapXml()
	loadMicrosoftApp()
	pages.InitTemplates()
	pages.InitBlogClient()
	router := setupRouter()
	cors := server.ConfigureCors(router)
	server.Start(cors, nil)
}
