package client

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"github.com/emvi/logbuch"
)

// DeleteClient deletes a client with given id and secret.
// This must only be called by trusted clients (check on API layer)!
func DeleteClient(clientId, clientSecret string) error {
	client := model.GetClientByClientIdAndClientSecret(clientId, clientSecret)

	if client == nil {
		return errs.ClientNotFound
	}

	if err := model.DeleteClientById(nil, client.ID); err != nil {
		logbuch.Error("Error deleting client", logbuch.Fields{"err": err, "client_id": clientId, "client_secret": clientSecret})
		return errs.Saving
	}

	return nil
}
