package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

const (
	feedRefBaseQuery = `SELECT "feed_ref".*,
		CASE WHEN "user".id IS NULL THEN 0 ELSE "user".id END "user.id",
		CASE WHEN "user".email IS NULL THEN '' ELSE "user".email END "user.email",
		CASE WHEN "user".firstname IS NULL THEN '' ELSE "user".firstname END "user.firstname",
		CASE WHEN "user".lastname IS NULL THEN '' ELSE "user".lastname END "user.lastname",
		CASE WHEN "user".picture IS NULL THEN '' ELSE "user".picture END "user.picture",
		CASE WHEN "organization_member".username IS NULL THEN '' ELSE "organization_member".username END "user.organization_member.username",
		CASE WHEN "user_group".id IS NULL THEN 0 ELSE "user_group".id END "group.id",
		CASE WHEN "user_group".info IS NULL THEN '' ELSE "user_group".info END "group.info",
		CASE WHEN "user_group".name IS NULL THEN '' ELSE "user_group".name END "group.name",
		CASE WHEN "article".id IS NULL THEN 0 ELSE "article".id END "article.id",
		CASE WHEN "article".wip IS NULL THEN 0 ELSE "article".wip END "article.wip",
		CASE WHEN "article".read_everyone IS NULL THEN FALSE ELSE "article".read_everyone END "article.read_everyone",
		CASE WHEN "article".write_everyone IS NULL THEN FALSE ELSE "article".write_everyone END "article.write_everyone",
		CASE WHEN "article".client_access IS NULL THEN FALSE ELSE "article".client_access END "article.client_access",
		CASE WHEN "article_content".id IS NULL THEN 0 ELSE "article_content".id END "article_content.id",
		CASE WHEN "article_content".title IS NULL THEN '' ELSE "article_content".title END "article_content.title",
		CASE WHEN "article_content"."commit" IS NULL THEN NULL ELSE "article_content"."commit" END "article_content.commit",
		CASE WHEN "article_list".id IS NULL THEN 0 ELSE "article_list".id END "article_list.id",
		CASE WHEN "article_list".public IS NULL THEN FALSE ELSE "article_list".public END "article_list.public"
		FROM "feed_ref"
		LEFT JOIN "user" ON "feed_ref".user_id = "user".id
		LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		LEFT JOIN "user_group" ON "feed_ref".user_group_id = "user_group".id
		LEFT JOIN "article" ON "feed_ref".article_id = "article".id
		LEFT JOIN "article_content" ON "feed_ref".article_content_id = "article_content".id
		LEFT JOIN "article_list" ON "feed_ref".article_list_id = "article_list".id `
)

type FeedRef struct {
	db.BaseEntity

	FeedId hide.ID     `db:"feed_id" json:"feed_id"`
	Key    null.String `json:"key"`
	Value  null.String `json:"value"`

	// referenced objects
	UserID           hide.ID `db:"user_id" json:"user_id"`                       // nullable
	UserGroupID      hide.ID `db:"user_group_id" json:"user_group_id"`           // nullable
	ArticleID        hide.ID `db:"article_id" json:"article_id"`                 // nullable
	ArticleContentID hide.ID `db:"article_content_id" json:"article_content_id"` // nullable
	ArticleListID    hide.ID `db:"article_list_id" json:"article_list_id"`       // nullable

	User           *User           `db:"user" json:"user"`
	UserGroup      *UserGroup      `db:"group" json:"user_group"`
	Article        *Article        `db:"article" json:"article"`
	ArticleContent *ArticleContent `db:"article_content" json:"article_content"`
	ArticleList    *ArticleList    `db:"article_list" json:"article_list"`
}

func FindFeedRefByFeedId(feedId hide.ID) []FeedRef {
	query := `SELECT * FROM "feed_ref" WHERE feed_id = $1`
	var entities []FeedRef

	if err := connection.Select(&entities, query, feedId); err != nil {
		logbuch.Error("Error reading feed references by feed id", logbuch.Fields{"err": err, "feed_id": feedId})
		return nil
	}

	return entities
}

func FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orgaId, langId, feedId hide.ID) []FeedRef {
	query := feedRefBaseQuery + `WHERE "feed_ref".feed_id = $2`
	var entities []FeedRef

	if err := connection.Select(&entities, query, orgaId, feedId); err != nil {
		logbuch.Error("Error reading feed references by organization id and feed id", logbuch.Fields{"err": err, "orga_id": orgaId, "feed_id": feedId})
		return nil
	}

	return joinFeedRefLatesttArticleContentAndArticleListNames(entities, orgaId, langId)
}

func FindFeedRefFeedIdsByArticleId(tx *sqlx.Tx, articleId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT feed_id FROM "feed_ref" WHERE article_id = $1`
	var ids []hide.ID

	if err := tx.Select(&ids, query, articleId); err != nil {
		logbuch.Error("Error reading feed reference feed ids by article id", logbuch.Fields{"err": err, "article_id": articleId})
		return nil
	}

	return ids
}

func FindFeedRefFeedIdsByArticleContentId(tx *sqlx.Tx, contentId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT feed_id FROM "feed_ref" WHERE article_content_id = $1`
	var ids []hide.ID

	if err := tx.Select(&ids, query, contentId); err != nil {
		logbuch.Error("Error reading feed reference feed ids by article content id", logbuch.Fields{"err": err, "content_id": contentId})
		return nil
	}

	return ids
}

func FindFeedRefFeedIdsByArticleListId(tx *sqlx.Tx, listId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT feed_id FROM "feed_ref" WHERE article_list_id = $1`
	var ids []hide.ID

	if err := tx.Select(&ids, query, listId); err != nil {
		logbuch.Error("Error reading feed reference feed ids by article list id", logbuch.Fields{"err": err, "list_id": listId})
		return nil
	}

	return ids
}

func FindFeedRefFeedIdsByUserGroupId(tx *sqlx.Tx, groupId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT feed_id FROM "feed_ref" WHERE user_group_id = $1`
	var ids []hide.ID

	if err := tx.Select(&ids, query, groupId); err != nil {
		logbuch.Error("Error reading feed reference feed ids by user group id", logbuch.Fields{"err": err, "group_id": groupId})
		return nil
	}

	return ids
}

func joinFeedRefLatesttArticleContentAndArticleListNames(refs []FeedRef, orgaId, langId hide.ID) []FeedRef {
	for i := range refs {
		if refs[i].ArticleID != 0 && refs[i].Article != nil {
			refs[i].Article.LatestArticleContent = GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageId(orgaId, refs[i].ArticleID, langId, true)
		}

		if refs[i].ArticleListID != 0 && refs[i].ArticleList != nil {
			refs[i].ArticleList.Name = GetArticleListNameByOrganizationIdAndArticleListIdAndLangId(orgaId, refs[i].ArticleListID, langId)
		}
	}

	return refs
}

func DeleteFeedRefByFeedIds(tx *sqlx.Tx, ids []hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query, args, err := sqlx.In(`DELETE FROM "feed_ref" WHERE feed_id IN (?)`, ids)

	if err != nil {
		logbuch.Error("Error deleting feed ref by feed ids", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	query = connection.Rebind(query)
	_, err = tx.Exec(query, args...)

	if err != nil {
		logbuch.Error("Error deleting feed ref by feed ids", logbuch.Fields{"err": err})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveFeedRef(tx *sqlx.Tx, entity *FeedRef) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "feed_ref" (feed_id,
			"key",
			"value",
			user_id,
			user_group_id,
			article_id,
			article_content_id,
			article_list_id)
			VALUES (:feed_id,
			:key,
			:value,
			:user_id,
			:user_group_id,
			:article_id,
			:article_content_id,
			:article_list_id) RETURNING id`,
		`UPDATE "feed_ref" SET feed_id = :feed_id,
			"key" = :key,
			"value" = :value,
			user_id = :user_id,
			user_group_id = :user_group_id,
			article_id = :article_id,
			article_content_id = :article_content_id,
			article_list_id = :article_list_id
			WHERE id = :id`)
}
