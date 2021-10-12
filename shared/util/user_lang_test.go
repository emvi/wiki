package util

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestDetermineLangNil(t *testing.T) {
	testutil.CleanBackendDb(t)

	if lang := DetermineLang(nil, 0, 0, 0); lang == nil || lang.ID != 0 {
		t.Fatal("Passing no parameters to function must return empty language")
	}
}

func TestDetermineLangByParameter(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	createTestLang(t, orga, "name", "na", true)
	de := createTestLang(t, orga, "Deutsch", "de", false)
	createTestLang(t, orga, "English", "en", false)

	if lang := DetermineLang(nil, orga.ID, user.ID, de.ID); lang == nil || lang.ID != de.ID {
		t.Fatal("'de' must be returned")
	}
}

func TestDetermineLangUserDefault(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	createTestLang(t, orga, "name", "na", true)
	createTestLang(t, orga, "Deutsch", "de", false)
	en := createTestLang(t, orga, "English", "en", false)

	if lang := DetermineLang(nil, orga.ID, user.ID, 0); lang == nil || lang.ID != en.ID {
		t.Fatal("'en' must be returned")
	}
}

func TestDetermineLangOrganizationDefault(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	na := createTestLang(t, orga, "name", "na", true)
	createTestLang(t, orga, "Deutsch", "de", false)
	createTestLang(t, orga, "English", "en", false)

	if lang := DetermineLang(nil, orga.ID, 0, 0); lang == nil || lang.ID != na.ID {
		t.Fatal("'na' must be returned")
	}
}

func TestDetermineSystemSupportedLangCode(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	setUserLang(t, user, "jp")

	if code := DetermineSystemSupportedLangCode(orga.ID, user.ID); code != "en" {
		t.Fatalf("Expected default supported lang, but was: %v", code)
	}

	setUserLang(t, user, "de")

	if code := DetermineSystemSupportedLangCode(orga.ID, user.ID); code != "de" {
		t.Fatalf("Expected supported lang 'de', but was: %v", code)
	}

	setUserLang(t, user, "en")

	if code := DetermineSystemSupportedLangCode(orga.ID, user.ID); code != "en" {
		t.Fatalf("Expected supported lang 'en', but was: %v", code)
	}
}

func TestGetSystemSupportedLangCode(t *testing.T) {
	if GetSystemSupportedLangCode("jp") != "en" {
		t.Fatalf("Must have returned default, but was: %v", GetSystemSupportedLangCode("jp"))
	}

	if GetSystemSupportedLangCode("de") != "de" {
		t.Fatalf("Must have returned supported code, but was: %v", GetSystemSupportedLangCode("jp"))
	}

	if GetSystemSupportedLangCode("") != "en" {
		t.Fatalf("Must have returned default for empty code, but was: %v", GetSystemSupportedLangCode("jp"))
	}
}

func setUserLang(t *testing.T, user *model.User, code string) {
	user.Language.SetValid(code)

	if err := model.SaveUser(nil, user, false); err != nil {
		t.Fatal(err)
	}
}

func createTestLang(t *testing.T, orga *model.Organization, name, code string, def bool) *model.Language {
	lang := &model.Language{OrganizationId: orga.ID, Name: name, Code: code, Default: def}

	if err := model.SaveLanguage(nil, lang); err != nil {
		t.Fatal(err)
	}

	return lang
}
