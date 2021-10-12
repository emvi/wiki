package jwt

import (
	"testing"
	"time"
)

func TestGetClientTokenClaims(t *testing.T) {
	now := time.Now()
	token, ttl, err := NewClientToken(&ClientTokenClaims{ClientId: "client", Trusted: true, Scopes: []string{"scope"}})

	if err != nil || token == "" || ttl.Before(now) {
		t.Fatal("New token must be created")
	}

	claims := GetClientTokenClaims(token)

	if claims == nil {
		t.Fatal("Claims must be returned")
	}

	if claims.ClientId != "client" || !claims.Trusted || claims.Scopes[0] != "scope" {
		t.Fatal("Claims attributes must be set")
	}
}
