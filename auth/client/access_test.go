package client

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/constants"
	"github.com/emvi/null"
	"testing"
)

func TestValidateClientCredentialsGrantType(t *testing.T) {
	if _, err := ValidateClientCredentials("invalid", "client", "secret"); err != errs.GrantTypeInvalid {
		t.Fatalf("Grant type must be invalid, but was: %v", err)
	}
}

func TestValidateClientCredentialsClientNotFound(t *testing.T) {
	if _, err := ValidateClientCredentials(constants.AuthGrantType, "invalid", "invalid"); err != errs.ClientCredentialsInvalid {
		t.Fatalf("Client credentials must be invalid, but was: %v", err)
	}
}

func TestValidateClientCredentialsSuccess(t *testing.T) {
	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "client"`); err != nil {
		t.Fatal(err)
	}

	client := &model.Client{Name: "test",
		ClientId:     "client",
		ClientSecret: "secret",
		RedirectURI:  null.NewString("https://test.com/", true),
		Trusted:      false}

	if err := model.SaveClient(nil, client); err != nil {
		t.Fatal(err)
	}

	scope1 := &model.Scope{ClientId: client.ID, Key: "key1", Value: "value1"}
	scope2 := &model.Scope{ClientId: client.ID, Key: "key2", Value: "value2"}

	if err := model.SaveScope(nil, scope1); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveScope(nil, scope2); err != nil {
		t.Fatal(err)
	}

	resp, err := ValidateClientCredentials(constants.AuthGrantType, "client", "secret")

	if err != nil {
		t.Fatalf("Client credentials must be valid, but was: %v", err)
	}

	if resp.TokenType != constants.AuthTokenType ||
		resp.ExpiresIn <= 0 ||
		len(resp.AccessToken) == 0 {
		t.Fatalf("Response object invalid, was: %v", resp)
	}
}
