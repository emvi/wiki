package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
)

type UserGroupMember struct {
	db.BaseEntity

	UserGroupId hide.ID `db:"user_group_id" json:"user_group_id"`
	UserId      hide.ID `db:"user_id" json:"user_id"`
	IsModerator bool    `db:"is_moderator" json:"is_moderator"`

	User *User `db:"user" json:"user"`
}

func GetUserGroupMemberById(id hide.ID) *UserGroupMember {
	return GetUserGroupMemberByIdTx(nil, id)
}

func GetUserGroupMemberByIdTx(tx *sqlx.Tx, id hide.ID) *UserGroupMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(UserGroupMember)

	if err := tx.Get(entity, `SELECT * FROM "user_group_member" WHERE id = $1`, id); err != nil {
		logbuch.Debug("User group member by id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func GetUserGroupMemberByGroupIdAndIdTx(tx *sqlx.Tx, groupId, id hide.ID) *UserGroupMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(UserGroupMember)

	if err := tx.Get(entity, `SELECT * FROM "user_group_member" WHERE user_group_id = $1 AND id = $2`, groupId, id); err != nil {
		logbuch.Debug("User group member by group id and id not found", logbuch.Fields{"err": err, "group_id": groupId, "id": id})
		return nil
	}

	return entity
}

func GetUserGroupMemberByGroupIdAndUserId(groupId, userId hide.ID) *UserGroupMember {
	return GetUserGroupMemberByGroupIdAndUserIdTx(nil, groupId, userId)
}

func GetUserGroupMemberByGroupIdAndUserIdTx(tx *sqlx.Tx, groupId, userId hide.ID) *UserGroupMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(UserGroupMember)

	if err := tx.Get(entity, `SELECT * FROM "user_group_member" WHERE user_group_id = $1 AND user_id = $2`, groupId, userId); err != nil {
		logbuch.Debug("User group member by group id and user id not found", logbuch.Fields{"err": err, "group_id": groupId, "user_id": userId})
		return nil
	}

	return entity
}

func CountUserGroupMemberByUserGroupId(groupId hide.ID) int {
	var count int

	if err := connection.Get(&count, `SELECT COUNT(1) FROM "user_group_member" WHERE user_group_id = $1`, groupId); err != nil {
		logbuch.Error("Error couting user group members by user group id", logbuch.Fields{"err": err, "group_id": groupId})
		return 0
	}

	return count
}

func FindUserGroupMemberOnlyByUserGroupIdTx(tx *sqlx.Tx, groupId hide.ID) []UserGroupMember {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "user_group_member".* FROM "user_group_member"
		JOIN "organization_member" ON "user_group_member".user_id = "organization_member".user_id
		WHERE user_group_id = $1
		AND "organization_member".active IS TRUE
		ORDER BY "user_group_member".is_moderator ASC`
	var entities []UserGroupMember

	if err := tx.Select(&entities, query, groupId); err != nil {
		logbuch.Error("Error reading only user group member by group id", logbuch.Fields{"err": err, "group_id": groupId})
		return nil
	}

	return entities
}

func FindUserGroupMemberUserIdByUserGroupId(groupId hide.ID) []hide.ID {
	return FindUserGroupMemberUserIdByUserGroupIdTx(nil, groupId)
}

func FindUserGroupMemberUserIdByUserGroupIdTx(tx *sqlx.Tx, groupId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "user_group_member".user_id FROM "user_group_member"
		JOIN "organization_member" ON "user_group_member".user_id = "organization_member".user_id
		WHERE user_group_id = $1
		AND "organization_member".active IS TRUE
		ORDER BY "user_group_member".is_moderator ASC`
	var ids []hide.ID

	if err := tx.Select(&ids, query, groupId); err != nil {
		logbuch.Error("Error reading user group member user id by group id", logbuch.Fields{"err": err, "group_id": groupId})
		return nil
	}

	return ids
}

func FindUserGroupMemberUserByOrganizationIdAndUserGroupId(orgaId, groupId hide.ID) []User {
	return FindUserGroupMemberUserByOrganizationIdAndUserGroupIdTx(nil, orgaId, groupId)
}

func FindUserGroupMemberUserByOrganizationIdAndUserGroupIdTx(tx *sqlx.Tx, orgaId, groupId hide.ID) []User {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := userBaseQueryHead + ` FROM "user_group_member"
		JOIN "user_group" ON "user_group_member".user_group_id = "user_group".id
		JOIN "user" ON "user_group_member".user_id = "user".id
		JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		WHERE "user_group".organization_id = $1
		AND user_group_id = $2`
	var entities []User

	if err := tx.Select(&entities, query, orgaId, groupId); err != nil {
		logbuch.Error("Error reading user group member by organization id and user group id", logbuch.Fields{"err": err, "orga_id": orgaId, "group_id": groupId})
		return nil
	}

	return entities
}

func FindUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimit(orgaId, groupId hide.ID, filter *SearchUserGroupMemberFilter) []UserGroupMember {
	query, params := buildUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimitQuery(orgaId, groupId, filter, false)
	var entities []UserGroupMember

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading user group member by organization id and group id and filter with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "group_id": groupId, "filter": filter})
		return nil
	}

	return entities
}

func CountUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimit(orgaId, groupId hide.ID, filter *SearchUserGroupMemberFilter) int {
	query, params := buildUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimitQuery(orgaId, groupId, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting user group member by organization id and group id and filter with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "group_id": groupId, "filter": filter})
		return 0
	}

	return count
}

func buildUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimitQuery(orgaId, groupId hide.ID, filter *SearchUserGroupMemberFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 2)
	params[0] = orgaId
	params[1] = groupId

	defaultFields := []SortValue{{`"organization_member".username`, sortDirectionASC}}
	customFields := []SortValue{{`"organization_member".username`, filter.SortUsername}, {`"user".email`, filter.SortEmail}, SortValue{`"user".firstname`, filter.SortFirstname}, SortValue{`"user".lastname`, filter.SortLastname}}
	selectOrCount := `COUNT(DISTINCT("user_group_member".id))`

	if !count {
		sortColumns := filter.getSortColumns("article", defaultFields, customFields...)
		selectOrCount = `DISTINCT ON ("user_group_member".id, ` + sortColumns + `)
			"user_group_member".*,
			"user".id "user.id",
			"user".email "user.email",
			"user".firstname "user.firstname",
			"user".lastname "user.lastname",
			"user".info "user.info",
			"user".picture "user.picture",
			"organization_member".username "user.organization_member.username",
			"organization_member".info "user.organization_member.info"`
	}

	var sb strings.Builder
	sb.WriteString(`SELECT `)
	sb.WriteString(selectOrCount)
	sb.WriteString(` FROM "user_group_member"
		JOIN "user" ON "user_group_member".user_id = "user".id
		JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		WHERE user_group_id = $2
		AND "organization_member".active IS TRUE `)

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
		sb.WriteString(filter.addSorting("user_group_member", defaultFields, customFields...))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	}

	return sb.String(), params
}

func DeleteUserGroupMemberByUserGroupIdAndId(tx *sqlx.Tx, groupId, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "user_group_member" WHERE user_group_id = $1 AND id = $2`, groupId, id); err != nil {
		logbuch.Error("Error deleting user group member by user group id and id", logbuch.Fields{"err": err, "group_id": groupId, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func DeleteUserGroupMemberByUserGroupId(tx *sqlx.Tx, groupId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := tx.Exec(`DELETE FROM "user_group_member" WHERE user_group_id = $1`, groupId); err != nil {
		logbuch.Error("Error deleting user group member by user group id", logbuch.Fields{"err": err, "group_id": groupId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveUserGroupMember(tx *sqlx.Tx, entity *UserGroupMember) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "user_group_member" (is_moderator,
			user_group_id,
			user_id)
			VALUES (:is_moderator,
			:user_group_id,
			:user_id)
			RETURNING id`,
		`UPDATE "user_group_member" SET is_moderator = :is_moderator,
			user_group_id = :user_group_id,
			user_id = :user_id
			WHERE id = :id`)
}
