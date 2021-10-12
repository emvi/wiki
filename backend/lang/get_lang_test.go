package lang

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestGetLang(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)

	if _, err := GetLang(orga, 0); err != errs.LanguageNotFound {
		t.Fatal("Language must not be found")
	}

	en := &model.Language{Code: "en", Name: "English", OrganizationId: orga.ID}

	if err := model.SaveLanguage(nil, en); err != nil {
		t.Fatal(err)
	}

	if lang, err := GetLang(orga, en.ID); err != nil || lang.Code != "en" || lang.Name != "English" {
		t.Fatal("Language en must be returned")
	}
}
