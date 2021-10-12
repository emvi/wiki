package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type EmailBlacklist struct {
	db.BaseEntity

	Domain string
}

func GetEmailBlacklistByDomain(domain string) *EmailBlacklist {
	entity := new(EmailBlacklist)

	if err := connection.Get(entity, `SELECT * FROM "email_blacklist" WHERE LOWER(domain) = LOWER($1)`, domain); err != nil {
		logbuch.Debug("Email blacklist by domain not found", logbuch.Fields{"err": err, "domain": domain})
		return nil
	}

	return entity
}

func SaveEmailBlacklist(tx *sqlx.Tx, entity *EmailBlacklist) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "email_blacklist" ("domain") VALUES (:domain) RETURNING id`,
		`UPDATE "email_blacklist" SET "domain" = :domain
			WHERE id = :id`)
}
