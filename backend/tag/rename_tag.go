package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func RenameTag(orga *model.Organization, userId, id hide.ID, name string) error {
	if _, err := perm.CheckUserIsAdminOrMod(orga.ID, userId); err != nil {
		return err
	}

	if err := perm.CheckUserTagAccess(orga.ID, userId, id); err != nil {
		return err
	}

	data := AddTagData{Tag: name}

	if err := data.validate(); err != nil {
		return err
	}

	tag := model.GetTagByOrganizationIdAndId(orga.ID, id)

	if tag == nil {
		return errs.TagNotFound
	}

	if existing := model.GetTagByOrganizationIdAndName(orga.ID, data.Tag); existing != nil && existing.ID != tag.ID {
		return errs.TagNameExistsAlready
	}

	tag.Name = data.Tag

	if err := model.SaveTag(nil, tag); err != nil {
		logbuch.Error("Error saving tag while renaming", logbuch.Fields{"err": err, "orga_id": orga.ID, "id": id, "name": name})
		return errs.Saving
	}

	return nil
}
