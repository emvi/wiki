package sso

import (
	"bytes"
	"emviwiki/shared/constants"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	githubOAuthLogin         = "https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s"
	githubUserEndpoint       = "https://api.github.com/user"
	githubUserEmailsEndpoint = "https://api.github.com/user/emails"
)

type GitHubSSOProvider struct {
	clientId     string
	clientSecret string
}

func NewGitHubSSOProvider(clientId, clientSecret string) *GitHubSSOProvider {
	return &GitHubSSOProvider{clientId, clientSecret}
}

func (provider *GitHubSSOProvider) GetToken(code string) (SSOTokenResponse, error) {
	if code == "" {
		return SSOTokenResponse{}, accessDeniedError
	}

	url := fmt.Sprintf(githubOAuthLogin, provider.clientId, provider.clientSecret, code)
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

func (provider *GitHubSSOProvider) GetUser(token string) (SSOUser, error) {
	req, err := http.NewRequest(http.MethodGet, githubUserEndpoint, nil)

	if err != nil {
		return SSOUser{}, err
	}

	req.Header.Add(constants.AuthHeader, fmt.Sprintf("token %s", token))
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
		Id        int    `json:"id"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}{}

	if err := decodeResponseBody(resp, &userResp); err != nil {
		return SSOUser{}, err
	}

	ssoUser := SSOUser{Id: strconv.Itoa(userResp.Id), Name: userResp.Name, Picture: userResp.AvatarURL}

	if err := getGitHubUserEmail(token, &ssoUser); err != nil {
		return SSOUser{}, err
	}

	if ssoUser.Id == "" || ssoUser.Name == "" || ssoUser.Email == "" {
		logbuch.Warn("User data incomplete", logbuch.Fields{"user": ssoUser})
		return SSOUser{}, userDataError
	}

	return ssoUser, nil
}

func getGitHubUserEmail(token string, user *SSOUser) error {
	req, err := http.NewRequest(http.MethodGet, githubUserEmailsEndpoint, nil)

	if err != nil {
		return err
	}

	req.Header.Add(constants.AuthHeader, fmt.Sprintf("token %s", token))
	var client http.Client
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for SSO user email request", logbuch.Fields{"err": err})
		}
	}()

	emails := []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}{}

	if err := decodeResponseBody(resp, &emails); err != nil {
		return err
	}

	for _, email := range emails {
		if email.Primary {
			user.Email = email.Email
			break
		}
	}

	if user.Email == "" {
		return errors.New("no primary email found")
	}

	return nil
}

func decodeResponseBody(resp *http.Response, obj interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logbuch.Error("Error reading response body", logbuch.Fields{"err": err})
		return err
	}

	logbuch.Debug("SSO response body", logbuch.Fields{"body": string(body)})
	decoder := json.NewDecoder(bytes.NewReader(body))

	if err := decoder.Decode(&obj); err != nil {
		logbuch.Warn("Error decoding response body", logbuch.Fields{"body": string(body)})
		return err
	}

	return nil
}
