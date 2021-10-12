package auth

import (
	"bytes"
	"emviwiki/shared/constants"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	clientIdParam        = "client_id"
	clientSecretParam    = "client_secret"
	tokenEndpoint        = "/api/v1/auth/token"
	meEndpoint           = "/api/v1/auth/user"
	newClientEndpoint    = "/api/v1/auth/client"
	deleteClientEndpoint = "/api/v1/auth/client"
)

var (
	TokenExpiredErr = errors.New("token expired")
)

type EmviAuthClient struct {
	ClientId     string
	ClientSecret string
	TokenType    string
	Token        string
	m            sync.RWMutex
}

func NewEmviAuthClient(clientId, clientSecret string) *EmviAuthClient {
	return &EmviAuthClient{ClientId: clientId, ClientSecret: clientSecret}
}

func (auth *EmviAuthClient) ValidateToken(r *http.Request) (*TokenResponse, error) {
	req, err := http.NewRequest(http.MethodGet, Host+tokenEndpoint, nil)

	if err != nil {
		return nil, err
	}

	authHeader := r.Header.Get(constants.AuthHeader)
	clientHeader := r.Header.Get(constants.AuthClientHeader)
	req.Header.Add(constants.AuthHeader, authHeader)
	req.Header.Add(constants.AuthClientHeader, clientHeader)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for token validation request", logbuch.Fields{"err": err})
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logbuch.Debug("Could not obtain user token", logbuch.Fields{
			"status":                   resp.StatusCode,
			constants.AuthHeader:       authHeader,
			constants.AuthClientHeader: clientHeader,
		})
		return nil, errors.New("could not obtain user token")
	}

	tokenResp := new(TokenResponse)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, tokenResp); err != nil {
		return nil, err
	}

	return tokenResp, nil
}

func (auth *EmviAuthClient) GetActiveUser(w http.ResponseWriter, r *http.Request) (*UserResponse, error) {
	client := http.Client{}
	req, _ := http.NewRequest(http.MethodGet, Host+meEndpoint, nil)
	req.Header.Add(constants.AuthHeader, r.Header.Get(constants.AuthHeader))
	resp, err := client.Do(req)

	if err != nil {
		logbuch.Error("Error requesting user data from authentication server", logbuch.Fields{"err": err})
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for active user request", logbuch.Fields{"err": err})
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logbuch.Error("Error reading user data from authentication server", logbuch.Fields{"err": err})
		return nil, err
	}

	errResp := &struct {
		Errors []string `json:"error"`
	}{}

	if err := json.Unmarshal(body, errResp); err != nil {
		logbuch.Error("Error decoding user data error response from authentication server")
		return nil, err
	}

	if len(errResp.Errors) != 0 {
		logbuch.Error("Unexpected errors from authentication server", logbuch.Fields{"errs": errResp.Errors})
		return nil, err
	}

	userData := &struct {
		User UserResponse `json:"user"`
	}{}

	if err := json.Unmarshal(body, userData); err != nil {
		logbuch.Error("Error decoding user data from authentication server", logbuch.Fields{"err": err})
		return nil, err
	}

	return &userData.User, nil
}

func (auth *EmviAuthClient) NewClient(name string, scopes map[string]string) (*NewClientResponse, error) {
	resp, err := auth.performNewClientRequest(name, scopes)

	// refresh token and try again
	if err == TokenExpiredErr {
		if e := auth.obtainToken(); e != nil {
			return nil, err
		}

		resp, err = auth.performNewClientRequest(name, scopes)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (auth *EmviAuthClient) DeleteClient(clientId, clientSecret string) error {
	err := auth.performDeleteClientRequest(clientId, clientSecret)

	// refresh token and try again
	if err == TokenExpiredErr {
		if e := auth.obtainToken(); e != nil {
			return err
		}

		err = auth.performDeleteClientRequest(clientId, clientSecret)
	}

	if err != nil {
		return err
	}

	return nil
}

func (auth *EmviAuthClient) performNewClientRequest(name string, scopes map[string]string) (*NewClientResponse, error) {
	data := struct {
		Name   string            `json:"name"`
		Scopes map[string]string `json:"scopes"`
	}{name, scopes}
	jsonData, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, Host+newClientEndpoint, bytes.NewReader(jsonData))

	if err != nil {
		return nil, err
	}

	auth.m.RLock()
	req.Header.Add(constants.AuthHeader, fmt.Sprintf("%s %s", constants.AuthTokenType, auth.Token))
	auth.m.RUnlock()
	req.Header.Add(constants.AuthClientHeader, auth.ClientId)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for new client request", logbuch.Fields{"err": err})
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, TokenExpiredErr
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error creating new client")
	}

	tokenResp := new(NewClientResponse)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, tokenResp); err != nil {
		return nil, err
	}

	return tokenResp, nil
}

func (auth *EmviAuthClient) performDeleteClientRequest(clientId, clientSecret string) error {
	logbuch.Debug("Performing delete client request", logbuch.Fields{"client_id": clientId, "client_secret": clientSecret})
	req, err := http.NewRequest(http.MethodDelete, Host+deleteClientEndpoint, nil)

	if err != nil {
		return err
	}

	auth.m.RLock()
	req.Header.Add(constants.AuthHeader, fmt.Sprintf("%s %s", constants.AuthTokenType, auth.Token))
	auth.m.RUnlock()
	req.Header.Add(constants.AuthClientHeader, auth.ClientId)
	query := req.URL.Query()
	query.Add(clientIdParam, clientId)
	query.Add(clientSecretParam, clientSecret)
	req.URL.RawQuery = query.Encode()
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for delete client request", logbuch.Fields{"err": err})
		}
	}()

	if resp.StatusCode == http.StatusUnauthorized {
		return TokenExpiredErr
	} else if resp.StatusCode != http.StatusOK {
		return errors.New("error deleting client")
	}

	return nil
}

func (auth *EmviAuthClient) obtainToken() error {
	auth.m.Lock()
	defer auth.m.Unlock()

	data := struct {
		GrantType    string `json:"grant_type"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{constants.AuthGrantType, auth.ClientId, auth.ClientSecret}
	jsonData, err := json.Marshal(data)

	if err != nil {
		logbuch.Error("Error marshalling token request", logbuch.Fields{"err": err})
		return err
	}

	req, err := http.NewRequest(http.MethodPost, Host+tokenEndpoint, bytes.NewReader(jsonData))

	if err != nil {
		logbuch.Error("Error building token request", logbuch.Fields{"err": err})
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logbuch.Error("Error while executing token request", logbuch.Fields{"err": err})
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing response body for token request", logbuch.Fields{"err": err})
		}
	}()

	if resp.StatusCode != http.StatusOK {
		logbuch.Debug("Could not obtain token", logbuch.Fields{"status": resp.StatusCode})
		return errors.New("could not obtain token")
	}

	tokenResp := &struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
		// we don't care about TTL here
	}{}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logbuch.Debug("Error reading token response body", logbuch.Fields{"status": resp.StatusCode})
		return err
	}

	if err := json.Unmarshal(body, tokenResp); err != nil {
		logbuch.Debug("Error unmarshalling token response body", logbuch.Fields{"status": resp.StatusCode})
		return err
	}

	auth.Token = tokenResp.AccessToken
	auth.TokenType = tokenResp.TokenType
	logbuch.Debug("Successfully obtained client token", logbuch.Fields{"token": auth.Token, "token_type": auth.TokenType})
	return nil
}
