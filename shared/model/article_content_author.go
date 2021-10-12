package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type ArticleContentAuthor struct {
	db.BaseEntity

	ArticleContentId hide.ID `db:"article_content_id"`
	UserId           hide.ID `db:"user_id"`
}

func FindArticleContentAuthorByArticleContentId(id hide.ID) []ArticleContentAuthor {
	return FindArticleContentAuthorByArticleContentIdTx(nil, id)
}

func FindArticleContentAuthorByArticleContentIdTx(tx *sqlx.Tx, id hide.ID) []ArticleContentAuthor {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var entities []ArticleContentAuthor

	if err := tx.Select(&entities, `SELECT * FROM "article_content_author" WHERE article_content_id = $1`, id); err != nil {
		logbuch.Error("Error finding article content authors by article content id", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entities
}

func FindArticleContentAuthorUserByArticleId(id hide.ID) []User {
	return FindArticleContentAuthorUserByArticleIdTx(nil, id)
}

func FindArticleContentAuthorUserByArticleIdTx(tx *sqlx.Tx, id hide.ID) []User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT DISTINCT ON (id) * FROM (SELECT
		` + userBaseQueryFields + ` FROM "article_content"
		JOIN "article_content_author" ON "article_content".id = "article_content_author".article_content_id
		JOIN "user" ON "article_content_author".user_id = "user".id
		JOIN "organization_member" ON "user".id = "organization_member".user_id
		WHERE article_id = $1
		ORDER BY "user".lastname, "user".firstname, "organization_member".username ASC
		) AS results`
	var entities []User

	if err := tx.Select(&entities, query, id); err != nil {
		logbuch.Error("Error finding article content author user by article content id", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entities
}

func SaveArticleContentAuthor(tx *sqlx.Tx, entity *ArticleContentAuthor) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_content_author" (article_content_id, user_id)
			VALUES (:article_content_id, :user_id)
			RETURNING id`,
		`UPDATE "article_content_author" SET article_content_id = :article_content_id,
			user_id = :user_id
			WHERE id = :id`)
}
