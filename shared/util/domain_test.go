package util

import (
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"testing"
)

func TestInjectSubdomain(t *testing.T) {
	in := "https://emvi.com/some/path"
	out := InjectSubdomain(in, "test")
	expected := "https://test.emvi.com/some/path"

	if out != expected {
		t.Fatalf("Expected subdomain to be injected, but was: %v", out)
	}
}

func TestGetOrganizationURL(t *testing.T) {
	config.Load()
	config.Get().Hosts.Frontend = "https://emvi.com"
	orga := &model.Organization{NameNormalized: "my-orga"}
	out := GetOrganizationURL(orga, "/some/file/path.txt")
	expected := "https://my-orga.emvi.com/some/file/path.txt"

	if out != expected {
		t.Fatalf("Expected propper organization URL, but was: %v", out)
	}
}
