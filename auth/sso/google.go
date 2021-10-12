package sso

import (
	"emviwiki/shared/constants"
	"fmt"
	"github.com/emvi/logbuch"
	"net/http"
)

const (
	googleOAuthLogin   = "https://oauth2.googleapis.com/token?client_id=%s&client_secret=%s&redirect_uri=%s&code=%s&grant_type=authorization_code"
	googleUserEndpoint = "https://www.googleapis.com/oauth2/v1/userinfo?alt=json"
)

type GoogleSSOProvider struct {
	clientId     string
	clientSecret string
}

func NewGoogleSSOProvider(clientId, clientSecret string) *GoogleSSOProvider {
	return &GoogleSSOProvider{clientId, clientSecret}
}

func (provider *GoogleSSOProvider) GetToken(code string) (SSOTokenResponse, error) {
	if code == "" {
		return SSOTokenResponse{}, accessDeniedError
	}

	redirectUri := fmt.Sprintf("%s/auth/sso/google", authHost)
	url := fmt.Sprintf(googleOAuthLogin, provider.clientId, provider.clientSecret, redirectUri, code)
	req, err := http.NewRequest(http.MethodPost, url, nil)

	if err != nil {
		return SSOTokenResponse{}, err
	}

	req.Header.Add("Accept", "application/json")
	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return SSOTokenResponse{}, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for SSO token request", logbuch.Fields{"err": err})
		}
	}()

	var tokenResp SSOTokenResponse

	if err := decodeResponseBody(resp, &tokenResp); err != nil {
		return SSOTokenResponse{}, err
	}

	if tokenResp.AccessToken == "" {
		return SSOTokenResponse{}, accessDeniedError
	}

	return tokenResp, nil
}

func (provider *GoogleSSOProvider) GetUser(token string) (SSOUser, error) {
	req, err := http.NewRequest(http.MethodGet, googleUserEndpoint, nil)

	if err != nil {
		return SSOUser{}, err
	}

	req.Header.Add(constants.AuthHeader, fmt.Sprintf("%s %s", constants.AuthTokenType, token))
	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return SSOUser{}, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for SSO user request", logbuch.Fields{"err": err})
		}
	}()

	var userResp SSOUser

	if err := decodeResponseBody(resp, &userResp); err != nil {
		return SSOUser{}, err
	}

	if userResp.Id == "" || userResp.Name == "" || userResp.Email == "" {
		logbuch.Warn("User data incomplete", logbuch.Fields{"user": userResp})
		return SSOUser{}, userDataError
	}

	return userResp, nil
}
