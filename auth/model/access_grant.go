package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type AccessGrant struct {
	db.BaseEntity

	UserId   hide.ID `db:"user_id"`
	ClientId hide.ID `db:"client_id"`
}

func GetAccessGrantByUserIdAndClientId(userId, clientId hide.ID) *AccessGrant {
	access := new(AccessGrant)

	if err := connection.Get(access, `SELECT * FROM "access_grant" WHERE user_id = $1 AND client_id = $2`, userId, clientId); err != nil {
		logbuch.Debug("Access grant by user ID and client ID not found", logbuch.Fields{"err": err, "user_id": userId, "client_id": clientId})
		return nil
	}

	return access
}

func SaveAccessGrant(tx *sqlx.Tx, entity *AccessGrant) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "access_grant" (user_id, client_id)
			VALUES (:user_id, :client_id) RETURNING id`,
		`UPDATE "access_grant" SET user_id = :user_id,
			client_id = :client_id
			WHERE id = :id`)
}
