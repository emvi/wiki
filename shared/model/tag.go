package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	tagUsageCountForUserQuery = `(SELECT COUNT(DISTINCT "article_tag".id) FROM "article_tag"
		JOIN "article" ON "article_tag".article_id = "article".id
		JOIN "article_access" ON "article_tag".article_id = "article_access".article_id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE "article_tag".tag_id = "tag".id
		AND ("article".read_everyone IS TRUE
		OR "article".write_everyone IS TRUE
		OR "article_access".user_id = $2
		OR "user_group_member".user_id = $2)
		) AS "usages"`
	tagUsageCountForClientQuery = `(SELECT COUNT(DISTINCT "article_tag".id) FROM "article_tag"
		JOIN "article" ON "article_tag".article_id = "article".id
		WHERE "article_tag".tag_id = "tag".id
		AND "article".client_access IS TRUE
		) AS "usages"`
)

type Tag struct {
	db.BaseEntity

	OrganizationId hide.ID `db:"organization_id" json:"organization_id"`
	Name           string  `json:"name"`

	Usages int `json:"usages"`
}

func GetTagByOrganizationIdAndId(orgaId, id hide.ID) *Tag {
	entity := new(Tag)

	if err := connection.Get(entity, `SELECT * FROM "tag" WHERE organization_id = $1 AND id = $2`, orgaId, id); err != nil {
		logbuch.Debug("Tag by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetTagByOrganizationIdAndName(orgaId hide.ID, name string) *Tag {
	entity := new(Tag)

	if err := connection.Get(entity, `SELECT * FROM "tag" WHERE organization_id = $1 AND LOWER(name) = LOWER($2)`, orgaId, name); err != nil {
		logbuch.Debug("Tag by organization id and name not found", logbuch.Fields{"err": err, "orga_id": orgaId, "name": name})
		return nil
	}

	return entity
}

func GetTagByOrganizationIdAndUserIdAndId(orgaId, userId, id hide.ID) *Tag {
	query := `SELECT DISTINCT ON ("tag".id) "tag".* FROM "tag"
		JOIN "article_tag" ON "tag".id = "article_tag".tag_id
		JOIN "article" ON "article_tag".article_id = "article".id
		JOIN "article_access" ON "article_tag".article_id = "article_access".article_id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE "tag"."organization_id" = $1
		AND "tag".id = $3
		AND ("article".client_access IS TRUE
		OR "article".read_everyone IS TRUE
		OR "article".write_everyone IS TRUE
		OR "article_access".user_id = $2
		OR "user_group_member".user_id = $2)`
	entity := new(Tag)

	if err := connection.Get(entity, query, orgaId, userId, id); err != nil {
		logbuch.Debug("Tag by organization id and user id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "id": id})
		return nil
	}

	return entity
}

func FindTagByOrganizationIdAndUserIdAndArticleId(orgaId, userId, articleId hide.ID) []Tag {
	return FindTagByOrganizationIdAndUserIdAndArticleIdTx(nil, orgaId, userId, articleId)
}

func FindTagByOrganizationIdAndUserIdAndArticleIdTx(tx *sqlx.Tx, orgaId, userId, articleId hide.ID) []Tag {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := ""
	params := make([]interface{}, 2)
	params[0] = orgaId

	if userId != 0 {
		params[1] = userId
		params = append(params, articleId)
		query = `SELECT "tag".*, ` + tagUsageCountForUserQuery + ` FROM "tag"
			JOIN "article_tag" ON "tag".id = "article_tag".tag_id
			WHERE organization_id = $1
			AND article_id = $3`
	} else {
		params[1] = articleId
		query = `SELECT "tag".*, ` + tagUsageCountForClientQuery + ` FROM "tag"
			JOIN "article_tag" ON "tag".id = "article_tag".tag_id
			WHERE organization_id = $1
			AND article_id = $2`
	}

	var entities []Tag

	if err := tx.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading tags by organization id and user id and article id", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "article_id": articleId})
		return nil
	}

	return entities
}

func FindTagByOrganizationIdAndTagLimit(orgaId, userId hide.ID, keywords string, filter *SearchTagFilter) []Tag {
	query, params := buildTagByOrganizationIdAndTagLimitQuery(orgaId, userId, keywords, filter, false)
	var entities []Tag

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading tags by organization id and query and filter", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "keywords": keywords, "filter": filter})
		return nil
	}

	return entities
}

func CountTagByOrganizationIdAndTagLimit(orgaId, userId hide.ID, keywords string, filter *SearchTagFilter) int {
	query, params := buildTagByOrganizationIdAndTagLimitQuery(orgaId, userId, keywords, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting tags by organization id and query and filter", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "keywords": keywords, "filter": filter})
		return 0
	}

	return count
}

func CountTagByOrganizationId(orgaId hide.ID) int {
	query := `SELECT COUNT(1) FROM "tag" WHERE organization_id = $1`
	var count int

	if err := connection.Get(&count, query, orgaId); err != nil {
		logbuch.Error("Error counting tags by organization id", logbuch.Fields{"err": err, "organization_id": orgaId})
		return 0
	}

	return count
}

func buildTagByOrganizationIdAndTagLimitQuery(orgaId, userId hide.ID, keywords string, filter *SearchTagFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 1)
	params[0] = orgaId
	defaultFields := []SortValue{{`"tag"."name"`, sortDirectionASC}}
	customFields := []SortValue{{`"usages"`, filter.SortUsages}, {`"tag"."name"`, filter.SortName}}
	selectOrCount := `COUNT(DISTINCT("tag".id))`

	if !count {
		sortColumns := filter.getSortColumns("tag", defaultFields, customFields...)
		selectOrCount = `DISTINCT ON ("tag".id,` + sortColumns + `) "tag".*, `

		if userId != 0 {
			selectOrCount += tagUsageCountForUserQuery
		} else {
			selectOrCount += tagUsageCountForClientQuery
		}
	}

	var sb strings.Builder

	if userId != 0 {
		params = append(params, userId)

		// if requested by a user check access to tag through article direct or indirect access
		sb.WriteString(`SELECT `)
		sb.WriteString(selectOrCount)
		sb.WriteString(` FROM "tag"
			JOIN "article_tag" ON "tag".id = "article_tag".tag_id
			JOIN "article" ON "article_tag".article_id = "article".id
			JOIN "article_access" ON "article_tag".article_id = "article_access".article_id
			LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id
			LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
			WHERE "tag"."organization_id" = $1 `)

		if keywords != "" {
			params = append(params, keywords)
			sb.WriteString(`AND ($3 = '' OR SIMILARITY("tag"."name", $3) > 0.2 OR LOWER("tag"."name") LIKE LOWER('%'||$3||'%')) `)
		}

		sb.WriteString(`AND ("article".read_everyone IS TRUE
			OR "article".write_everyone IS TRUE
			OR "article_access".user_id = $2
			OR "user_group_member".user_id = $2) `)
	} else {
		// if requested by a client check the article has set client access to true
		sb.WriteString(`SELECT `)
		sb.WriteString(selectOrCount)
		sb.WriteString(` FROM "tag"
			JOIN "article_tag" ON "tag".id = "article_tag".tag_id
			JOIN "article" ON "article_tag".article_id = "article".id
			WHERE "tag"."organization_id" = $1 `)

		if keywords != "" {
			params = append(params, keywords)
			sb.WriteString(`AND ($2 = '' OR SIMILARITY("tag"."name", $2) > 0.2 OR LOWER("tag"."name") LIKE LOWER('%'||$2||'%')) `)
		}

		sb.WriteString(`AND "article".client_access IS TRUE `)
	}

	// add date filter
	index := len(params) + 1
	dateFilter, index, params := filter.addDateFilter("tag", index, params)
	sb.WriteString(dateFilter)

	if !count {
		// sorting
		sb.WriteString(filter.addSorting("tag", defaultFields, SortValue{"usages", filter.SortUsages}, SortValue{"name", filter.SortName}))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	}

	return sb.String(), params
}

func DeleteTagByOrganizationIdAndId(tx *sqlx.Tx, orgaId, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "tag"
		WHERE organization_id = $1
		AND id = $2`, orgaId, id)

	if err != nil {
		logbuch.Error("Error deleting tag by organization id and id", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func DeleteTagUnusedByOrganizationId(tx *sqlx.Tx, orgaId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "tag"
		WHERE organization_id = $1
		AND NOT EXISTS (SELECT 1 FROM "article_tag" WHERE tag_id = "tag".id)`, orgaId)

	if err != nil {
		logbuch.Error("Error deleting unused tag by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveTag(tx *sqlx.Tx, entity *Tag) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "tag" (organization_id, name)
			VALUES (:organization_id, :name) RETURNING id`,
		`UPDATE "tag" SET organization_id = :organization_id,
			name = :name
			WHERE id = :id`)
}
