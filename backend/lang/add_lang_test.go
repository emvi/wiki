package lang

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestAddLang(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if err := AddLang(orga, user.ID, "ja"); err != errs.LanguageExists {
		t.Fatalf("Expected language to exist already, but was: %v", err)
	}

	if err := AddLang(orga, user.ID, ""); err != errs.LanguageInvalid {
		t.Fatalf("Expected language to be invalid, but was: %v", err)
	}

	if err := AddLang(orga, user.ID, "en"); err != nil {
		t.Fatalf("Expected new language to be added, but was: %v", err)
	}
}

func TestAddLangNonExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := AddLang(orga, user.ID, "en"); err != errs.RequiresExpertVersion {
		t.Fatalf("Expected requires expert, but was: %v", err)
	}
}
