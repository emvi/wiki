package rest

import (
	"net/http"
	"testing"
)

func TestGetLangCode(t *testing.T) {
	input := []string{"", "en", "e", "0000", "de", "de-DE,de;q=0.8,ro;q=0.6"}
	expected := []string{"en", "en", "en", "en", "de", "de"}

	for i, in := range input {
		r, _ := http.NewRequest(http.MethodGet, "", nil)
		r.Header.Add("accept-language", in)

		if GetLangCode(r) != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], in)
		}
	}
}

func TestGetSupportedLangCode(t *testing.T) {
	input := []string{"en", "de", "jp"}
	expected := []string{"en", "de", "en"}

	for i, in := range input {
		r, _ := http.NewRequest(http.MethodGet, "", nil)
		r.Header.Add("accept-language", in)

		if GetSupportedLangCode(r) != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], in)
		}
	}
}
