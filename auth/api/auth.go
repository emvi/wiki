package api

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	"emviwiki/shared/auth"
	"emviwiki/shared/constants"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
	"time"
)

type TokenResponse struct {
	ExpiresIn int64    `json:"expires_in"`
	Scopes    []string `json:"scopes"`
	UserId    hide.ID  `json:"user_id"`
	Language  *string  `json:"language"`
	ClientId  string   `json:"client_id"`
	Trusted   bool     `json:"trusted"`
	IsSSOUser bool     `json:"is_sso_user"`
}

func (resp *TokenResponse) IsClient() bool {
	return resp.UserId == 0
}

type AuthContext struct {
	Token string
	TokenResponse
}

type AuthHandler func(*AuthContext, http.ResponseWriter, *http.Request) []error

func AuthMiddleware(next AuthHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := getTokenResponse(w, r)

		if ctx == nil {
			return
		}

		err := next(&AuthContext{r.Header.Get(constants.AuthHeader), *ctx}, w, r)
		rest.HandleErrors(w, r, err...)
	})
}

func TrustedClientMiddleware(next AuthHandler) AuthHandler {
	return func(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
		if !ctx.IsClient() || model.GetClientByClientIdAndTrusted(ctx.ClientId) == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return nil
		}

		return next(ctx, w, r)
	}
}

func ValidateTokenHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	if ctx.IsClient() {
		rest.WriteResponse(w, struct {
			ExpiresIn int64    `json:"expires_in"`
			Scopes    []string `json:"scopes"`
			ClientId  string   `json:"client_id"`
			Trusted   bool     `json:"trusted"`
		}{ctx.TokenResponse.ExpiresIn, ctx.Scopes, ctx.TokenResponse.ClientId, ctx.TokenResponse.Trusted})
	} else {
		rest.WriteResponse(w, struct {
			TokenResponse
			ClientId string `json:"-"`
		}{TokenResponse: ctx.TokenResponse})
	}

	return nil
}

func getTokenResponse(w http.ResponseWriter, r *http.Request) *TokenResponse {
	token := auth.GetAccessToken(r)

	if token == "" {
		rest.WriteErrorResponse(w, http.StatusUnauthorized, errs.TokenInvalid)
		return nil
	}

	client := r.Header.Get(constants.AuthClientHeader)
	response := new(TokenResponse)

	if client == "" {
		claims := jwt.GetUserTokenClaims(token)

		// check for user id to prevent user token to be confused with client token
		if claims == nil || claims.UserId == 0 {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, errs.TokenInvalid)
			return nil
		}

		response.ExpiresIn = claims.ExpiresAt - time.Now().Unix()
		response.Scopes = claims.Scopes
		response.UserId = claims.UserId
		response.Language = &claims.Language
		response.IsSSOUser = claims.IsSSOUser

		if claims.Language == "" {
			response.Language = nil
		}
	} else {
		claims := jwt.GetClientTokenClaims(token)

		if claims == nil || claims.ClientId != client {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, errs.TokenInvalid)
			return nil
		}

		response.ExpiresIn = claims.ExpiresAt - time.Now().Unix()
		response.Scopes = claims.Scopes
		response.ClientId = claims.ClientId
		response.Trusted = claims.Trusted
	}

	return response
}
