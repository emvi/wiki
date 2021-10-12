package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type ArticleAccess struct {
	db.BaseEntity

	UserId      hide.ID `db:"user_id" json:"user_id"`             // nullable
	UserGroupId hide.ID `db:"user_group_id" json:"user_group_id"` // nullable
	ArticleId   hide.ID `db:"article_id" json:"article_id"`
	Write       bool    `json:"write"`

	User  *User      `db:"user" json:"user"`
	Group *UserGroup `db:"group" json:"group"`
}

func FindArticleAccessByArticleIdAndUserId(articleId, userId hide.ID) []ArticleAccess {
	return FindArticleAccessByArticleIdAndUserIdTx(nil, articleId, userId)
}

func FindArticleAccessByArticleIdAndUserIdTx(tx *sqlx.Tx, articleId, userId hide.ID) []ArticleAccess {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var entities []ArticleAccess

	if err := tx.Select(&entities, `SELECT "article_access".*
		FROM "article_access"
		LEFT JOIN "user" ON "article_access".user_id = "user".id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE article_id = $1
		AND ("user_group_member".user_id = $2 OR "user".id = $2)`, articleId, userId); err != nil {
		logbuch.Debug("Article access by article id and user id or user group id not found", logbuch.Fields{"err": err, "article_id": articleId, "user_or_user_group_id": userId})
		return nil
	}

	return entities
}

func FindArticleAccessByArticleIdAndUserIdAndWrite(articleId, userId hide.ID, writePermission bool) []ArticleAccess {
	query := `SELECT "article_access".*
		FROM "article_access"
		LEFT JOIN "user" ON "article_access".user_id = "user".id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE article_id = $1 `

	if writePermission {
		query += `AND "article_access".write IS TRUE `
	}

	query += `AND ("user_group_member".user_id = $2 OR "user".id = $2)`
	var entities []ArticleAccess

	if err := connection.Select(&entities, query, articleId, userId); err != nil {
		logbuch.Debug("Article access by article id and user id or user group id and write permission not found", logbuch.Fields{"err": err, "article_id": articleId, "user_or_user_group_id": userId, "write": writePermission})
		return nil
	}

	return entities
}

func FindArticleAccessByOrganizationIdAndArticleId(orgaId, articleId hide.ID) []ArticleAccess {
	return FindArticleAccessByOrganizationIdAndArticleIdTx(nil, orgaId, articleId)
}

func FindArticleAccessByOrganizationIdAndArticleIdTx(tx *sqlx.Tx, orgaId, articleId hide.ID) []ArticleAccess {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	// also used to copy an article, so no limit!
	query := `SELECT "article_access".*,
		CASE WHEN "user".id IS NULL THEN 0 ELSE "user".id END "user.id",
		CASE WHEN "user".email IS NULL THEN '' ELSE "user".email END "user.email",
		CASE WHEN "user".firstname IS NULL THEN '' ELSE "user".firstname END "user.firstname",
		CASE WHEN "user".lastname IS NULL THEN '' ELSE "user".lastname END "user.lastname",
		CASE WHEN "user".language IS NULL THEN '' ELSE "user".language END "user.language",
		CASE WHEN "user".info IS NULL THEN '' ELSE "user".info END "user.info",
		CASE WHEN "user".picture IS NULL THEN NULL ELSE "user".picture END "user.picture",
		CASE WHEN "user_group".organization_id IS NULL THEN 0 ELSE "user_group".organization_id END "group.organization_id",
		CASE WHEN "user_group".id IS NULL THEN 0 ELSE "user_group".id END "group.id",
		CASE WHEN "user_group".info IS NULL THEN '' ELSE "user_group".info END "group.info",
		CASE WHEN "user_group".name IS NULL THEN '' ELSE "user_group".name END "group.name",
		CASE WHEN "organization_member".username IS NULL THEN '' ELSE "organization_member".username END "user.organization_member.username"
		FROM "article_access"
		LEFT JOIN "user" ON "article_access".user_id = "user".id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id
		LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		WHERE article_id = $2
		AND ("organization_member" IS NULL OR "organization_member".active IS TRUE)`
	var entities []ArticleAccess

	if err := tx.Select(&entities, query, orgaId, articleId); err != nil {
		logbuch.Error("Article access by organization id and article id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId})
		return nil
	}

	for i := range entities {
		if entities[i].User.ID == 0 {
			entities[i].User = nil
		} else {
			entities[i].Group = nil
		}
	}

	return entities
}

func DeleteArticleAccessByArticleId(tx *sqlx.Tx, articleId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "article_access" WHERE article_id = $1`, articleId)

	if err != nil {
		logbuch.Error("Error deleting article access by article id", logbuch.Fields{"err": err, "article_id": articleId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleAccess(tx *sqlx.Tx, entity *ArticleAccess) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_access" (user_id,
			user_group_id,
			article_id,
			write)
			VALUES (:user_id,
			:user_group_id,
			:article_id,
			:write)
			RETURNING id`,
		`UPDATE "article_access" SET user_id = :user_id,
			user_group_id = :user_group_id,
			article_id = :article_id,
			write = :write
			WHERE id = :id`)
}
