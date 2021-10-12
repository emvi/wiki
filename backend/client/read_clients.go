package client

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func ReadClients(orga *model.Organization, userId hide.ID) ([]model.Client, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return nil, err
	}

	clients := model.FindClientByOrganizationId(orga.ID)

	for i := range clients {
		clients[i].Scopes = model.FindClientScopeByClientId(clients[i].ID)
	}

	return clients, nil
}

func ReadClient(orga *model.Organization, userId, id hide.ID) (*model.Client, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return nil, err
	}

	client := model.GetClientByOrganizationIdAndId(orga.ID, id)

	if client == nil {
		return nil, errs.ClientNotFound
	}

	client.Scopes = model.FindClientScopeByClientId(client.ID)
	return client, nil
}
