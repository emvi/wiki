package auth

import (
	"github.com/emvi/hide"
	"net/http"
	"time"
)

// AuthClient is an interface to an authentication service.
type AuthClient interface {
	// ValidateToken validates the token within a request and returns the token data or an error.
	ValidateToken(*http.Request) (*TokenResponse, error)

	// GetActiveUser returns the user data for a request or an error.
	GetActiveUser(http.ResponseWriter, *http.Request) (*UserResponse, error)

	// NewClient creates a new client if this client is trusted or else an error is returned.
	// The name must be unique. The scopes are passed as a map of key and values (e.g. object:rw).
	NewClient(string, map[string]string) (*NewClientResponse, error)

	// DeleteClient deletes a client if this client is trusted
	// for given client_id and client_secret or else an error is returned.
	DeleteClient(string, string) error
}

// TokenResponse wraps the response from a token validation request.
type TokenResponse struct {
	ExpiresIn int64    `json:"expires_in"`
	Scopes    []string `json:"scopes"`
	UserId    hide.ID  `json:"user_id"`
	Language  *string  `json:"language"`
	ClientId  string   `json:"client_id"`
	Trusted   bool     `json:"trusted"`
	IsSSOUser bool     `json:"is_sso_user"`
}

// IsClient returns true if the response contains a valid user ID,
// which means it was used for an implicit grant or false if it was used for a client credential grant type.
func (resp *TokenResponse) IsClient() bool {
	return resp.UserId == 0
}

// UserResponse wraps the response for retrieving a user.
type UserResponse struct {
	Id              hide.ID   `json:"id"`
	Email           string    `json:"email"`
	Firstname       string    `json:"firstname"`
	Lastname        string    `json:"lastname"`
	Language        string    `json:"language"`
	Picture         string    `json:"picture"`
	Info            string    `json:"info"`
	AcceptMarketing bool      `json:"accept_marketing"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
	IsSSOUser       bool      `json:"is_sso_user"`
}

// NewClientResponse wraps the response for creating a new client.
type NewClientResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
