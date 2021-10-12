package client

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestNewClientFailure(t *testing.T) {
	testutil.CleanAuthDb(t)
	existingClient := &model.Client{Name: "existing"}

	if err := model.SaveClient(nil, existingClient); err != nil {
		t.Fatal(err)
	}

	input := []struct {
		name   string
		scopes map[string]string
	}{
		{"", map[string]string{}},
		{"01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891", map[string]string{}},
		{"name", map[string]string{
			"key": "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891",
		}},
		{"name", map[string]string{
			"01234567890123456789012345678901234567891": "value",
		}},
		{"existing", map[string]string{
			"key": "value",
		}},
	}
	expected := []error{
		errs.ClientNameInvalid,
		errs.ClientNameInvalid,
		errs.ScopeInvalid,
		errs.ScopeInvalid,
		errs.ClientNameInUse,
	}

	for i, in := range input {
		if _, err := NewClient(in.name, in.scopes); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}

func TestNewClientSuccess(t *testing.T) {
	testutil.CleanAuthDb(t)
	scopes := map[string]string{
		"scope1": "rw",
		"scope2": "r",
	}
	client, err := NewClient("name", scopes)

	if err != nil {
		t.Fatalf("Expected new client to be created, but was: %v", err)
	}

	if client.Name != "name" ||
		len(client.ClientId) != 20 ||
		len(client.ClientSecret) != 64 ||
		client.Trusted {
		t.Fatalf("New client not as expected: %v", client)
	}
}
