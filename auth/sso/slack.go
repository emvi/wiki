package sso

import (
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	"net/http"
)

const (
	slackOAuthLogin   = "https://slack.com/api/oauth.access?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s"
	slackUserEndpoint = "https://slack.com/api/users.identity?token=%s"
)

type SlackSSOProvider struct {
	clientId     string
	clientSecret string
}

func NewSlackSSOProvider(clientId, clientSecret string) *SlackSSOProvider {
	return &SlackSSOProvider{clientId, clientSecret}
}

func (provider *SlackSSOProvider) GetToken(code string) (SSOTokenResponse, error) {
	if code == "" {
		return SSOTokenResponse{}, accessDeniedError
	}

	redirectUri := fmt.Sprintf("%s/auth/sso/slack", authHost)
	url := fmt.Sprintf(slackOAuthLogin, provider.clientId, provider.clientSecret, code, redirectUri)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return SSOTokenResponse{}, err
	}

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

	if !tokenResp.OK {
		return SSOTokenResponse{}, errors.New("Slack returned an error")
	}

	if tokenResp.AccessToken == "" {
		return SSOTokenResponse{}, accessDeniedError
	}

	return SSOTokenResponse{AccessToken: tokenResp.AccessToken}, nil
}

func (provider *SlackSSOProvider) GetUser(token string) (SSOUser, error) {
	url := fmt.Sprintf(slackUserEndpoint, token)
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return SSOUser{}, err
	}

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

	userResp := struct {
		OK   bool `json:"ok"`
		User struct {
			Id       string `json:"id"`
			Name     string `json:"name"`
			Email    string `json:"email"`
			Image512 string `json:"image_512"`
		} `json:"user"`
	}{}

	if err := decodeResponseBody(resp, &userResp); err != nil {
		return SSOUser{}, err
	}

	if !userResp.OK {
		logbuch.Warn("Error on user data request", logbuch.Fields{"user": userResp})
		return SSOUser{}, errors.New("Slack returned an error")
	}

	if userResp.User.Id == "" || userResp.User.Name == "" || userResp.User.Email == "" {
		logbuch.Warn("User data incomplete", logbuch.Fields{"user": userResp})
		return SSOUser{}, userDataError
	}

	return SSOUser{Id: userResp.User.Id, Name: userResp.User.Name, Email: userResp.User.Email, Picture: userResp.User.Image512}, nil
}
