package client

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func DeleteClient(orga *model.Organization, userId, id hide.ID, auth auth.AuthClient) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	client := model.GetClientByOrganizationIdAndId(orga.ID, id)

	if client == nil {
		return errs.ClientNotFound
	}

	if err := auth.DeleteClient(client.ClientId, client.ClientSecret); err != nil {
		logbuch.Error("Error deleting client on auth", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "id": id})
		return errs.Saving
	}

	if err := model.DeleteClientById(nil, id); err != nil {
		return errs.Saving
	}

	return nil
}
