package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type DomainBlacklist struct {
	db.BaseEntity

	Name string
}

func GetDomainBlacklistByName(name string) *DomainBlacklist {
	entity := new(DomainBlacklist)

	if err := connection.Get(entity, `SELECT * FROM "domain_blacklist" WHERE LOWER(name) = LOWER($1)`, name); err != nil {
		logbuch.Debug("Domain blacklist by name not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return entity
}

func SaveDomainBlacklist(tx *sqlx.Tx, entity *DomainBlacklist) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "domain_blacklist" (name)
			VALUES (:name) RETURNING id`,
		`UPDATE "domain_blacklist" SET 	name = :name,
			WHERE id = :id`)
}
