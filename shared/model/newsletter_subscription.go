package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

type NewsletterSubscription struct {
	db.BaseEntity

	Email     string
	List      null.String
	Confirmed bool
	Code      string
}

func GetNewsletterSubscriptionByEmailAndList(email, list string) *NewsletterSubscription {
	query := `SELECT * FROM "newsletter_subscription" WHERE LOWER(email) = LOWER($1) AND list = $2`
	params := make([]interface{}, 1)
	params[0] = email

	if list == "" {
		query = `SELECT * FROM "newsletter_subscription" WHERE LOWER(email) = LOWER($1) AND list IS NULL`
	} else {
		params = append(params, list)
	}

	entity := new(NewsletterSubscription)

	if err := connection.Get(entity, query, params...); err != nil {
		logbuch.Debug("Newsletter by email and list not found", logbuch.Fields{"err": err, "email": email, "list": list})
		return nil
	}

	return entity
}

func GetNewsletterSubscriptionByCode(code string) *NewsletterSubscription {
	entity := new(NewsletterSubscription)

	if err := connection.Get(entity, `SELECT * FROM "newsletter_subscription" WHERE code = $1`, code); err != nil {
		logbuch.Debug("Newsletter by code not found", logbuch.Fields{"err": err, "code": code})
		return nil
	}

	return entity
}

func DeleteNewsletterSubscriptionById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "newsletter_subscription" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting newsletter by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveNewsletterSubscription(tx *sqlx.Tx, entity *NewsletterSubscription) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "newsletter_subscription" (email,
			list,
			confirmed,
			code)
			VALUES (:email,
			:list,
			:confirmed,
			:code) RETURNING id`,
		`UPDATE "newsletter_subscription" SET email = :email,
			list = :list,
			confirmed = :confirmed,
			code = :code
			WHERE id = :id`)
}
