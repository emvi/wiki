package client

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	nameMaxLen = 40
)

type SaveClientData struct {
	Id     hide.ID `json:"id"`
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
}

func (data *SaveClientData) validate(orgaId hide.ID) []error {
	data.Name = strings.TrimSpace(data.Name)
	err := make([]error, 0)

	if len(data.Name) == 0 {
		err = append(err, errs.NameEmpty)
	} else if utf8.RuneCountInString(data.Name) > nameMaxLen {
		err = append(err, errs.NameLen)
	} else if model.GetClientByOrganizationIdAndName(orgaId, data.Name) != nil {
		err = append(err, errs.ClientExistsAlready)
	}

	newScopes := make([]Scope, 0, len(data.Scopes))

	for i := range data.Scopes {
		if !data.Scopes[i].Read && !data.Scopes[i].Write {
			continue
		}

		data.Scopes[i].Name = strings.ToLower(data.Scopes[i].Name)
		found := false

		for _, s := range newScopes {
			if s.Name == data.Scopes[i].Name {
				found = true
			}

			data.Scopes[i].Read = data.Scopes[i].Read || data.Scopes[i].Write
		}

		if found {
			continue
		}

		scope, ok := Scopes[data.Scopes[i].Name]

		if !ok {
			err = append(err, errs.ScopeInvalid)
			break
		}

		if (data.Scopes[i].Read && !scope.Read) || (data.Scopes[i].Write && !scope.Write) {
			err = append(err, errs.ScopeInvalid)
			break
		}

		newScopes = append(newScopes, data.Scopes[i])
	}

	data.Scopes = newScopes

	if len(err) == 0 {
		return nil
	}

	return err
}

func SaveClient(orga *model.Organization, userId hide.ID, data *SaveClientData, auth auth.AuthClient) []error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return []error{err}
	}

	if err := data.validate(orga.ID); err != nil {
		return err
	}

	if data.Id != 0 {
		// updating the name does not affect the name on auth,
		// because the name on auth is different anyway to identify it for this organization
		if err := updateClient(orga.ID, data.Id, data.Name); err != nil {
			return []error{err}
		}
	} else {
		if err := newClient(orga.ID, data, auth); err != nil {
			return []error{err}
		}
	}

	return nil
}

func updateClient(orgaId, id hide.ID, name string) error {
	client := model.GetClientByOrganizationIdAndId(orgaId, id)

	if client == nil {
		return errs.ClientNotFound
	}

	client.Name = name

	if err := model.SaveClient(nil, client); err != nil {
		logbuch.Error("Error saving client while updating", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return errs.Saving
	}

	return nil
}

func newClient(orgaId hide.ID, data *SaveClientData, auth auth.AuthClient) error {
	clientId, clientSecret, err := newClientCredentials(orgaId, data, auth)

	if err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to save new client", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	client := &model.Client{OrganizationId: orgaId,
		Name:         data.Name,
		ClientId:     clientId,
		ClientSecret: clientSecret}

	if err := model.SaveClient(tx, client); err != nil {
		logbuch.Error("Error saving client while creating new client", logbuch.Fields{"err": err, "orga_id": orgaId})
		return errs.Saving
	}

	for _, s := range data.Scopes {
		scope := &model.ClientScope{ClientId: client.ID,
			Name:  s.Name,
			Read:  s.Read,
			Write: s.Write}

		if err := model.SaveClientScope(tx, scope); err != nil {
			logbuch.Error("Error saving client while creating new client", logbuch.Fields{"err": err, "orga_id": orgaId})
			return errs.Saving
		}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction while saving new client", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func newClientCredentials(orgaId hide.ID, data *SaveClientData, auth auth.AuthClient) (string, string, error) {
	scopes := make(map[string]string)

	for _, scope := range data.Scopes {
		scopes[scope.Name] = "r"

		if scope.Write {
			scopes[scope.Name] += "w"
		}
	}

	resp, err := auth.NewClient(strconv.Itoa(int(orgaId))+"_"+data.Name, scopes)

	if err != nil {
		logbuch.Error("Error creating new client on auth", logbuch.Fields{"err": err})
		return "", "", errs.Saving
	}

	return resp.ClientId, resp.ClientSecret, nil
}
