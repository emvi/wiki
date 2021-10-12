package sso

import (
	"bytes"
	"emviwiki/shared/constants"
	"fmt"
	"github.com/emvi/logbuch"
	"net/http"
	"strings"
)

const (
	microsoftOAuthLogin     = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	microsoftOAuthTokenBody = "client_id=%s&client_secret=%s&scope=user.read&code=%s&redirect_uri=%s&grant_type=authorization_code"
	microsoftUserEndpoint   = "https://graph.microsoft.com/v1.0/me"
	microsoftDefaultName    = "Name Surname"
)

type MicrosoftSSOProvider struct {
	clientId     string
	clientSecret string
}

func NewMicrosoftSSOProvider(clientId, clientSecret string) *MicrosoftSSOProvider {
	return &MicrosoftSSOProvider{clientId, clientSecret}
}

func (provider *MicrosoftSSOProvider) GetToken(code string) (SSOTokenResponse, error) {
	if code == "" {
		return SSOTokenResponse{}, accessDeniedError
	}

	redirectUri := fmt.Sprintf("%s/auth/sso/microsoft", authHost)
	body := fmt.Sprintf(microsoftOAuthTokenBody, provider.clientId, provider.clientSecret, code, redirectUri)
	req, err := http.NewRequest(http.MethodPost, microsoftOAuthLogin, bytes.NewReader([]byte(body)))

	if err != nil {
		return SSOTokenResponse{}, err
	}

	req.Header.Add("Accept", "application/x-www-form-urlencoded")
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

func (provider *MicrosoftSSOProvider) GetUser(token string) (SSOUser, error) {
	req, err := http.NewRequest(http.MethodGet, microsoftUserEndpoint, nil)

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

	userResp := struct {
		Id                string `json:"id"`
		GivenName         string `json:"givenName"`
		Surname           string `json:"surname"`
		DisplayName       string `json:"displayName"` // in case GivenName and Surname are not set
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"` // in case the mail is not set
	}{}

	if err := decodeResponseBody(resp, &userResp); err != nil {
		return SSOUser{}, err
	}

	if userResp.Id == "" || (userResp.Mail == "" && userResp.UserPrincipalName == "") {
		logbuch.Warn("User data incomplete", logbuch.Fields{"user": userResp})
		return SSOUser{}, userDataError
	}

	if userResp.Mail == "" {
		userResp.Mail = userResp.UserPrincipalName
	}

	userResp.GivenName = strings.TrimSpace(userResp.GivenName)
	userResp.Surname = strings.TrimSpace(userResp.Surname)
	userResp.DisplayName = strings.TrimSpace(userResp.DisplayName)
	name := ""

	if userResp.GivenName != "" && userResp.Surname != "" {
		name = fmt.Sprintf(`%s %s`, userResp.GivenName, userResp.Surname)
	} else if userResp.DisplayName != "" {
		name = userResp.DisplayName
	} else {
		name = microsoftDefaultName
	}

	return SSOUser{Id: userResp.Id,
		Email: userResp.Mail,
		Name:  name}, nil
}
