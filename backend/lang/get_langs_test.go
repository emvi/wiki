package lang

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestGetLangs(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	en := &model.Language{Code: "en", Name: "English", OrganizationId: orga.ID}
	de := &model.Language{Code: "de", Name: "Deutsch", OrganizationId: orga.ID}

	if err := model.SaveLanguage(nil, en); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveLanguage(nil, de); err != nil {
		t.Fatal(err)
	}

	langs := GetLangs(orga)

	if len(langs) != 3 {
		t.Fatalf("Expected 3 languages to be returned, but was: %v", len(langs))
	}

	if langs[0].Code != "de" || langs[0].Name != "Deutsch" {
		t.Fatal("First language returned must be de")
	}

	if langs[1].Code != "en" || langs[1].Name != "English" {
		t.Fatal("Second language returned must be en")
	}
}
