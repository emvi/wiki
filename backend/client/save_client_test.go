package client

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSaveClientDataValidateName(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateClient(t, orga, "name", "id", "secret")
	data := &SaveClientData{Name: ""}

	if err := data.validate(orga.ID); len(err) != 1 && err[0] != errs.NameEmpty {
		t.Fatalf("Expected client name to be invalid, but was: %v", err)
	}

	data.Name = "01234567890123456789012345678901234567891"

	if err := data.validate(orga.ID); len(err) != 1 && err[0] != errs.UsernameInvalid {
		t.Fatalf("Expected client name to be invalid, but was: %v", err)
	}

	data.Name = "name"

	if err := data.validate(orga.ID); len(err) != 1 && err[0] != errs.NameInUse {
		t.Fatalf("Expected client name to be in use, but was: %v", err)
	}

	data.Name = "newclient"

	if err := data.validate(orga.ID); err != nil {
		t.Fatalf("Expected client name to be valid, but was: %v", err)
	}
}

func TestSaveClientDataValidateScopes(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateClient(t, orga, "name", "id", "secret")
	data := &SaveClientData{Name: "valid", Scopes: []Scope{
		{"unknown", true, false},
	}}

	if err := data.validate(orga.ID); len(err) != 1 && err[0] != errs.ScopeInvalid {
		t.Fatalf("Expected scope to be invalid, but was: %v", err)
	}

	data.Scopes = []Scope{
		{"articles", true, true},
	}

	if err := data.validate(orga.ID); len(err) != 1 && err[0] != errs.ScopeInvalid {
		t.Fatalf("Expected scope to be invalid, but was: %v", err)
	}

	data.Scopes = []Scope{
		{"articles", false, false},
	}

	if err := data.validate(orga.ID); err != nil {
		t.Fatalf("Expected scope to be valid, but was: %v", err)
	}

	if len(data.Scopes) != 0 {
		t.Fatalf("No scopes should be remaining, but was: %v", len(data.Scopes))
	}

	data.Scopes = []Scope{
		{"articles", true, false},
		{"articles", true, false},
	}

	if err := data.validate(orga.ID); err != nil {
		t.Fatalf("Expected scope to be valid, but was: %v", err)
	}

	if len(data.Scopes) != 1 {
		t.Fatalf("One scopes should be remaining, but was: %v", len(data.Scopes))
	}
}

func TestSaveClientNotAdmin(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if err := SaveClient(orga, user.ID+1, nil, nil); len(err) != 1 && err[0] != errs.PermissionDenied {
		t.Fatalf("Expected access to be denied, but was: %v", err)
	}
}

func TestSaveClientUpdateClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	authProvider := auth.NewMockAuthClient()
	orga, user := testutil.CreateOrgaAndUser(t)
	client := testutil.CreateClient(t, orga, "name", "id", "secret")
	data := &SaveClientData{Id: client.ID + 1, Name: "newname"}

	if err := SaveClient(orga, user.ID, data, authProvider); len(err) != 1 || err[0] != errs.ClientNotFound {
		t.Fatalf("Expected client not to be found, but was: %v", err)
	}

	data.Id = client.ID

	if err := SaveClient(orga, user.ID, data, authProvider); err != nil {
		t.Fatalf("Expected client name to be updated, but was: %v", err)
	}

	client = model.GetClientByOrganizationIdAndId(orga.ID, client.ID)

	if client.Name != "newname" {
		t.Fatalf("Client name must have been changed, but was: %v", client.Name)
	}
}

func TestSaveClientNewClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	authProvider := auth.NewMockAuthClient()
	authProvider.NewClientMockResponse = &auth.NewClientResponse{ClientId: "01234567890123456789",
		ClientSecret: "0123456789012345678901234567890123456789012345678901234567891234"}
	orga, user := testutil.CreateOrgaAndUser(t)
	scopes := []Scope{
		{"articles", true, false},
		{"articles", true, false}, // this is on purpose
	}
	data := &SaveClientData{Name: "name", Scopes: scopes}

	if err := SaveClient(orga, user.ID, data, authProvider); err != nil {
		t.Fatalf("New client must have been created, but was: %v", err)
	}

	client := model.GetClientByOrganizationIdAndName(orga.ID, "name")

	if client == nil {
		t.Fatal("Client must have been created")
	}

	if client.Name != "name" ||
		len(client.ClientId) != 20 ||
		len(client.ClientSecret) != 64 {
		t.Fatalf("Client not as expected: %v", client)
	}

	clientScopes := model.FindClientScopeByClientId(client.ID)

	if len(clientScopes) != 1 ||
		clientScopes[0].Name != "articles" ||
		!clientScopes[0].Read ||
		clientScopes[0].Write {
		t.Fatalf("Scope not as expected: %v", clientScopes)
	}

	if authProvider.NewClientCalls != 1 {
		t.Fatalf("New client must have been called once, but was: %v", authProvider.NewClientCalls)
	}
}

func TestSaveClientNameExistsDifferentOrga(t *testing.T) {
	testutil.CleanBackendDb(t)
	authProvider := auth.NewMockAuthClient()
	authProvider.NewClientMockResponse = &auth.NewClientResponse{ClientId: "01234567890123456789",
		ClientSecret: "0123456789012345678901234567890123456789012345678901234567891234"}
	orga, user := testutil.CreateOrgaAndUser(t)
	orga2, _ := testutil.CreateOrga(t, user, "orga2")
	testutil.CreateClient(t, orga2, "name", "id", "secret")
	scopes := []Scope{{"articles", true, false}}
	data := &SaveClientData{Name: "name", Scopes: scopes}

	if err := SaveClient(orga, user.ID, data, authProvider); err != nil {
		t.Fatalf("New client must have been created, but was: %v", err)
	}
}
