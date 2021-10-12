package client

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadClients(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	user := testutil.CreateUser(t, orga, 321, "user@test.com")
	client1 := testutil.CreateClient(t, orga, "client1", "id1", "secret1")
	testutil.CreateClient(t, orga, "client2", "id2", "secret2")
	testutil.CreateClientScope(t, client1, "scope1", true, false)
	testutil.CreateClientScope(t, client1, "scope2", false, true)

	if _, err := ReadClients(orga, user.ID); err != errs.PermissionDenied {
		t.Fatalf("User must not be allowed to read clients, but was: %v", err)
	}

	clients, err := ReadClients(orga, admin.ID)

	if err != nil || len(clients) != 2 {
		t.Fatalf("Admin must be allowed to read clients, but was: %v %v", err, len(clients))
	}

	if clients[0].Name != "client1" || len(clients[0].Scopes) != 2 {
		t.Fatalf("First client must have two scopes, but was: %v", len(clients[0].Scopes))
	}
}

func TestReadClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	user := testutil.CreateUser(t, orga, 321, "user@test.com")
	client := testutil.CreateClient(t, orga, "client", "id", "secret")
	testutil.CreateClientScope(t, client, "scope1", false, true)
	testutil.CreateClientScope(t, client, "scope2", true, false)

	if _, err := ReadClient(orga, user.ID, 0); err != errs.PermissionDenied {
		t.Fatalf("User must not be allowed to read client, but was: %v", err)
	}

	if _, err := ReadClient(orga, admin.ID, 0); err != errs.ClientNotFound {
		t.Fatalf("Client must not be found, but was: %v", err)
	}

	result, err := ReadClient(orga, admin.ID, client.ID)

	if err != nil {
		t.Fatalf("Admin must be allowed to read client, but was: %v", err)
	}

	if result.ID != client.ID || len(result.Scopes) != 2 {
		t.Fatalf("Client not as expected: %v", result)
	}
}
