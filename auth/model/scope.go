package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type Scope struct {
	db.BaseEntity

	ClientId hide.ID `db:"client_id"`
	Key      string
	Value    string
}

func FindScopeByClientId(clientId hide.ID) []Scope {
	var entities []Scope

	if err := connection.Select(&entities, `SELECT * FROM "scope" WHERE client_id = $1 ORDER BY key ASC`, clientId); err != nil {
		logbuch.Error("Error finding scopes by client id", logbuch.Fields{"err": err, "client_id": clientId})
		return nil
	}

	return entities
}

func SaveScope(tx *sqlx.Tx, entity *Scope) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "scope" (client_id, key, value)
			VALUES (:client_id, :key, :value) RETURNING id`,
		`UPDATE "scope" SET client_id = :client_id,
			key = :key,
			value = :value
			WHERE id = :id`)
}
