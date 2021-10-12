package main

import (
	authmodel "emviwiki/auth/model"
	"emviwiki/dashboard/auth"
	"emviwiki/dashboard/content"
	"emviwiki/dashboard/model"
	"emviwiki/dashboard/pages"
	"emviwiki/shared/config"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	backend "emviwiki/shared/model"
	"emviwiki/shared/server"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	staticDir = "static"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	server.ServeStaticFiles(router, "", staticDir, "")
	router.Handle("/login", http.HandlerFunc(pages.LoginPageHandler)).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/logout", http.HandlerFunc(pages.LogoutPageHandler)).Methods(http.MethodGet)
	router.Handle("/newsletter/edit/{id}", auth.Middleware(pages.NewsletterEditPageHandler)).Methods(http.MethodGet, http.MethodPost, http.MethodPut)
	router.Handle("/newsletter/create", auth.Middleware(pages.NewsletterEditPageHandler)).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/newsletter/{id}", auth.Middleware(pages.NewsletterPageHandler)).Methods(http.MethodDelete)
	router.Handle("/newsletter", auth.Middleware(pages.NewsletterPageHandler)).Methods(http.MethodGet)
	router.Handle("/file/{id}", auth.Middleware(pages.FilesPageHandler)).Methods(http.MethodDelete)
	router.Handle("/files", auth.Middleware(pages.FilesPageHandler)).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/content/{filename}", http.HandlerFunc(content.ReadContentHandler)).Methods(http.MethodGet)
	router.Handle("/", auth.Middleware(pages.StartPageHandler)).Methods(http.MethodGet)
	return router
}

func connectDB() (*db.Connection, *db.Connection, *db.Connection) {
	authdb := config.Get().AuthDB
	backenddb := config.Get().BackendDB
	dashboarddb := config.Get().DashboardDB
	return db.NewConnection(db.ConnectionData{
			Host:               backenddb.Host,
			Port:               backenddb.Port,
			User:               backenddb.User,
			Password:           backenddb.Password,
			Schema:             backenddb.Schema,
			SSLMode:            backenddb.SSLMode,
			SSLCert:            backenddb.SSLCert,
			SSLKey:             backenddb.SSLKey,
			SSLRootCert:        backenddb.SSLRootCert,
			MaxOpenConnections: backenddb.MaxOpenConnections,
		}), db.NewConnection(db.ConnectionData{
			Host:               authdb.Host,
			Port:               authdb.Port,
			User:               authdb.User,
			Password:           authdb.Password,
			Schema:             authdb.Schema,
			SSLMode:            authdb.SSLMode,
			SSLCert:            authdb.SSLCert,
			SSLKey:             authdb.SSLKey,
			SSLRootCert:        authdb.SSLRootCert,
			MaxOpenConnections: authdb.MaxOpenConnections,
		}), db.NewConnection(db.ConnectionData{
			Host:               dashboarddb.Host,
			Port:               dashboarddb.Port,
			User:               dashboarddb.User,
			Password:           dashboarddb.Password,
			Schema:             dashboarddb.Schema,
			SSLMode:            dashboarddb.SSLMode,
			SSLCert:            dashboarddb.SSLCert,
			SSLKey:             dashboarddb.SSLKey,
			SSLRootCert:        dashboarddb.SSLRootCert,
			MaxOpenConnections: dashboarddb.MaxOpenConnections,
		})
}

func main() {
	config.Load()
	stdout, stderr := server.ConfigureLogging()
	defer server.CloseLogger(stdout, stderr)
	db.Migrate()
	auth.LoadRSAKeys()
	mail.LoadConfig()
	i18n.LoadConfig()
	pages.LoadConfig()
	pages.InitTemplates()
	content.LoadConfig()
	backendDB, authDB, dashboardDB := connectDB()
	backend.SetConnection(backendDB)
	authmodel.SetConnection(authDB)
	model.SetConnection(backendDB, authDB, dashboardDB)
	defer backendDB.Disconnect()
	defer authDB.Disconnect()
	defer dashboardDB.Disconnect()
	router := setupRouter()
	cors := server.ConfigureCors(router)
	server.Start(cors, nil)
}
