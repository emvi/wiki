package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type Language struct {
	db.BaseEntity

	OrganizationId hide.ID `db:"organization_id" json:"organization_id"`
	Name           string  `json:"name"`
	Code           string  `json:"code"`
	Default        bool    `json:"default"`
}

func GetLanguageByOrganizationIdAndId(orgaId, id hide.ID) *Language {
	return GetLanguageByOrganizationIdAndIdTx(nil, orgaId, id)
}

func GetLanguageByOrganizationIdAndIdTx(tx *sqlx.Tx, orgaId, id hide.ID) *Language {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(Language)

	if err := tx.Get(entity, `SELECT * FROM "language" WHERE organization_id = $1 AND id = $2`, orgaId, id); err != nil {
		logbuch.Debug("Language by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetLanguageByOrganizationIdAndCode(id hide.ID, code string) *Language {
	return GetLanguageByOrganizationIdAndCodeTx(nil, id, code)
}

func GetLanguageByOrganizationIdAndCodeTx(tx *sqlx.Tx, id hide.ID, code string) *Language {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(Language)

	if err := tx.Get(entity, `SELECT * FROM "language" WHERE organization_id = $1 AND LOWER(code) = LOWER($2)`, id, code); err != nil {
		// don't log here, this gets called too often
		return nil
	}

	return entity
}

func GetDefaultLanguageByOrganizationId(id hide.ID) *Language {
	return GetDefaultLanguageByOrganizationIdTx(nil, id)
}

func GetDefaultLanguageByOrganizationIdTx(tx *sqlx.Tx, id hide.ID) *Language {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(Language)

	if err := tx.Get(entity, `SELECT * FROM "language" WHERE organization_id = $1 AND "default" IS TRUE`, id); err != nil {
		logbuch.Debug("Default language by organization id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func FindLanguagesByOrganizationId(id hide.ID) []Language {
	var entities []Language

	if err := connection.Select(&entities, `SELECT * FROM "language" WHERE organization_id = $1 ORDER BY name ASC`, id); err != nil {
		logbuch.Error("Error reading languages by organization id", logbuch.Fields{"err": err, "organization_id": id})
		return nil
	}

	return entities
}

func SaveLanguage(tx *sqlx.Tx, entity *Language) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "language" (organization_id,
			name,
			code,
			"default")
			VALUES (:organization_id,
			:name,
			:code,
			:default) RETURNING id`,
		`UPDATE "language" SET organization_id = :organization_id,
			name = :name,
			code = :code,
			"default" = :default
			WHERE id = :id`)
}
