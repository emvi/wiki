package main

import (
	"emviwiki/auth/api"
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	"emviwiki/auth/pages"
	"emviwiki/auth/sso"
	"emviwiki/auth/user"
	"emviwiki/shared/config"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/recaptcha"
	"emviwiki/shared/rest"
	"emviwiki/shared/server"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	staticDir       = "static"
	staticDirPrefix = "/auth/static/"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// pages
	router.NotFoundHandler = http.HandlerFunc(pages.NotFoundPageHandler)
	router.HandleFunc("/auth/login", pages.LoginPageHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/auth/logout", pages.LogoutPageHandler).Methods(http.MethodGet)
	router.HandleFunc("/auth/authorize", pages.AuthorizationPageHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/auth/password", pages.PasswordPageHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/auth/passwordreset", pages.PasswordResetPageHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/auth/email", pages.UpdateUserEmailPageHandler).Methods(http.MethodGet)
	router.HandleFunc("/auth/sso/{provider}", pages.SSOPageHandler).Methods(http.MethodGet)

	// REST endpoints
	router.Handle("/api/v1/auth/token", api.AuthMiddleware(api.ValidateTokenHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/auth/token", rest.ErrorMiddleware(api.ClientCredentialsHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/registration", rest.ErrorMiddleware(api.ConfirmRegistrationHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/auth/registration", rest.ErrorMiddleware(api.RegistrationHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/registration", rest.ErrorMiddleware(api.CancelRegistrationHandler)).Methods(http.MethodDelete)
	router.Handle("/api/v1/auth/registration/password", rest.ErrorMiddleware(api.RegistrationPasswordHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/registration/personal", rest.ErrorMiddleware(api.RegistrationPersonalHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/registration/completion", rest.ErrorMiddleware(api.RegistrationCompletionHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/user", api.AuthMiddleware(api.GetUserHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/auth/user/email", api.AuthMiddleware(api.UpdateUserEmailHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/user/email", rest.ErrorMiddleware(api.UpdateUserEmailConfirmationHandler)).Methods(http.MethodGet)
	router.Handle("/api/v1/auth/user/data", api.AuthMiddleware(api.UpdateUserDataHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/user/password", api.AuthMiddleware(api.UpdatePasswordHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/client", api.AuthMiddleware(api.TrustedClientMiddleware(api.RegisterClientHandler))).Methods(http.MethodPost)
	router.Handle("/api/v1/auth/client", api.AuthMiddleware(api.TrustedClientMiddleware(api.DeleteClientHandler))).Methods(http.MethodDelete)

	// static content
	server.ServeStaticFiles(router, staticDirPrefix, staticDir, config.Get().Version)

	return router
}

func connectDB() *db.Connection {
	auth := config.Get().AuthDB
	return db.NewConnection(db.ConnectionData{
		Host:               auth.Host,
		Port:               auth.Port,
		User:               auth.User,
		Password:           auth.Password,
		Schema:             auth.Schema,
		SSLMode:            auth.SSLMode,
		SSLCert:            auth.SSLCert,
		SSLKey:             auth.SSLKey,
		SSLRootCert:        auth.SSLRootCert,
		MaxOpenConnections: auth.MaxOpenConnections,
	})
}

func main() {
	config.Load()
	stdout, stderr := server.ConfigureLogging()
	defer server.CloseLogger(stdout, stderr)
	db.Migrate()
	mail.LoadConfig()
	i18n.LoadConfig()
	user.LoadConfig()
	sso.LoadConfig()
	jwt.LoadRSAKeys()
	api.LoadConfig()
	pages.LoadConfig()
	pages.InitTemplates()
	user.InitTemplates()
	recaptcha.LoadConfig()
	connection := connectDB()
	model.SetConnection(connection)
	defer connection.Disconnect()
	router := setupRouter()
	cors := server.ConfigureCors(router)
	server.Start(cors, nil)
}
