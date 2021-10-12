package i18n

import "testing"

func TestGetVars(t *testing.T) {
	if GetVars("jp", mailEndI18n)["end_terms"] != "Terms" {
		t.Fatal("Unexpected vars")
	}

	if GetVars("en", mailEndI18n)["end_terms"] != "Terms" {
		t.Fatal("Unexpected vars")
	}

	if GetVars("de", mailEndI18n)["end_terms"] != "Nutzungsbedingungen" {
		t.Fatal("Unexpected vars")
	}
}
