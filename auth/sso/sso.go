package sso

import (
	"errors"
)

var (
	accessDeniedError = errors.New("access denied")
	userDataError     = errors.New("missing user data")
)

// SSOTokenResponse is the response of a token request for a SSO provider.
type SSOTokenResponse struct {
	OK          bool   `json:"ok"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// SSOUser is the information provided by a SSO provider for a given user.
type SSOUser struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// SSOProvider is a SSO provider to authenticate users.
type SSOProvider interface {
	// GetToken returns the access token for given code after the user has been redirected to Emvi.
	GetToken(string) (SSOTokenResponse, error)

	// GetUser returns the user information for the given access token.
	GetUser(string) (SSOUser, error)
}
