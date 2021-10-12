package client

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDeleteClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	authProvider := auth.NewMockAuthClient()
	orga, admin := testutil.CreateOrgaAndUser(t)
	user := testutil.CreateUser(t, orga, 321, "user@test.com")
	client := testutil.CreateClient(t, orga, "client", "id", "secret")
	testutil.CreateClientScope(t, client, "scope1", false, true)
	testutil.CreateClientScope(t, client, "scope2", true, false)

	if err := DeleteClient(orga, user.ID, 0, authProvider); err != errs.PermissionDenied {
		t.Fatalf("User must not be allowed to delete client, but was: %v", err)
	}

	if err := DeleteClient(orga, admin.ID, 0, authProvider); err != errs.ClientNotFound {
		t.Fatalf("Client must not be found, but was: %v", err)
	}

	if err := DeleteClient(orga, admin.ID, client.ID, authProvider); err != nil {
		t.Fatalf("Client must be deleted, but was: %v", err)
	}

	if model.GetClientByOrganizationIdAndId(orga.ID, client.ID) != nil {
		t.Fatal("Client must not exist anymore")
	}

	if authProvider.DeleteClientCalls != 1 {
		t.Fatalf("Delete client must have been called once")
	}
}
