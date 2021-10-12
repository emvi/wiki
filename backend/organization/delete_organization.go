package organization

import (
	"emviwiki/backend/billing"
	"emviwiki/backend/client"
	"emviwiki/backend/content"
	"emviwiki/shared/auth"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"path/filepath"

	"emviwiki/backend/errs"
	"emviwiki/shared/model"
)

func DeleteOrganization(orgaId, userId hide.ID, name string, auth auth.AuthClient) error {
	organization := model.GetOrganizationByUserIdAndIdAndIsAdmin(orgaId, userId)

	if organization == nil {
		return errs.PermissionDenied
	}

	if organization.Name != name {
		return errs.NameDoesNotMatch
	}

	// won't do anything if no subscription exists, but makes sure the customer is deleted
	if err := billing.DeleteSubscriptionAndCustomer(organization, false); err != nil {
		logbuch.Error("Error cancelling subscription/deleting customer while deleting organization", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return err
	}

	if err := deleteOrganizationPicture(organization, userId); err != nil {
		logbuch.Error("Error deleting organization picture while deleting organization", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
	}

	// read files before deleting organization to enable rollback
	files := model.FindFileByOrganizationId(orgaId)

	if err := deleteClients(organization, userId, auth); err != nil {
		logbuch.Error("Error deleting organization clients while deleting organization", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return errs.Saving
	}

	if err := model.DeleteOrganizationById(nil, orgaId); err != nil {
		logbuch.Error("Error deleting organization", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return errs.Saving
	}

	deleteOrganizationFiles(organization, userId, files)
	return nil
}

func deleteClients(orga *model.Organization, userId hide.ID, auth auth.AuthClient) error {
	clients := model.FindClientByOrganizationId(orga.ID)

	for _, c := range clients {
		if err := client.DeleteClient(orga, userId, c.ID, auth); err != nil {
			return err
		}
	}

	return nil
}

func deleteOrganizationFiles(orga *model.Organization, userId hide.ID, files []model.File) {
	go func() {
		for _, file := range files {
			content.DeleteFileInStore(orga.ID, userId, filepath.Join(file.Path, file.UniqueName))
		}
	}()

	go func() {
		if orga.Picture.Valid {
			if err := content.DeleteFile(orga, userId, orga.Picture.String); err != nil {
				logbuch.Error("Error deleting organization picture", logbuch.Fields{"err": err})
			}
		}
	}()
}
