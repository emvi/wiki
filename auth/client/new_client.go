package client

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
	"unicode/utf8"
)

const (
	nameMaxLen       = 100
	scopeKeyMaxLen   = 40
	scopeValueMaxLen = 100
	clientIdLen      = 20
	clientSecretLen  = 64
)

// NewClient creates a new client. This must only be called by trusted clients (check on API layer)!
func NewClient(name string, scopes map[string]string) (*model.Client, error) {
	// validate name and scopes
	name = strings.TrimSpace(name)

	if len(name) == 0 || utf8.RuneCountInString(name) > nameMaxLen {
		return nil, errs.ClientNameInvalid
	}

	if model.GetClientByName(name) != nil {
		return nil, errs.ClientNameInUse
	}

	for key, value := range scopes {
		if utf8.RuneCountInString(key) > scopeKeyMaxLen || utf8.RuneCountInString(value) > scopeValueMaxLen {
			return nil, errs.ScopeInvalid
		}
	}

	// create new client
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to create new client", logbuch.Fields{"err": err})
		return nil, errs.TxBegin
	}

	client := &model.Client{Name: name,
		ClientId:     newUniqueClientId(tx),
		ClientSecret: newUniqueClientSecret(tx)}

	if err := model.SaveClient(tx, client); err != nil {
		logbuch.Error("Error saving client while creating new client", logbuch.Fields{"err": err})
		return nil, errs.Saving
	}

	for key, value := range scopes {
		scope := &model.Scope{ClientId: client.ID,
			Key:   key,
			Value: value}

		if err := model.SaveScope(tx, scope); err != nil {
			logbuch.Error("Error saving scope while creating new client", logbuch.Fields{"err": err})
			return nil, errs.Saving
		}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction while creating new client", logbuch.Fields{"err": err})
		return nil, errs.TxCommit
	}

	return client, nil
}

func newUniqueClientId(tx *sqlx.Tx) string {
	id := util.GenRandomString(clientIdLen)

	for model.GetClientByClientIdTx(tx, id) != nil {
		id = util.GenRandomString(clientIdLen)
	}

	return id
}

func newUniqueClientSecret(tx *sqlx.Tx) string {
	secret := util.GenRandomString(clientSecretLen)

	for model.GetClientByClientSecretTx(tx, secret) != nil {
		secret = util.GenRandomString(clientSecretLen)
	}

	return secret
}
