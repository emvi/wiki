package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	articleListBaseForUserQuery = `FROM "article_list"
		JOIN "article_list_member" ON "article_list".id = "article_list_member".article_list_id
		JOIN "article_list_name" ON "article_list".id = "article_list_name".article_list_id
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE "article_list".organization_id = $1 `
	articleListBaseForClientQuery = `FROM "article_list"
		JOIN "article_list_name" ON "article_list".id = "article_list_name".article_list_id
		WHERE "article_list".organization_id = $1 `
	articleListEntryCountForUserQuery = `(SELECT COUNT(1) FROM "article_list_entry"
		JOIN "article" ON "article_list_entry".article_id = "article".id
		WHERE article_list_id = "article_list".id
		AND ("article".write_everyone IS TRUE
			OR "article".read_everyone IS TRUE
			OR EXISTS (SELECT 1 FROM "article_access"
			LEFT JOIN "user_group_member" ON "article_access".user_group_id = "user_group_member".user_group_id
			WHERE article_id = "article".id
			AND ("article_access".user_id = $2 OR "user_group_member".user_id = $2)))
		) AS "article_count"`
	articleListEntryCountForClientQuery = `(SELECT COUNT(1) FROM "article_list_entry"
		JOIN "article" ON "article_list_entry".article_id = "article".id
		WHERE article_list_id = "article_list".id
		AND "article".client_access IS TRUE
		) AS "article_count"`
	articleListSelectNameQuery = `AND (
		(
			(SELECT EXISTS (SELECT 1 FROM article_list_name AS list WHERE list.article_list_id = article_list.id AND list.language_id = $%d)) IS TRUE
			AND language_id = $%d
		)
		OR (
			(SELECT EXISTS (SELECT 1 FROM article_list_name AS list WHERE list.article_list_id = article_list.id AND list.language_id = $%d)) IS FALSE
			AND (SELECT EXISTS (SELECT list.language_id FROM article_list_name AS list WHERE list.article_list_id = article_list.id AND list.language_id = (SELECT id FROM "language" WHERE "language".organization_id = $1 AND "language".default IS TRUE))) IS TRUE
			AND language_id = (SELECT id FROM "language" WHERE "language".organization_id = $1 AND "language".default IS TRUE)
		)
	) `
)

type ArticleList struct {
	db.BaseEntity

	OrganizationId hide.ID `db:"organization_id" json:"organization_id"`
	Public         bool    `json:"public"`
	Pinned         bool    `json:"pinned"`
	ClientAccess   bool    `db:"client_access" json:"client_access"`

	ArticleCount uint              `db:"article_count" json:"articles"`
	Name         *ArticleListName  `db:"name" json:"name"`
	Names        []ArticleListName `db:"-" json:"names"`
}

func GetArticleListByOrganizationIdAndId(orgaId, id hide.ID) *ArticleList {
	return GetArticleListByOrganizationIdAndIdTx(nil, orgaId, id)
}

func GetArticleListByOrganizationIdAndIdTx(tx *sqlx.Tx, orgaId, id hide.ID) *ArticleList {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "article_list".*
		FROM "article_list"
		WHERE organization_id = $1 AND id = $2`
	entity := new(ArticleList)

	if err := tx.Get(entity, query, orgaId, id); err != nil {
		logbuch.Debug("Article list by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetArticleListByOrganizationIdAndUserIdAndId(orgaId, userId, id hide.ID) *ArticleList {
	query := ""
	params := make([]interface{}, 2)
	params[0] = orgaId

	if userId != 0 {
		params[1] = userId
		params = append(params, id)
		query = `SELECT "article_list".*,` + articleListEntryCountForUserQuery + `
			FROM "article_list"
			WHERE organization_id = $1 AND id = $3`
	} else {
		params[1] = id
		query = `SELECT "article_list".*,` + articleListEntryCountForClientQuery + `
			FROM "article_list"
			WHERE organization_id = $1 AND id = $2
			AND "article_list".client_access IS TRUE`
	}

	entity := new(ArticleList)

	if err := connection.Get(entity, query, params...); err != nil {
		logbuch.Debug("Article list by organization id and user id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "id": id})
		return nil
	}

	return entity
}

func GetArticleListIdByOrganizationIdAndUserIdAndIdAndPinned(orgaId, id hide.ID) *ArticleList {
	query := `SELECT "article_list".id
		FROM "article_list"
		WHERE organization_id = $1 AND id = $2 AND pinned IS TRUE`
	entity := new(ArticleList)

	if err := connection.Get(entity, query, orgaId, id); err != nil {
		logbuch.Debug("Article list by organization id and user id and id and pinned not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func FindArticleListByOrganizationIdAndUserIdAndLanguageIdAndPrivateWithLimit(orgaId, userId, langId hide.ID, offset, n int) []ArticleList {
	query := `SELECT DISTINCT ON ("article_list".id) "article_list".*,
		"article_list_name".name "name.name", "article_list_name".info "name.info",
		` + articleListEntryCountForUserQuery + `
		FROM "article_list"
		JOIN "article_list_name" ON "article_list".id = "article_list_name".article_list_id
		JOIN "article_list_member" ON "article_list".id = "article_list_member".article_list_id
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE "article_list".organization_id = $1
		AND "article_list".public IS FALSE
		AND ("article_list_member".user_id = $2 OR "user_group_member".user_id = $2)
		AND (SELECT COUNT(1) FROM "article_list_member" m WHERE m.article_list_id = "article_list".id) = 1
		` + fmt.Sprintf(articleListSelectNameQuery, 5, 5, 5) + `
		LIMIT $4 OFFSET $3`
	var entities []ArticleList

	if err := connection.Select(&entities, query, orgaId, userId, offset, n, langId); err != nil {
		logbuch.Error("Article list by organization id and user id and language id and private not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "offset": offset, "n": n})
		return nil
	}

	return entities
}

func FindArticleListByOrganizationIdAndUserIdAndLanguageIdAndObservedWithLimit(orgaId, userId, langId hide.ID, offset, n int) []ArticleList {
	query := `SELECT DISTINCT ON ("article_list".id) "article_list".*,
		"article_list_name".name "name.name", "article_list_name".info "name.info",
		` + articleListEntryCountForUserQuery + `
		FROM "article_list"
		JOIN "article_list_name" ON "article_list".id = "article_list_name".article_list_id
		LEFT JOIN "article_list_member" ON "article_list".id = "article_list_member".article_list_id
		LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		JOIN "observed_object" ON "article_list".id = "observed_object".article_list_id AND "observed_object".user_id = $2
		WHERE "article_list".organization_id = $1
		AND (public IS TRUE OR "article_list_member".user_id = $2 OR "user_group_member".user_id = $2)
		` + fmt.Sprintf(articleListSelectNameQuery, 5, 5, 5) + `
		LIMIT $4 OFFSET $3`
	var entities []ArticleList

	if err := connection.Select(&entities, query, orgaId, userId, offset, n, langId); err != nil {
		logbuch.Error("Article list by organization id and user id and language id and observed not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "offset": offset, "n": n})
		return nil
	}

	return entities
}

func FindArticleListsByOrganizationIdAndUserIdAndLanguageIdAndPinnedLimit(orgaId, userId, langId hide.ID, offset, n int) []ArticleList {
	query := ""
	params := make([]interface{}, 4)
	params[0] = orgaId

	if userId != 0 {
		params[1] = userId
		params[2] = offset
		params[3] = n
		params = append(params, langId)
		query = `SELECT DISTINCT ON ("article_list".id) "article_list".*,
			"article_list_name".name "name.name", "article_list_name".info "name.info",
			` + articleListEntryCountForUserQuery + `
			FROM "article_list"
			JOIN "article_list_name" ON "article_list".id = "article_list_name".article_list_id
			LEFT JOIN "article_list_member" ON "article_list".id = "article_list_member".article_list_id
			LEFT JOIN "user_group" ON "article_list_member".user_group_id = "user_group".id
			LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
			WHERE "article_list".organization_id = $1
			AND (public IS TRUE OR "article_list_member".user_id = $2 OR "user_group_member".user_id = $2)
			AND pinned IS TRUE
			` + fmt.Sprintf(articleListSelectNameQuery, 5, 5, 5) + `
			LIMIT $4 OFFSET $3`
	} else {
		params[1] = offset
		params[2] = n
		params[3] = langId
		query = `SELECT DISTINCT ON ("article_list".id) "article_list".*,
			"article_list_name".name "name.name", "article_list_name".info "name.info",
			` + articleListEntryCountForClientQuery + `
			FROM "article_list"
			JOIN "article_list_name" ON "article_list".id = "article_list_name".article_list_id
			WHERE "article_list".organization_id = $1
			AND pinned IS TRUE
			AND client_access IS TRUE
			` + fmt.Sprintf(articleListSelectNameQuery, 4, 4, 4) + `
			LIMIT $3 OFFSET $2`
	}

	var entities []ArticleList

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Article list by organization id and user id and language id and pinned not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId})
		return nil
	}

	return entities
}

func FindArticleListsByOrganizationId(orgaId hide.ID) []ArticleList {
	var entities []ArticleList

	if err := connection.Select(&entities, `SELECT * FROM "article_list" WHERE organization_id = $1`, orgaId); err != nil {
		logbuch.Error("Article list by organization id not found", logbuch.Fields{"err": err, "orga_id": orgaId})
		return nil
	}

	return entities
}

// FindArticleListsByOrganizationIdAndAccessByUserIdOnlyTx returns lists where only the given user has access to.
func FindArticleListsByOrganizationIdAndAccessByUserIdOnlyTx(tx *sqlx.Tx, orgaId, userId hide.ID) []ArticleList {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "article_list".*,` + articleListEntryCountForUserQuery + `
		FROM "article_list"
		JOIN "article_list_member" ON "article_list".id = "article_list_member".article_list_id
		WHERE organization_id = $1
		AND "article_list_member".user_id = $2
		AND "article_list_member".is_moderator IS TRUE
		AND 1 = (SELECT COUNT(1) FROM "article_list_member"
		WHERE article_list_id = "article_list".id
		AND "article_list_member".is_moderator IS TRUE)`
	var entities []ArticleList

	if err := tx.Select(&entities, query, orgaId, userId); err != nil {
		logbuch.Error("Article list by organization id and access user id only not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entities
}

func FindArticleListsByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfo(orgaId, userId, langId hide.ID, keywords string, filter *SearchArticleListFilter) []ArticleList {
	query, params := buildArticleListByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfoQuery(orgaId, userId, langId, keywords, filter, false)
	var entities []ArticleList

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading article lists by organization id and user id and name or info", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "keywords": keywords})
		return nil
	}

	return entities
}

func CountArticleListsByOrganizationIdAndUserIdAndClientAccessAndPinned(orgaId, userId hide.ID, clientAccess bool) int {
	query := `SELECT COUNT(DISTINCT "article_list".id) FROM "article_list"
		LEFT JOIN "article_list_member" ON "article_list".id = "article_list_member".article_list_id
		WHERE organization_id = $1
		AND (public IS TRUE OR "article_list_member".user_id = $2)
		AND pinned IS TRUE `

	if clientAccess {
		query += `AND client_access IS TRUE`
	}

	var count int

	if err := connection.Get(&count, query, orgaId, userId); err != nil {
		logbuch.Error("Error counting article lists by organization id and user id", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return 0
	}

	return count
}

func CountArticleListsByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfo(orgaId, userId hide.ID, keywords string, filter *SearchArticleListFilter) int {
	query, params := buildArticleListByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfoQuery(orgaId, userId, 0, keywords, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting article lists by organization id and user id and name or info", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "keywords": keywords})
		return 0
	}

	return count
}

func CountArticleListByOrganizationId(orgaId hide.ID) int {
	query := `SELECT COUNT(1) FROM "article_list" WHERE organization_id = $1`
	var count int

	if err := connection.Get(&count, query, orgaId); err != nil {
		logbuch.Error("Error counting article lists by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return 0
	}

	return count
}

func buildArticleListByOrganizationIdAndUserIdAndLanguageIdAndNameOrInfoQuery(orgaId, userId, langId hide.ID, keywords string, filter *SearchArticleListFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 1)
	params[0] = orgaId
	var sb strings.Builder

	if count {
		sb.WriteString(`SELECT COUNT(1) FROM (SELECT DISTINCT "article_list".id `)
	} else {
		sb.WriteString(`SELECT id, organization_id, public, client_access, article_count, def_time, mod_time,
			"name" "name.name", "info" "name.info"
			FROM (SELECT DISTINCT ON ("article_list".id) "article_list".*, "article_list_name"."name", "article_list_name"."info", `)

		if userId != 0 {
			sb.WriteString(articleListEntryCountForUserQuery)
		} else {
			sb.WriteString(articleListEntryCountForClientQuery)
		}
	}

	if userId != 0 {
		params = append(params, userId)
		sb.WriteString(articleListBaseForUserQuery)

		if keywords != "" {
			sb.WriteString(`AND (SIMILARITY("article_list_name".name, $3) > 0.2
			OR LOWER("article_list_name".name) LIKE LOWER('%'||$3||'%')
			OR SIMILARITY("article_list_name".info, $3) > 0.2
			OR LOWER("article_list_name".info) LIKE LOWER('%'||$3||'%'))`)
		}

		sb.WriteString(`AND (public IS TRUE OR "article_list_member".user_id = $2 OR "user_group_member".user_id = $2) `)
	} else {
		sb.WriteString(articleListBaseForClientQuery)

		if keywords != "" {
			sb.WriteString(`AND (SIMILARITY("article_list_name".name, $2) > 0.2
			OR LOWER("article_list_name".name) LIKE LOWER('%'||$2||'%')
			OR SIMILARITY("article_list_name".info, $2) > 0.2
			OR LOWER("article_list_name".info) LIKE LOWER('%'||$2||'%'))`)
		}

		sb.WriteString(`AND "article_list".client_access IS TRUE `)
	}

	if keywords != "" {
		params = append(params, keywords)
	}

	// add field filter (joined with "AND")
	index := len(params) + 1
	fieldFilter, index, params := filter.addFieldFilter("article_list_name", index, params, []string{filter.Name, filter.Info}, "name", "info")
	sb.WriteString(fieldFilter)

	if len(filter.UserIds) > 0 {
		var userFilter string
		userFilter, index, params = filter.filterInIds(filter.UserIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "article_list_member".user_id %v`, userFilter))
	}

	if len(filter.UserGroupIds) > 0 {
		var userGroupFilter string
		userGroupFilter, index, params = filter.filterInIds(filter.UserGroupIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "article_list_member".user_group_id %v`, userGroupFilter))
	}

	// add date filter
	dateFilter, index, params := filter.addDateFilter("article_list", index, params)
	sb.WriteString(dateFilter)

	if !count {
		// select name based on user preference and availability
		params = append(params, langId)
		sb.WriteString(fmt.Sprintf(articleListSelectNameQuery, index, index, index))
		index++

		// close inner query
		sb.WriteString(") AS result_set ")

		// add order (to outer query)
		defaultFields := []SortValue{{"id", sortDirectionASC}}
		sb.WriteString(filter.addSorting("result_set", defaultFields, SortValue{"name", filter.SortName}, SortValue{"info", filter.SortInfo}))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	} else {
		sb.WriteString(") AS list_count")
	}

	return sb.String(), params
}

func DeleteArticleListById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "bookmark" WHERE article_list_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting bookmark by article list id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "observed_object" WHERE article_list_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting observed object by article list id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_list_entry" WHERE article_list_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article list entries", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_list_member" WHERE article_list_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article list member", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_list_name" WHERE article_list_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article list names", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_list" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article list", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleList(tx *sqlx.Tx, entity *ArticleList) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_list" (organization_id,
			public,
			pinned,
			client_access)
			VALUES (:organization_id,
			:public,
			:pinned,
			:client_access) RETURNING id`,
		`UPDATE "article_list" SET organization_id = :organization_id,
			public = :public,
			pinned = :pinned,
			client_access = :client_access
			WHERE id = :id`)
}
