package jwt

import (
	"testing"
	"time"
)

func TestGetUserTokenClaims(t *testing.T) {
	now := time.Now()
	token, ttl, err := NewUserToken(&UserTokenClaims{UserId: 5, Language: "en", Scopes: []string{"scope"}, IsSSOUser: true})

	if err != nil || token == "" || ttl.Before(now) {
		t.Fatal("New token must be created")
	}

	claims := GetUserTokenClaims(token)

	if claims == nil {
		t.Fatal("Claims must be returned")
	}

	if claims.UserId != 5 || claims.Language != "en" || claims.Scopes[0] != "scope" || !claims.IsSSOUser {
		t.Fatal("Claims attributes must be set")
	}
}
