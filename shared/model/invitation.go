package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type Invitation struct {
	db.BaseEntity

	OrganizationId hide.ID `db:"organization_id" json:"organization_id"`
	Email          string  `json:"email"`
	Code           string  `json:"code"`
	ReadOnly       bool    `db:"read_only" json:"read_only"`

	OrganizationName string `db:"organization_name" json:"organization_name"`
}

func GetInvitationByOrganizationIdAndEmail(orgaId hide.ID, email string) *Invitation {
	entity := new(Invitation)

	if err := connection.Get(entity, `SELECT * FROM "invitation" WHERE organization_id = $1 AND LOWER(email) = LOWER($2)`, orgaId, email); err != nil {
		logbuch.Debug("Invitation by organization id an email not found", logbuch.Fields{"err": err, "orga_id": orgaId, "email": email})
		return nil
	}

	return entity
}

func GetInvitationByEmailAndCode(email, code string) *Invitation {
	entity := new(Invitation)

	if err := connection.Get(entity, `SELECT * FROM "invitation" WHERE LOWER(email) = LOWER($1) AND code = $2`, email, code); err != nil {
		logbuch.Debug("Invitation by email and code not found", logbuch.Fields{"err": err, "email": email, "code": code})
		return nil
	}

	return entity
}

func GetInvitationByEmailAndId(email string, id hide.ID) *Invitation {
	entity := new(Invitation)

	if err := connection.Get(entity, `SELECT * FROM "invitation" WHERE LOWER(email) = LOWER($1) AND id = $2`, email, id); err != nil {
		logbuch.Debug("Invitation by email and id not found", logbuch.Fields{"err": err, "email": email, "id": id})
		return nil
	}

	return entity
}

func GetInvitationByOrganizationIdAndId(orgaId, id hide.ID) *Invitation {
	entity := new(Invitation)

	if err := connection.Get(entity, `SELECT * FROM "invitation" WHERE organization_id = $1 AND id = $2`, orgaId, id); err != nil {
		logbuch.Debug("Invitation by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func FindInvitationsByEmail(email string) []Invitation {
	query := `SELECT "invitation".*, "organization"."name" "organization_name"
		FROM "invitation"
		JOIN "organization" ON "invitation".organization_id = "organization".id
		WHERE LOWER("email") = LOWER($1)
		AND "invitation".def_time > NOW() - INTERVAL '30 days'`
	var entities []Invitation

	if err := connection.Select(&entities, query, email); err != nil {
		logbuch.Error("Error reading invitations by email", logbuch.Fields{"err": err, "email": email})
		return nil
	}

	return entities
}

func FindInvitationsByOrganizationId(orgaId hide.ID) []Invitation {
	query := `SELECT * FROM "invitation" WHERE organization_id = $1 ORDER BY email ASC`
	var entities []Invitation

	if err := connection.Select(&entities, query, orgaId); err != nil {
		logbuch.Error("Error reading invitations by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return nil
	}

	return entities
}

func DeleteInvitationByOrganizationIdAndEmail(tx *sqlx.Tx, orgaId hide.ID, email string) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "invitation" WHERE organization_id = $1 AND LOWER(email) = LOWER($2)`, orgaId, email); err != nil {
		logbuch.Error("Error deleting invitation by organization id and email", logbuch.Fields{"err": err, "orga_id": orgaId, "email": email})
		db.Rollback(tx)
		return err
	}

	return nil
}

func DeleteInvitationByDefTimeBeforeOneMonth(tx *sqlx.Tx) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "invitation" WHERE def_time < NOW() - INTERVAL '1 month'`); err != nil {
		logbuch.Error("Error deleting invitation by def time before one month", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	return nil
}

func DeleteInvitationById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "invitation" WHERE id = $1`, id); err != nil {
		logbuch.Error("Error deleting invitation by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveInvitation(tx *sqlx.Tx, entity *Invitation) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "invitation" (organization_id,
			email,
			code,
			read_only)
			VALUES (:organization_id,
			:email,
			:code,
			:read_only) RETURNING id`,
		`UPDATE "invitation" SET organization_id = :organization_id,
			email = :email,
			code = :code,
			read_only = :read_only
			WHERE id = :id`)
}
