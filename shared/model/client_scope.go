package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type ClientScope struct {
	db.BaseEntity

	ClientId hide.ID `db:"client_id" json:"client_id"`
	Name     string  `json:"name"`
	Read     bool    `json:"read"`
	Write    bool    `json:"write"`
}

func FindClientScopeByClientId(clientId hide.ID) []ClientScope {
	query := `SELECT * FROM "client_scope" WHERE client_id = $1`
	var entities []ClientScope

	if err := connection.Select(&entities, query, clientId); err != nil {
		logbuch.Error("Error reading client scopes by client id", logbuch.Fields{"err": err, "client_id": clientId})
		return nil
	}

	return entities
}

func SaveClientScope(tx *sqlx.Tx, entity *ClientScope) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "client_scope" (client_id, name, read, write)
			VALUES (:client_id, :name, :read, :write) RETURNING id`,
		`UPDATE "client_scope" SET client_id = :client_id,
			name = :name,
			read = :read,
			write = :write
			WHERE id = :id`)
}
