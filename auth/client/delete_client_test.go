package client

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDeleteClient(t *testing.T) {
	testutil.CleanAuthDb(t)
	scopes := map[string]string{
		"scope1": "rw",
		"scope2": "r",
	}
	client, err := NewClient("name", scopes)

	if err != nil {
		t.Fatal(err)
	}

	if err := DeleteClient("", ""); err != errs.ClientNotFound {
		t.Fatalf("Client must not be found, but was: %v", err)
	}

	if err := DeleteClient(client.ClientId, client.ClientSecret); err != nil {
		t.Fatalf("Client must be deleted, but was: %v", err)
	}

	if model.GetClientByClientId(client.ClientId) != nil {
		t.Fatal("Client must not exist anymore")
	}
}
