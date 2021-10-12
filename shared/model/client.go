package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	db.BaseEntity

	OrganizationId hide.ID `db:"organization_id" json:"organization_id"`
	Name           string  `json:"name"`
	ClientId       string  `db:"client_id" json:"client_id"`
	ClientSecret   string  `db:"client_secret" json:"client_secret"`

	Scopes []ClientScope `json:"scopes" db:"-"`
}

func GetClientByOrganizationIdAndId(orgaId, id hide.ID) *Client {
	entity := new(Client)

	if err := connection.Get(entity, `SELECT * FROM "client" WHERE organization_id = $1 AND id = $2`, orgaId, id); err != nil {
		logbuch.Debug("Client by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetClientByOrganizationIdAndName(orgaId hide.ID, name string) *Client {
	entity := new(Client)

	if err := connection.Get(entity, `SELECT * FROM "client" WHERE organization_id = $1 AND LOWER(name) = LOWER($2)`, orgaId, name); err != nil {
		logbuch.Debug("Client by organization id and name not found", logbuch.Fields{"err": err, "orga_id": orgaId, "name": name})
		return nil
	}

	return entity
}

func GetClientByOrganizationIdAndClientId(orgaId hide.ID, clientId string) *Client {
	entity := new(Client)

	if err := connection.Get(entity, `SELECT * FROM "client" WHERE organization_id = $1 AND client_id = $2`, orgaId, clientId); err != nil {
		logbuch.Debug("Client by organization id and client id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "client_id": clientId})
		return nil
	}

	return entity
}

func FindClientByOrganizationId(orgaId hide.ID) []Client {
	query := `SELECT * FROM "client" WHERE organization_id = $1`
	var entities []Client

	if err := connection.Select(&entities, query, orgaId); err != nil {
		logbuch.Error("Error reading clients by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return nil
	}

	return entities
}

func DeleteClientById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "client_scope" WHERE client_id = $1`, id)

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
		`INSERT INTO "client" (organization_id, name, client_id, client_secret)
			VALUES (:organization_id, :name, :client_id, :client_secret) RETURNING id`,
		`UPDATE "client" SET organization_id = :organization_id,
			name = :name,
			client_id = :client_id,
			client_secret = :client_secret
			WHERE id = :id`)
}
