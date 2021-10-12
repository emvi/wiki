package lang

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	iso6391 "github.com/emvi/iso-639-1"
	"strings"
)

func AddLang(orga *model.Organization, userId hide.ID, code string) error {
	if !orga.Expert {
		return errs.RequiresExpertVersion
	}

	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	code = strings.TrimSpace(code)

	if !iso6391.ValidCode(code) {
		return errs.LanguageInvalid
	}

	if model.GetLanguageByOrganizationIdAndCode(orga.ID, code) != nil {
		return errs.LanguageExists
	}

	isoLang := iso6391.FromCode(code)
	lang := &model.Language{OrganizationId: orga.ID,
		Name:    isoLang.NativeName,
		Code:    isoLang.Code,
		Default: false}

	if err := model.SaveLanguage(nil, lang); err != nil {
		return errs.Saving
	}

	return nil
}
