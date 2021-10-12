package user

import (
	"emviwiki/shared/config"
	"emviwiki/shared/constants"
	"net/http"
	"time"
)

// SetSessionCookie sets the access token cookie.
// The cookie is accessible from JavaScript and valid on subdomains.
func SetSessionCookie(w http.ResponseWriter, value string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.AuthCookieName,
		Value:    value,
		Expires:  expires,
		Secure:   config.Get().Server.HTTP.SecureCookies,
		HttpOnly: false,
		Domain:   config.Get().JWT.CookieDomainName,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}
