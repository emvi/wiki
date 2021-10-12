package lang

import (
	"emviwiki/shared/model"
)

func GetLangs(organization *model.Organization) []model.Language {
	return model.FindLanguagesByOrganizationId(organization.ID)
}
