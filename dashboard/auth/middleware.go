package auth

import (
	"emviwiki/shared/constants"
	"errors"
	"net/http"
	"strings"
)

type AuthHandler func(claims *UserTokenClaims, w http.ResponseWriter, r *http.Request)

func Middleware(next AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)

		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		claims := GetUserTokenClaims(token)

		if claims == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		next(claims, w, r)
	}
}

func getTokenFromRequest(r *http.Request) (string, error) {
	bearer, err := r.Cookie(constants.AuthHeader)

	if err != nil {
		return "", err
	}

	parts := strings.Split(bearer.Value, " ")

	if len(parts) != 2 || parts[0] != constants.AuthTokenType {
		return "", errors.New("bearer content not as expected")
	}

	return parts[1], nil
}
