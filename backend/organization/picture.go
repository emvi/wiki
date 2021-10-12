package organization

import (
	"emviwiki/backend/content"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"net/http"
)

const (
	organizationPicturePath = "organization/pictures"
)

func UploadOrganizationPicture(r *http.Request, orgaId, userId hide.ID) error {
	// read organization
	organization := model.GetOrganizationByUserIdAndIdAndIsAdmin(orgaId, userId)

	if organization == nil {
		return errs.PermissionDenied
	}

	if err := deleteOrganizationPicture(organization, userId); err != nil {
		logbuch.Debug("Error deleting old organization picture", logbuch.Fields{"err": err})
	}

	uniqueName, err := content.UploadFile(&content.File{
		Request:       r,
		Organization:  organization,
		UserId:        userId,
		Path:          organizationPicturePath,
		RequiresImage: true,
	})

	if err != nil {
		return errs.UploadingFile
	}

	organization.Picture = null.NewString(uniqueName, true)

	if err := model.SaveOrganization(nil, organization); err != nil {
		logbuch.Error("Error saving organization when uploading picture", logbuch.Fields{"err": err, "orga_id": orgaId})
		return errs.Saving
	}

	return nil
}

func DeleteOrganizationPicture(orgaId, userId hide.ID) error {
	organization := model.GetOrganizationByUserIdAndIdAndIsAdmin(orgaId, userId)

	if organization == nil {
		return errs.PermissionDenied
	}

	if err := deleteOrganizationPicture(organization, userId); err != nil {
		return err
	}

	return nil
}

func deleteOrganizationPicture(organization *model.Organization, userId hide.ID) error {
	if organization.Picture.Valid {
		if err := content.DeleteFile(organization, userId, organization.Picture.String); err != nil {
			logbuch.Error("Error deleting organization picture", logbuch.Fields{"err": err})
		}

		// save organization
		organization.Picture = null.String{}

		if err := model.SaveOrganization(nil, organization); err != nil {
			logbuch.Error("Error saving organization when deleting picture", logbuch.Fields{"err": err, "orga_id": organization.ID})
			return errs.Saving
		}
	}

	return nil
}
