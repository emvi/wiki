package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	memberCount = `(SELECT COUNT(DISTINCT "user_group_member".id) FROM "user_group_member"
		JOIN "organization_member" ON "user_group_member".user_id = "organization_member".user_id
		WHERE user_group_id = "user_group".id
		AND "organization_member".active IS TRUE) AS "member_count"`
)

type UserGroup struct {
	db.BaseEntity

	OrganizationId hide.ID     `db:"organization_id" json:"organization_id"`
	Name           string      `json:"name"`
	Info           null.String `json:"info"`
	Immutable      bool        `json:"immutable"`

	MemberCount uint `db:"member_count" json:"member"`
}

func GetUserGroupByOrganizationIdAndId(orgaId, id hide.ID) *UserGroup {
	return GetUserGroupByOrganizationIdAndIdTx(nil, orgaId, id)
}

func GetUserGroupByOrganizationIdAndIdTx(tx *sqlx.Tx, orgaId, id hide.ID) *UserGroup {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "user_group".*,
		` + memberCount + `
		FROM "user_group"
		WHERE organization_id = $1 AND id = $2`
	entity := new(UserGroup)

	if err := tx.Get(entity, query, orgaId, id); err != nil {
		logbuch.Debug("User group by organization id and id not found", logbuch.Fields{"err": err, "organization_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetUserGroupByOrganizationIdAndName(orgaId hide.ID, name string) *UserGroup {
	return GetUserGroupByOrganizationIdAndNameTx(nil, orgaId, name)
}

func GetUserGroupByOrganizationIdAndNameTx(tx *sqlx.Tx, orgaId hide.ID, name string) *UserGroup {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "user_group".*,
		` + memberCount + `
		FROM "user_group"
		WHERE organization_id = $1 AND LOWER(name) = LOWER($2)`
	entity := new(UserGroup)

	if err := tx.Get(entity, query, orgaId, name); err != nil {
		logbuch.Debug("User group by organization id and name not found", logbuch.Fields{"err": err, "organization_id": orgaId, "name": name})
		return nil
	}

	return entity
}

func FindUserGroupByOrganizationIdAndUserIdAndObservedWithLimit(orgaId, userId hide.ID, offset, n int) []UserGroup {
	query := `SELECT DISTINCT ON ("user_group".id) "user_group".*,
		` + memberCount + `
		FROM "user_group"
		JOIN "observed_object" ON "user_group".id = "observed_object".user_group_id AND "observed_object".user_id = $2 
		WHERE organization_id = $1
		LIMIT $4 OFFSET $3`
	var entities []UserGroup

	if err := connection.Select(&entities, query, orgaId, userId, offset, n); err != nil {
		logbuch.Error("User groups by organization id and user id and observed not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "offset": offset, "n": n})
		return nil
	}

	return entities
}

// FindUserGroupsByOrganizationIdAndAccessByUserIdOnlyTx returns groups where only the given user has access to.
func FindUserGroupsByOrganizationIdAndAccessByUserIdOnlyTx(tx *sqlx.Tx, orgaId, userId hide.ID) []UserGroup {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "user_group".*
		FROM "user_group"
		JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE organization_id = $1
		AND "user_group_member".user_id = $2
		AND "user_group_member".is_moderator IS TRUE
		AND 1 = (SELECT COUNT(1) FROM "user_group_member"
		WHERE user_group_id = "user_group".id
		AND "user_group_member".is_moderator IS TRUE)`
	var entities []UserGroup

	if err := tx.Select(&entities, query, orgaId, userId); err != nil {
		logbuch.Error("User groups by organization id and access user id only not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entities
}

func FindUserGroupsByOrganizationIdAndNameOrInfo(orgaId hide.ID, keywords string, filter *SearchUserGroupFilter) []UserGroup {
	query, params := buildUserGroupsByOrganizationIdAndNameOrInfoQuery(orgaId, keywords, filter, false)
	var entities []UserGroup

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading user groups by organization id and query string", logbuch.Fields{"err": err, "organization_id": orgaId, "keywords": keywords})
		return nil
	}

	return entities
}

func CountUserGroupsByOrganizationIdAndNameOrInfo(orgaId hide.ID, keywords string, filter *SearchUserGroupFilter) int {
	query, params := buildUserGroupsByOrganizationIdAndNameOrInfoQuery(orgaId, keywords, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting user groups by organization id and query string", logbuch.Fields{"err": err, "organization_id": orgaId, "keywords": keywords})
		return 0
	}

	return count
}

func CountUserGroupByOrganizationId(orgaId hide.ID) int {
	query := `SELECT COUNT(1) FROM "user_group" WHERE organization_id = $1`
	var count int

	if err := connection.Get(&count, query, orgaId); err != nil {
		logbuch.Error("Error counting user groups by organization id", logbuch.Fields{"err": err, "organization_id": orgaId})
		return 0
	}

	return count
}

func buildUserGroupsByOrganizationIdAndNameOrInfoQuery(orgaId hide.ID, keywords string, filter *SearchUserGroupFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 1)
	params[0] = orgaId
	var sb strings.Builder

	if count {
		sb.WriteString(`SELECT COUNT(1) `)
	} else {
		sb.WriteString(`SELECT "user_group".*, ` + memberCount)
	}

	sb.WriteString(`FROM "user_group" `)

	if len(filter.UserIds) > 0 {
		// distinct is not required when we filter for user ids anyway
		sb.WriteString(`JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id `)
	}

	sb.WriteString(`WHERE organization_id = $1 `)

	if keywords != "" {
		params = append(params, keywords)
		sb.WriteString(`AND (SIMILARITY("name", $2) > 0.2
			OR LOWER("name") LIKE LOWER('%'||$2||'%')
			OR SIMILARITY("info", $2) > 0.2
			OR LOWER("info") LIKE LOWER('%'||$2||'%')) `)
	}

	// add field filter (joined with "AND")
	index := len(params) + 1
	fieldFilter, index, params := filter.addFieldFilter("", index, params, []string{filter.Name, filter.Info}, "name", "info")
	sb.WriteString(fieldFilter)

	if len(filter.UserIds) > 0 {
		var userFilter string
		userFilter, index, params = filter.filterInIds(filter.UserIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "user_group_member".user_id %v`, userFilter))
	}

	// add date filter
	dateFilter, index, params := filter.addDateFilter("user_group", index, params)
	sb.WriteString(dateFilter)

	if !count {
		// sorting
		defaultFields := []SortValue{{"name", sortDirectionASC}}
		sb.WriteString(filter.addSorting("user_group", defaultFields, SortValue{"name", filter.SortName}, SortValue{"info", filter.SortInfo}))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	}

	return sb.String(), params
}

func DeleteUserGroupById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "observed_object" WHERE user_group_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting observed object by user group id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "user_group_member" WHERE user_group_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting user group members", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_access" WHERE user_group_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article access", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "user_group" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting user group", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveUserGroup(tx *sqlx.Tx, entity *UserGroup) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "user_group" (name,
			info,
			immutable,
			organization_id)
			VALUES (:name,
			:info,
			:immutable,
			:organization_id)
			RETURNING id`,
		`UPDATE "user_group" SET name = :name,
			info = :info,
			immutable = :immutable,
			organization_id = :organization_id
			WHERE id = :id`)
}
