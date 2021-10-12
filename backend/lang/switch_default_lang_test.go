package lang

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSwitchDefaultLanguage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)

	if err := SwitchDefaultLanguage(orga, user.ID, langEn.ID+langDe.ID); err != errs.LanguageNotFound {
		t.Fatalf("Language must not have been found, but was: %v", err)
	}

	if err := SwitchDefaultLanguage(orga, user.ID, langEn.ID); err != nil {
		t.Fatalf("Language must have been changed, but was: %v", err)
	}

	defaultLang := model.GetDefaultLanguageByOrganizationId(orga.ID)

	if defaultLang.Code != "en" {
		t.Fatalf("Default language must be en, but was: %v", defaultLang.Code)
	}

	if err := SwitchDefaultLanguage(orga, user.ID, langDe.ID); err != nil {
		t.Fatalf("Language must have been changed, but was: %v", err)
	}

	defaultLang = model.GetDefaultLanguageByOrganizationId(orga.ID)

	if defaultLang.Code != "de" {
		t.Fatalf("Default language must be de, but was: %v", defaultLang.Code)
	}

	langEn = model.GetLanguageByOrganizationIdAndId(orga.ID, langEn.ID)

	if langEn.Default {
		t.Fatal("Old default language must not be default anymore")
	}
}

func TestSwitchDefaultLanguageNonExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)

	if err := SwitchDefaultLanguage(orga, user.ID, langDe.ID); err != errs.RequiresExpertVersion {
		t.Fatalf("Expected requires expert, but was: %v", err)
	}
}
