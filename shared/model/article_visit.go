package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"time"
)

type ArticleVisit struct {
	db.BaseEntity

	ArticleId hide.ID `db:"article_id" json:"article_id"`
	UserId    hide.ID `db:"user_id" json:"user_id"`
}

func GetArticleVisitByArticleIdAndUserIdAndDefTimeAfter(articleId, userId hide.ID, defTime time.Time) *ArticleVisit {
	entity := new(ArticleVisit)

	if err := connection.Get(entity, `SELECT * FROM article_visit WHERE article_id = $1 AND user_id = $2 AND def_time > $3`, articleId, userId, defTime); err != nil {
		logbuch.Debug("Article visit by article id and user id and def time after not found", logbuch.Fields{"err": err, "article_id": articleId, "user_id": userId, "def_time": defTime})
		return nil
	}

	return entity
}

func SaveArticleVisit(tx *sqlx.Tx, entity *ArticleVisit) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_visit" (article_id, user_id)
			VALUES (:article_id, :user_id)
			RETURNING id`,
		`UPDATE "article_visit" SET article_id = :article_id,
			user_id = :user_id
			WHERE id = :id`)
}
