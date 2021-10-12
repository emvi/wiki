package server

import (
	"emviwiki/shared/config"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	"strings"
)

// ConfigureCors confgures CORS log level and origins.
// The configuration is non restrictive by default.
func ConfigureCors(router *mux.Router) http.Handler {
	logbuch.Info("Configuring CORS...")

	origins := strings.Split(config.Get().Cors.Origins, ",")
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            strings.ToLower(config.Get().Cors.Loglevel) == "debug",
	})
	return c.Handler(router)
}
