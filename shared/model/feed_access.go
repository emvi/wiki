package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type FeedAccess struct {
	db.BaseEntity

	UserId       hide.ID `db:"user_id" json:""`
	FeedId       hide.ID `db:"feed_id" json:""`
	Notification bool    `json:"notification"`
	Read         bool    `json:"read"`
}

func CountFeedAccessByOrganizationIdAndUserIdAndNotificationAndRead(orgaId, userId hide.ID, notification, read bool) int {
	query := `SELECT COUNT(1) FROM "feed_access"
		JOIN "feed" ON "feed_access".feed_id = "feed".id AND "feed".organization_id = $1
		WHERE "feed_access".user_id = $2
		AND notification = $3
		AND read = $4 `

	var count int

	if err := connection.Get(&count, query, orgaId, userId, notification, read); err != nil {
		logbuch.Error("Error counting feed access by organization id and user id and notification and read", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "notification": notification, "read": read})
		return 0
	}

	return count
}

func GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotification(orgaId, userId, feedId hide.ID, notification bool) *FeedAccess {
	return GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotificationTx(nil, orgaId, userId, feedId, notification)
}

func GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotificationTx(tx *sqlx.Tx, orgaId, userId, feedId hide.ID, notification bool) *FeedAccess {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(FeedAccess)

	if err := tx.Get(entity, `SELECT "feed_access".* FROM "feed_access"
		JOIN "feed" ON "feed_access".feed_id = "feed".id
		WHERE "feed".organization_id = $1
		AND "feed_access".user_id = $2
		AND "feed".id = $3
		AND notification = $4`, orgaId, userId, feedId, notification); err != nil {
		logbuch.Debug("Feed access by organization id and user id and feed id and notification not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "feed_id": feedId, "notification": notification})
		return nil
	}

	return entity
}

func FindFeedAccessByFeedId(feedId hide.ID) []FeedAccess {
	var entities []FeedAccess

	if err := connection.Select(&entities, `SELECT * FROM "feed_access" WHERE feed_id = $1`, feedId); err != nil {
		logbuch.Error("Error reading feed access by feed id", logbuch.Fields{"err": err, "feed_id": feedId})
		return nil
	}

	return entities
}

func UpdateFeedAccessNotificationByOrganizationIdAndUserId(tx *sqlx.Tx, orgaId, userId hide.ID, read bool) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`UPDATE "feed_access" SET read = $3
		FROM "feed"
		WHERE "feed".id = "feed_access".feed_id
		AND "feed".organization_id = $1
		AND "feed_access".user_id = $2`, orgaId, userId, read)

	if err != nil {
		db.Rollback(tx)
		logbuch.Error("Error updating feed access notification", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "read": read})
		return err
	}

	return nil
}

func DeleteFeedAccessByFeedIds(tx *sqlx.Tx, ids []hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query, args, err := sqlx.In(`DELETE FROM "feed_access" WHERE feed_id IN (?)`, ids)

	if err != nil {
		logbuch.Error("Error deleting feed access by feed ids", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	query = connection.Rebind(query)
	_, err = tx.Exec(query, args...)

	if err != nil {
		logbuch.Error("Error deleting feed access by feed ids", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveFeedAccess(tx *sqlx.Tx, entity *FeedAccess) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "feed_access" (notification,
			read,
			user_id,
			feed_id)
			VALUES (:notification,
			:read,
			:user_id,
			:feed_id) RETURNING id`,
		`UPDATE "feed_access" SET notification = :notification,
			read = :read,
			user_id = :user_id,
			feed_id = :feed_id
			WHERE id = :id`)
}
