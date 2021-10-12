package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

type Bookmark struct {
	db.BaseEntity

	OrganizationId hide.ID    `db:"organization_id" json:"organization_id"`
	UserId         hide.ID    `db:"user_id" json:"user_id"`
	ArticleId      null.Int64 `db:"article_id" json:"article_id"`
	ArticleListId  null.Int64 `db:"article_list_id" json:"article_list_id"`

	Article     *Article     `db:"article" json:"article"`
	ArticleList *ArticleList `db:"article_list" json:"article_list"`
}

func GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orgaId, userId, articleId, listId hide.ID) *Bookmark {
	entity := new(Bookmark)

	if err := connection.Get(entity, `SELECT * FROM "bookmark" WHERE organization_id = $1 AND user_id = $2 AND (article_id = $3 OR article_list_id = $4)`, orgaId, userId, articleId, listId); err != nil {
		logbuch.Debug("Bookmark by organization id and user id and article id or article list id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "article_id": articleId, "list_id": listId})
		return nil
	}

	return entity
}

func GetBookmarkByUserIdAndArticleId(userId, articleId hide.ID) *Bookmark {
	entity := new(Bookmark)

	if err := connection.Get(entity, `SELECT * FROM "bookmark" WHERE user_id = $1 AND article_id = $2`, userId, articleId); err != nil {
		logbuch.Debug("Bookmark by user id and article id not found", logbuch.Fields{"err": err, "user_id": userId, "article_id": articleId})
		return nil
	}

	return entity
}

func GetBookmarkByUserIdAndArticleListId(userId, listId hide.ID) *Bookmark {
	entity := new(Bookmark)

	if err := connection.Get(entity, `SELECT * FROM "bookmark" WHERE user_id = $1 AND article_list_id = $2`, userId, listId); err != nil {
		logbuch.Debug("Bookmark by user id and article list id not found", logbuch.Fields{"err": err, "user_id": userId, "list_id": listId})
		return nil
	}

	return entity
}

func FindBookmarkByOrganizationIdAndUserIdAndLanguageIdArticleIdSetWithLimit(orgaId, userId, langId hide.ID, offset, n int) []Bookmark {
	query := `SELECT "article".id "article.id",
		"article".organization_id "article.organization_id",
		"article".views "article.views",
		"article".wip "article.wip",
		"article".read_everyone "article.read_everyone",
		"article".write_everyone "article.write_everyone",
		"article".client_access "article.client_access",
		"article".archived "article.archived",
		"article".published "article.published"
		FROM "bookmark"
		JOIN "article" ON "bookmark".article_id = "article".id
		WHERE "bookmark".organization_id = $1
		AND user_id = $2
		ORDER BY "bookmark"."def_time" DESC
		LIMIT $3 OFFSET $4`
	var entities []Bookmark

	if err := connection.Select(&entities, query, orgaId, userId, n, offset); err != nil {
		logbuch.Error("Error reading bookmarks by organization id and user id and article id set with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "offset": offset, "n": n})
		return nil
	}

	for i := range entities {
		entities[i].Article.LatestArticleContent = GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageId(orgaId, entities[i].Article.ID, langId, false)
	}

	return entities
}

func FindBookmarkByOrganizationIdAndUserIdAndLanguageIdArticleListIdSetWithLimit(orgaId, userId, langId hide.ID, offset, n int) []Bookmark {
	query := `SELECT "article_list".id "article_list.id",
		"article_list".organization_id "article_list.organization_id",
		"article_list".public "article_list.public",
		"article_list".def_time "article_list.def_time",
		"article_list".mod_time "article_list.mod_time",
		(SELECT COUNT(1) FROM "article_list_entry" WHERE article_list_id = "article_list".id) AS "article_list.article_count"
		FROM "bookmark"
		JOIN "article_list" ON "bookmark".article_list_id = "article_list".id
		WHERE "bookmark".organization_id = $1
		AND user_id = $2
		ORDER BY "bookmark"."def_time" DESC
		LIMIT $3 OFFSET $4`
	var entities []Bookmark

	if err := connection.Select(&entities, query, orgaId, userId, n, offset); err != nil {
		logbuch.Error("Error reading bookmarks by organization id and user id and article list id set with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "offset": offset, "n": n})
		return nil
	}

	for i := range entities {
		entities[i].ArticleList.Name = GetArticleListNameByOrganizationIdAndArticleListIdAndLangId(orgaId, entities[i].ArticleList.ID, langId)
	}

	return entities
}

func DeleteBookmarkById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "bookmark" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting bookmark by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveBookmark(tx *sqlx.Tx, entity *Bookmark) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "bookmark" (organization_id, user_id, article_id, article_list_id)
			VALUES (:organization_id, :user_id, :article_id, :article_list_id) RETURNING id`,
		`UPDATE "bookmark" SET organization_id = :organization_id,
			user_id = :user_id,
			article_id = :article_id,
			article_list_id = :article_list_id
			WHERE id = :id`)
}
