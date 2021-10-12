package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type ArticleTag struct {
	db.BaseEntity

	ArticleId hide.ID `db:"article_id"`
	TagId     hide.ID `db:"tag_id"`
}

func GetArticleTagByOrganizationIdAndArticleIdAndName(orgaId, articleId hide.ID, name string) *ArticleTag {
	entity := new(ArticleTag)

	if err := connection.Get(entity, `SELECT "article_tag".* FROM "article_tag"
		JOIN "tag" ON "article_tag".tag_id = "tag".id
		WHERE organization_id = $1
		AND article_id = $2
		AND LOWER("name") = LOWER($3)`, orgaId, articleId, name); err != nil {
		logbuch.Debug("Article tag by organization id and article id and name not found", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId, "name": name})
		return nil
	}

	return entity
}

func CountArticleTagByArticleId(articleId hide.ID) int {
	query := `SELECT COUNT(1) FROM "article_tag" WHERE article_id = $1`
	var count int

	if err := connection.Get(&count, query, articleId); err != nil {
		logbuch.Error("Error counting tags by article id", logbuch.Fields{"err": err, "article_id": articleId})
		return 0
	}

	return count
}

func DeleteArticleTagByOrganizationIdAndTagId(tx *sqlx.Tx, orgaId, tagId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "article_tag"
		WHERE tag_id = $2
		AND EXISTS (SELECT 1 FROM "tag" WHERE id = "article_tag".tag_id AND organization_id = $1)`, orgaId, tagId)

	if err != nil {
		logbuch.Error("Error deleting article tag by organization id and tag id", logbuch.Fields{"err": err, "orga_id": orgaId, "tag_id": tagId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func DeleteArticleTagByOrganizationIdAndId(tx *sqlx.Tx, orgaId, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "article_tag"
		WHERE id = $2
		AND EXISTS (SELECT 1 FROM "tag" WHERE id = "article_tag".tag_id AND organization_id = $1)`, orgaId, id)

	if err != nil {
		logbuch.Error("Error deleting article tag by organization id and id", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleTag(tx *sqlx.Tx, entity *ArticleTag) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_tag" (article_id, tag_id)
			VALUES (:article_id, :tag_id) RETURNING id`,
		`UPDATE "article_tag" SET article_id = :article_id,
			tag_id = :tag_id
			WHERE id = :id`)
}
