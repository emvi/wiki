package auth

import (
	"emviwiki/shared/constants"
	"net/http"
	"strings"
)

// GetAccessToken returns the access token from the authorization header or cookie, depending on what is set.
// The token is not validated! This must be done in the afterwards.
func GetAccessToken(r *http.Request) string {
	authHeader := r.Header.Get(constants.AuthHeader)

	if authHeader != "" {
		bearer := strings.Split(r.Header.Get(constants.AuthHeader), " ")

		if len(bearer) != 2 || bearer[0] != constants.AuthTokenType {
			return ""
		}

		return bearer[1]
	} else {
		cookie, err := r.Cookie(constants.AuthCookieName)

		if err != nil {
			return ""
		}

		return cookie.Value
	}
}
