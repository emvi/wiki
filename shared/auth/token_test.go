package auth

import (
	"emviwiki/shared/constants"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccessTokenInvalid(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	if GetAccessToken(r) != "" {
		t.Fatalf("Token must not be set, but was: %v", GetAccessToken(r))
	}
}

func TestGetAccessTokenHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(constants.AuthHeader, constants.AuthTokenType+" valid_token")

	if GetAccessToken(r) != "valid_token" {
		t.Fatalf("Token must have been returned, but was: %v", GetAccessToken(r))
	}
}

func TestGetAccessTokenCookie(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  constants.AuthCookieName,
		Value: "valid_token",
	})

	if GetAccessToken(r) != "valid_token" {
		t.Fatalf("Token must have been returned, but was: %v", GetAccessToken(r))
	}
}
