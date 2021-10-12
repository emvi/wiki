package auth

import "net/http"

type MockAuthClient struct {
	ValidateTokenCalls int
	GetActiveUserCalls int
	NewClientCalls     int
	DeleteClientCalls  int

	NewClientMockResponse *NewClientResponse
}

func NewMockAuthClient() *MockAuthClient {
	return new(MockAuthClient)
}

func (auth *MockAuthClient) ValidateToken(r *http.Request) (*TokenResponse, error) {
	auth.ValidateTokenCalls++
	return nil, nil
}

func (auth *MockAuthClient) GetActiveUser(w http.ResponseWriter, r *http.Request) (*UserResponse, error) {
	auth.GetActiveUserCalls++
	return nil, nil
}

func (auth *MockAuthClient) NewClient(name string, scopes map[string]string) (*NewClientResponse, error) {
	auth.NewClientCalls++
	return auth.NewClientMockResponse, nil
}

func (auth *MockAuthClient) DeleteClient(clientId, clientSecret string) error {
	auth.DeleteClientCalls++
	return nil
}
