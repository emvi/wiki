package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type ObservedObject struct {
	db.BaseEntity

	UserId        hide.ID `db:"user_id" json:"user_id"`
	ArticleId     hide.ID `db:"article_id" json:"article_id"`           // nullable
	ArticleListId hide.ID `db:"article_list_id" json:"article_list_id"` // nullable
	UserGroupId   hide.ID `db:"user_group_id" json:"user_group_id"`     // nullable
}

func GetObservedObjectByUserIdAndArticleId(userId, articleId hide.ID) *ObservedObject {
	entity := new(ObservedObject)

	if err := connection.Get(entity, `SELECT * FROM "observed_object" WHERE user_id = $1 AND article_id = $2`, userId, articleId); err != nil {
		logbuch.Debug("Observed object by user id and article id not found", logbuch.Fields{"err": err, "user_id": userId, "article_id": articleId})
		return nil
	}

	return entity
}

func GetObservedObjectByUserIdAndArticleListId(userId, articleListId hide.ID) *ObservedObject {
	entity := new(ObservedObject)

	if err := connection.Get(entity, `SELECT * FROM "observed_object" WHERE user_id = $1 AND article_list_id = $2`, userId, articleListId); err != nil {
		logbuch.Debug("Observed object by user id and article list id not found", logbuch.Fields{"err": err, "user_id": userId, "article_list_id": articleListId})
		return nil
	}

	return entity
}

func GetObservedObjectByUserIdAndUserGroupId(userId, userGroupId hide.ID) *ObservedObject {
	entity := new(ObservedObject)

	if err := connection.Get(entity, `SELECT * FROM "observed_object" WHERE user_id = $1 AND user_group_id = $2`, userId, userGroupId); err != nil {
		logbuch.Debug("Observed object by user id and user group id not found", logbuch.Fields{"err": err, "user_id": userId, "user_group_id": userGroupId})
		return nil
	}

	return entity
}

func GetObservedObjectByUserIdAndArticleIdOrArticleListIdOrUserGroupId(userId, articleId, articleListId, userGroupId hide.ID) *ObservedObject {
	entities := new(ObservedObject)

	if err := connection.Get(entities, `SELECT * FROM "observed_object" WHERE user_id = $1 AND (article_id = $2 OR article_list_id = $3 OR user_group_id = $4)`, userId, articleId, articleListId, userGroupId); err != nil {
		logbuch.Debug("Observed object by user id and article id and article list id or user group id not found", logbuch.Fields{"err": err, "user_id": userId, "article_id": articleId, "article_list_id": articleListId, "user_group_id": userGroupId})
		return nil
	}

	return entities
}

func FindObservedObjectUserIdByArticleIdOrArticleListId(articleId, articleListId hide.ID) []hide.ID {
	return FindObservedObjectUserIdByArticleIdOrArticleListIdTx(nil, articleId, articleListId)
}

func FindObservedObjectUserIdByArticleIdOrArticleListIdTx(tx *sqlx.Tx, articleId, articleListId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var ids []hide.ID

	if err := tx.Select(&ids, `SELECT user_id FROM "observed_object" WHERE article_id = $1 OR article_list_id = $2`, articleId, articleListId); err != nil {
		logbuch.Error("Error reading observed object user ids by article id or article list id", logbuch.Fields{"err": err, "article_id": articleId, "article_list_id": articleListId})
		return nil
	}

	return ids
}

func FindObservedObjectUserIdByArticleListIdTx(tx *sqlx.Tx, articleListId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var ids []hide.ID

	if err := tx.Select(&ids, `SELECT user_id FROM "observed_object" WHERE article_list_id = $1`, articleListId); err != nil {
		logbuch.Error("Error reading observed object user ids by article list id", logbuch.Fields{"err": err, "article_list_id": articleListId})
		return nil
	}

	return ids
}

func FindObservedObjectUserIdByUserGroupIdTx(tx *sqlx.Tx, groupId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var ids []hide.ID

	if err := tx.Select(&ids, `SELECT user_id FROM "observed_object" WHERE user_group_id = $1`, groupId); err != nil {
		logbuch.Error("Error reading observed object user ids by user group id", logbuch.Fields{"err": err, "group_id": groupId})
		return nil
	}

	return ids
}

func DeleteObservedObjectById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "observed_object" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting observed object by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveObservedObject(tx *sqlx.Tx, entity *ObservedObject) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "observed_object" (user_id,
			article_id,
			article_list_id,
			user_group_id)
			VALUES (:user_id,
			:article_id,
			:article_list_id,
			:user_group_id) RETURNING id`,
		`UPDATE "observed" SET user_id = :user_id,
			article_id = :article_id,
			article_list_id = :article_list_id,
			user_group_id = :user_group_id
			WHERE id = :id`)
}
