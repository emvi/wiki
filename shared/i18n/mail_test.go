package i18n

import "testing"

func TestGetMailTitle(t *testing.T) {
	if GetMailTitle("jp")["password_mail"] != "Your password at Emvi has been reset" {
		t.Fatal("Unexpected title")
	}

	if GetMailTitle("en")["password_mail"] != "Your password at Emvi has been reset" {
		t.Fatal("Unexpected title")
	}

	if GetMailTitle("de")["password_mail"] != "Dein Passwort bei Emvi wurde zurückgesetzt" {
		t.Fatal("Unexpected title")
	}
}

func TestGetMailEndI18n(t *testing.T) {
	if GetMailEndI18n("jp")["end_terms"] != "Terms" {
		t.Fatal("Unexpected vars")
	}

	if GetMailEndI18n("en")["end_terms"] != "Terms" {
		t.Fatal("Unexpected vars")
	}

	if GetMailEndI18n("de")["end_terms"] != "Nutzungsbedingungen" {
		t.Fatal("Unexpected vars")
	}
}
