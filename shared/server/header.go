package server

import (
	"net/http"
)

// AddSecurityHeaders adds security HTTP response headers.
func AddSecurityHeaders(w http.ResponseWriter) {
	w.Header().Add("X-Frame-Options", "deny")
}

// SecurityHeadersMiddleware adds security HTTP response headers to the handler function.
func SecurityHeadersMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AddSecurityHeaders(w)
		next(w, r)
	}
}
