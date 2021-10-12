package lang

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func GetLang(organization *model.Organization, id hide.ID) (*model.Language, error) {
	lang := model.GetLanguageByOrganizationIdAndId(organization.ID, id)

	if lang == nil {
		return nil, errs.LanguageNotFound
	}

	return lang, nil
}
