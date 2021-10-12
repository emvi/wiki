package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"time"
)

type Newsletter struct {
	db.BaseEntity

	Subject   string    `json:"subject"`
	Content   string    `json:"content"`
	Scheduled time.Time `json:"scheduled"`
	Status    string    `json:"status"`
}

func GetNewsletterById(id hide.ID) *Newsletter {
	query := `SELECT * FROM "newsletter" WHERE id = $1`
	entity := new(Newsletter)

	if err := dashboardDB.Get(entity, query, id); err != nil {
		logbuch.Debug("Newsletter by id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func GetLatestNewsletterToSend() *Newsletter {
	query := `SELECT * FROM "newsletter" WHERE scheduled < NOW() AND status = 'planned' LIMIT 1`
	entity := new(Newsletter)

	if err := dashboardDB.Get(entity, query); err != nil {
		logbuch.Debug("Latest newsletter not found", logbuch.Fields{"err": err})
		return nil
	}

	return entity
}

func FindNewsletter() []Newsletter {
	query := `SELECT * FROM "newsletter" ORDER BY def_time DESC`
	var entities []Newsletter

	if err := dashboardDB.Select(&entities, query); err != nil {
		logbuch.Error("Error finding newsletter", logbuch.Fields{"err": err})
		return nil
	}

	return entities
}

func SaveNewsletter(tx *sqlx.Tx, entity *Newsletter) error {
	return dashboardDB.SaveEntity(tx, entity,
		`INSERT INTO "newsletter" (subject,
			content,
			scheduled,
			status) VALUES (:subject,
			:content,
			:scheduled,
			:status) RETURNING id`,
		`UPDATE "newsletter" SET subject = :subject,
			content = :content,
			scheduled = :scheduled,
			status = :status
			WHERE id = :id`)
}

func DeleteNewsletterById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		var err error
		tx, err = dashboardDB.Beginx()

		if err != nil {
			return err
		}

		defer func() {
			if err := tx.Commit(); err != nil {
				logbuch.Error("Error committing transaction to delete newsletter", logbuch.Fields{"err": err, "id": id})
			}
		}()
	}

	if _, err := tx.Exec(`DELETE FROM "newsletter" WHERE id = $1`, id); err != nil {
		logbuch.Error("Error deleting newsletter by id", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
