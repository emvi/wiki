package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	db.BaseEntity

	Name         string      `json:"-"`
	ClientId     string      `db:"client_id" json:"client_id"`
	ClientSecret string      `db:"client_secret" json:"client_secret"`
	RedirectURI  null.String `db:"redirect_uri" json:"-"`
	Trusted      bool        `json:"-"`
}

func GetClientByClientIdAndRedirectURI(clientId, redirectURI string) *Client {
	client := new(Client)

	if err := connection.Get(client, `SELECT * FROM "client" WHERE client_id = $1 AND redirect_uri = $2`, clientId, redirectURI); err != nil {
		logbuch.Debug("Client by client id and redirect uri not found", logbuch.Fields{"err": err, "clientId": clientId, "redirectURI": redirectURI})
		return nil
	}

	return client
}

func GetClientByClientIdAndClientSecret(clientId, clientSecret string) *Client {
	client := new(Client)

	if err := connection.Get(client, `SELECT * FROM "client" WHERE client_id = $1 AND client_secret = $2`, clientId, clientSecret); err != nil {
		logbuch.Debug("Client by client id and secret not found", logbuch.Fields{"err": err, "clientId": clientId, "clientSecret": clientSecret})
		return nil
	}

	return client
}

func GetClientByClientIdAndTrusted(clientId string) *Client {
	client := new(Client)

	if err := connection.Get(client, `SELECT * FROM "client" WHERE client_id = $1 AND trusted IS TRUE`, clientId); err != nil {
		logbuch.Debug("Client by client id and trusted not found", logbuch.Fields{"err": err, "clientId": clientId})
		return nil
	}

	return client
}

func GetClientByClientId(clientId string) *Client {
	return GetClientByClientIdTx(nil, clientId)
}

func GetClientByClientIdTx(tx *sqlx.Tx, clientId string) *Client {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	client := new(Client)

	if err := tx.Get(client, `SELECT * FROM "client" WHERE client_id = $1`, clientId); err != nil {
		logbuch.Debug("Client by client id not found", logbuch.Fields{"err": err, "clientId": clientId})
		return nil
	}

	return client
}

func GetClientByClientSecretTx(tx *sqlx.Tx, secret string) *Client {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	client := new(Client)

	if err := tx.Get(client, `SELECT * FROM "client" WHERE client_secret = $1`, secret); err != nil {
		logbuch.Debug("Client by client secret not found", logbuch.Fields{"err": err, "secret": secret})
		return nil
	}

	return client
}

func GetClientByName(name string) *Client {
	client := new(Client)

	if err := connection.Get(client, `SELECT * FROM "client" WHERE LOWER("name") = LOWER($1)`, name); err != nil {
		logbuch.Debug("Client by name not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return client
}

func DeleteClientById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "scope" WHERE client_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting client scope by client id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "client" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting client by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveClient(tx *sqlx.Tx, entity *Client) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "client" ("name", client_id, client_secret, redirect_uri, trusted)
			VALUES (:name, :client_id, :client_secret, :redirect_uri, :trusted) RETURNING id`,
		`UPDATE "client" SET "name" = :name,
			client_id = :client_id,
			client_secret = :client_secret,
			redirect_uri = :redirect_uri,
			trusted = :trusted
			WHERE id = :id`)
}
