package tag

import (
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func GetTagByIdOrName(orga *model.Organization, userId, id hide.ID, name string) *model.Tag {
	var tag *model.Tag

	if id != 0 {
		tag = model.GetTagByOrganizationIdAndId(orga.ID, id)
	} else {
		tag = model.GetTagByOrganizationIdAndName(orga.ID, name)
	}

	if tag == nil {
		return nil
	}

	if err := perm.CheckUserTagAccess(orga.ID, userId, tag.ID); err != nil {
		return nil
	}

	return tag
}
