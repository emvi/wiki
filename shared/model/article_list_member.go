package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ArticleListMember struct {
	db.BaseEntity

	ArticleListId hide.ID `db:"article_list_id" json:"article_list_id"`
	UserId        hide.ID `db:"user_id" json:"user_id"`             // nullable
	UserGroupId   hide.ID `db:"user_group_id" json:"user_group_id"` // nullable
	IsModerator   bool    `db:"is_moderator" json:"is_moderator"`

	User      *User      `db:"user" json:"user"`
	UserGroup *UserGroup `db:"user_group" json:"user_group"`
}

func GetArticleListMemberByArticleListIdAndId(listId, id hide.ID) *ArticleListMember {
	return GetArticleListMemberByArticleListIdAndIdTx(nil, listId, id)
}

func GetArticleListMemberByArticleListIdAndIdTx(tx *sqlx.Tx, listId, id hide.ID) *ArticleListMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(ArticleListMember)

	if err := tx.Get(entity, `SELECT * FROM "article_list_member"
		WHERE article_list_id = $1
		AND id = $2`, listId, id); err != nil {
		logbuch.Debug("Article list member by article list id and id not found", logbuch.Fields{"err": err, "list_id": listId, "id": id})
		return nil
	}

	return entity
}

func GetArticleListMemberByArticleListIdAndUserIdTx(tx *sqlx.Tx, listId, userId hide.ID) *ArticleListMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(ArticleListMember)

	if err := tx.Get(entity, `SELECT * FROM "article_list_member"
		WHERE article_list_id = $1
		AND user_id = $2`, listId, userId); err != nil {
		logbuch.Debug("Article list member by article list id and user id not found", logbuch.Fields{"err": err, "list_id": listId, "user_id": userId})
		return nil
	}

	return entity
}

func GetArticleListMemberByArticleListIdAndUserGroupIdTx(tx *sqlx.Tx, listId, groupId hide.ID) *ArticleListMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(ArticleListMember)

	if err := tx.Get(entity, `SELECT * FROM "article_list_member"
		WHERE article_list_id = $1
		AND user_group_id = $2`, listId, groupId); err != nil {
		logbuch.Debug("Article list member by article list id and user group id not found", logbuch.Fields{"err": err, "list_id": listId, "group_id": groupId})
		return nil
	}

	return entity
}

func FindArticleListMemberModeratorByArticleListIdAndUserId(listId, userId hide.ID) []ArticleListMember {
	query := `SELECT "article_list_member".* FROM "article_list_member"
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		LEFT JOIN "user" ON "article_list_member".user_id = "user".id OR "user_group_member".user_id = "user".id
		LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id
		WHERE article_list_id = $1
		AND "article_list_member".is_moderator IS TRUE
		AND "organization_member".active IS TRUE
		AND ("article_list_member".user_id = $2 OR "user_group_member".user_id = $2)`
	var entities []ArticleListMember

	if err := connection.Select(&entities, query, listId, userId); err != nil {
		logbuch.Error("Article list member by article list id and user id not found", logbuch.Fields{"err": err, "list_id": listId, "user_id": userId})
		return nil
	}

	return entities
}

func FindArticleListMemberModeratorByArticleListIdTx(tx *sqlx.Tx, listId hide.ID) []ArticleListMember {
	query := `SELECT DISTINCT ON ("article_list_member".id) "article_list_member".* FROM "article_list_member"
		LEFT JOIN "user" ON "article_list_member".user_id = "user".id
		LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE article_list_id = $1
		AND "article_list_member".is_moderator IS TRUE
		AND ("organization_member".id IS NULL OR "organization_member".active IS TRUE)`
	var entities []ArticleListMember

	if err := tx.Select(&entities, query, listId); err != nil {
		logbuch.Error("Article list member by article list id not found", logbuch.Fields{"err": err, "list_id": listId})
		return nil
	}

	return entities
}

func FindArticleListMemberByArticleListId(listId hide.ID) []ArticleListMember {
	var entities []ArticleListMember

	if err := connection.Select(&entities, `SELECT * FROM "article_list_member" WHERE article_list_id = $1`, listId); err != nil {
		logbuch.Error("Error reading article list member by article list id", logbuch.Fields{"err": err, "list_id": listId})
		return nil
	}

	return entities
}

func FindArticleListMemberByArticleListIdAndUserIdIncludingUserGroupMember(tx *sqlx.Tx, listId, userId hide.ID) []ArticleListMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "article_list_member".* FROM "article_list_member"
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE article_list_id = $1
		AND ("article_list_member".user_id = $2 OR "user_group_member".user_id = $2)`
	var entities []ArticleListMember

	if err := tx.Select(&entities, query, listId, userId); err != nil {
		logbuch.Error("Article list member by article list id and user id not found", logbuch.Fields{"err": err, "list_id": listId, "user_id": userId})
		return nil
	}

	return entities
}

func FindArticleListMemberUserIdByArticleListId(listId hide.ID) []hide.ID {
	return FindArticleListMemberUserIdByArticleListIdTx(nil, listId)
}

func FindArticleListMemberUserIdByArticleListIdTx(tx *sqlx.Tx, listId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "article_list_member".user_id
		FROM "article_list_member"
		WHERE article_list_id = $1
		AND "article_list_member".user_id IS NOT NULL
		UNION
		SELECT "user_group_member".user_id
		FROM "article_list_member"
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE article_list_id = $1
		AND "article_list_member".user_id IS NULL`
	var ids []hide.ID

	if err := tx.Select(&ids, query, listId); err != nil {
		logbuch.Error("Error reading article list member ids by article list id", logbuch.Fields{"err": err, "list_id": listId})
		return nil
	}

	return ids
}

func FindArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimit(orgaId, listId hide.ID, filter *SearchArticleListMemberFilter) []ArticleListMember {
	query, params := buildArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimitQuery(orgaId, listId, filter, false)
	var entities []ArticleListMember

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading article list member by organization id and list id", logbuch.Fields{"err": err, "orga_id": orgaId, "list_id": listId, "filter": filter})
		return nil
	}

	return entities
}

func CountArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimit(orgaId, listId hide.ID, filter *SearchArticleListMemberFilter) int {
	query, params := buildArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimitQuery(orgaId, listId, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting article list member by organization id and list id", logbuch.Fields{"err": err, "orga_id": orgaId, "list_id": listId, "filter": filter})
		return 0
	}

	return count
}

func buildArticleListMemberByOrganizationIdAndArticleListIdAndFilterLimitQuery(orgaId, listId hide.ID, filter *SearchArticleListMemberFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 2)
	params[0] = orgaId
	params[1] = listId
	selectOrCount := `COUNT(DISTINCT("article_list_member".id))`

	if !count {
		selectOrCount = `"article_list_member".*,
			CASE WHEN "user".id IS NULL THEN 0 ELSE "user".id END "user.id",
			CASE WHEN "user".email IS NULL THEN '' ELSE "user".email END "user.email",
			CASE WHEN "user".firstname IS NULL THEN '' ELSE "user".firstname END "user.firstname",
			CASE WHEN "user".lastname IS NULL THEN '' ELSE "user".lastname END "user.lastname",
			"user".language "user.language",
			"user".info "user.info",
			"user".picture "user.picture",
			CASE WHEN "organization_member".id IS NULL THEN 0 ELSE "organization_member".id END "user.organization_member.id",
			CASE WHEN "organization_member".username IS NULL THEN '' ELSE "organization_member".username END "user.organization_member.username",
			CASE WHEN "organization_member".is_moderator IS NULL THEN FALSE ELSE "organization_member".is_moderator END "user.organization_member.is_moderator",
			CASE WHEN "organization_member".is_admin IS NULL THEN FALSE ELSE "organization_member".is_admin END "user.organization_member.is_admin",
			CASE WHEN "user_group".id IS NULL THEN 0 ELSE "user_group".id END "user_group.id",
			CASE WHEN "user_group".name IS NULL THEN '' ELSE "user_group".name END "user_group.name",
			CASE WHEN "user_group".info IS NULL THEN '' ELSE "user_group".info END "user_group.info",
			CASE WHEN "user_group".def_time IS NULL THEN TO_TIMESTAMP(0) ELSE "user_group".def_time END "user_group.def_time",
			CASE WHEN "user_group".mod_time IS NULL THEN TO_TIMESTAMP(0) ELSE "user_group".mod_time END "user_group.mod_time", ` +
			strings.Replace(memberCount, `"member_count"`, `"user_group.member_count"`, 1)
	}

	var sb strings.Builder
	sb.WriteString(`SELECT `)
	sb.WriteString(selectOrCount)
	sb.WriteString(` FROM "article_list_member"
		LEFT JOIN "user" ON "article_list_member".user_id = "user".id
		LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		WHERE article_list_id = $2
		AND ("organization_member" IS NULL OR "organization_member".active IS TRUE) `)

	// add field filter (joined with "AND")
	index := 3

	if len(filter.UserIds) > 0 {
		var userFilter string
		userFilter, index, params = filter.filterInIds(filter.UserIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "user".id %v`, userFilter))
	}

	// add date filter
	dateFilter, _, params := filter.addDateFilter("user", index, params)
	sb.WriteString(dateFilter)

	if !count {
		// sorting
		defaultFields := []SortValue{{`"organization_member".username`, sortDirectionASC}, {`"user_group".name`, sortDirectionASC}}
		sb.WriteString(filter.addSorting("article_list_member", defaultFields, SortValue{`"organization_member".username`, filter.SortUsername}, SortValue{`"user".email`, filter.SortEmail}, SortValue{`"user".firstname`, filter.SortFirstname}, SortValue{`"user".lastname`, filter.SortLastname}))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	}

	return sb.String(), params
}

func DeleteArticleListMemberByArticleListIdAndId(tx *sqlx.Tx, listId, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "article_list_member" WHERE article_list_id = $1 AND id = $2`, listId, id); err != nil {
		logbuch.Error("Error deleting article list member by article list id and id", logbuch.Fields{"err": err, "list_id": listId, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func DeleteArticleListMemberByArticleListId(tx *sqlx.Tx, listId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "article_list_member" WHERE article_list_id = $1`, listId); err != nil {
		logbuch.Error("Error deleting article list member by article list id", logbuch.Fields{"err": err, "list_id": listId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleListMember(tx *sqlx.Tx, entity *ArticleListMember) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_list_member" (article_list_id,
			user_id,
			user_group_id,
			is_moderator)
			VALUES (:article_list_id,
			:user_id,
			:user_group_id,
			:is_moderator) RETURNING id`,
		`UPDATE "article_list_member" SET article_list_id = :article_list_id,
			user_id = :user_id,
			user_group_id = :user_group_id,
			is_moderator = :is_moderator
			WHERE id = :id`)
}
