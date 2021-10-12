package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

type ArticleListName struct {
	db.BaseEntity

	ArticleListId hide.ID     `db:"article_list_id" json:"article_list_id"`
	LanguageId    hide.ID     `db:"language_id" json:"language_id"`
	Name          string      `json:"name"`
	Info          null.String `json:"info"`
}

func GetArticleListNameByOrganizationIdAndArticleListIdAndLangId(orgaId, listId, langId hide.ID) *ArticleListName {
	return GetArticleListNameByOrganizationIdAndArticleListIdAndLangIdTx(nil, orgaId, listId, langId)
}

func GetArticleListNameByOrganizationIdAndArticleListIdAndLangIdTx(tx *sqlx.Tx, orgaId, listId, langId hide.ID) *ArticleListName {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(ArticleListName)

	if err := tx.Get(entity, `SELECT * FROM "article_list_name" WHERE article_list_id = $2 AND language_id = $3
		UNION ALL
		SELECT * FROM "article_list_name" WHERE article_list_id = $2 AND language_id = (SELECT id FROM "language" WHERE organization_id = $1 AND "default" IS TRUE)
		UNION ALL
		SELECT * FROM "article_list_name" WHERE article_list_id = $2
		LIMIT 1`, orgaId, listId, langId); err != nil {
		logbuch.Debug("Error reading article list name by organization id and article list id and language id", logbuch.Fields{"err": err, "orga_id": orgaId, "list_id": listId, "lang_id": langId})
		return nil
	}

	return entity
}

func FindArticleListNamesByArticleListId(listId hide.ID) []ArticleListName {
	var entities []ArticleListName

	if err := connection.Select(&entities, `SELECT * FROM "article_list_name"
		WHERE article_list_id = $1`, listId); err != nil {
		logbuch.Error("Error reading article list names by article list id", logbuch.Fields{"err": err, "list_id": listId})
		return nil
	}

	return entities
}

func DeleteArticleListNameByArticleListId(tx *sqlx.Tx, listId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "article_list_name" WHERE article_list_id = $1`, listId)

	if err != nil {
		logbuch.Error("Error deleting article list name by article list id", logbuch.Fields{"err": err, "list_id": listId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleListName(tx *sqlx.Tx, entity *ArticleListName) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_list_name" (article_list_id,
			language_id,
			name,
			info)
			VALUES (:article_list_id,
			:language_id,
			:name,
			:info) RETURNING id`,
		`UPDATE "article_list_name" SET article_list_id = :article_list_id,
			language_id = :language_id,
			name = :name,
			info = :info
			WHERE id = :id`)
}
