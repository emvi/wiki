package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type Login struct {
	db.BaseEntity

	UserId hide.ID `db:"user_id"`
}

func FindLoginByUserId(userId hide.ID) []Login {
	var entities []Login

	if err := connection.Select(&entities, `SELECT * FROM "login" WHERE user_id = $1`, userId); err != nil {
		logbuch.Error("Error finding login by user id", logbuch.Fields{"err": err, "user_id": userId})
		return nil
	}

	return entities
}

func SaveLogin(tx *sqlx.Tx, entity *Login) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "login" (user_id) VALUES (:user_id) RETURNING id`,
		`UPDATE "login" SET user_id = :user_id WHERE id = :id`)
}
