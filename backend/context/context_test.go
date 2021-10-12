package context

import (
	"emviwiki/backend/client"
	"emviwiki/shared/model"
	"testing"
)

func TestNewEmviContext(t *testing.T) {
	scopes := []string{"unknown:rw", "articles", "articles:r", "lists:r", "nope:r", "tags:r"}
	ctx := NewEmviContext(nil, 0, scopes, false)

	if len(ctx.Scopes) != 3 {
		t.Fatalf("Context must have three valid scopes")
	}
}

func TestEmviContext_HasScopes(t *testing.T) {
	needle := client.Scope{"name", true, false}
	scopes := []string{"articles:r", "lists:rw", "tags:r"}
	orga := &model.Organization{Expert: true}
	ctx := NewEmviContext(orga, 0, scopes, false)

	if ctx.HasScopes(needle) {
		t.Fatalf("Haystack must not contain needle")
	}

	needle.Name = "tags"

	if !ctx.HasScopes(needle) {
		t.Fatalf("Haystack must contain needle")
	}

	needle2 := client.Scope{"lists", true, true}

	if !ctx.HasScopes(needle, needle2) {
		t.Fatalf("Haystack must contain needles")
	}

	needle.Name = "name"
	ctx.TrustedClient = true

	if !ctx.HasScopes(needle, needle2) {
		t.Fatalf("Client must be trusted")
	}
}
