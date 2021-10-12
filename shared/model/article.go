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
	articleSelectNameQuery = `AND (
			(
				(SELECT EXISTS (SELECT 1 FROM article_content c WHERE c.article_id = article.id AND c.language_id = $%d)) IS TRUE
				AND article_content.language_id = $%d
			)
			OR (
				(SELECT EXISTS (SELECT 1 FROM article_content c WHERE c.article_id = article.id AND c.language_id = $%d)) IS FALSE
				AND (SELECT EXISTS (SELECT 1 FROM article_content c WHERE c.article_id = article.id AND c.language_id = (SELECT id FROM "language" WHERE "language".organization_id = $1 AND "language".default IS TRUE))) IS TRUE
				AND article_content.language_id = (SELECT id FROM "language" WHERE "language".organization_id = $1 AND "language".default IS TRUE)
			)
			OR (
				(SELECT EXISTS (SELECT 1 FROM article_content c WHERE c.article_id = article.id AND c.language_id = $%d)) IS FALSE
				AND (SELECT EXISTS (SELECT 1 FROM article_content c WHERE c.article_id = article.id AND c.language_id = (SELECT id FROM "language" WHERE "language".organization_id = $1 AND "language".default IS TRUE))) IS FALSE
			)
		)`
	articleContentFieldsQuery = `"article_content".id "latest_article_content.id",
		"article_content".title "latest_article_content.title",
		"article_content".version "latest_article_content.version",
		"article_content".commit "latest_article_content.commit",
		"article_content".wip "latest_article_content.wip",
		"article_content".reading_time "latest_article_content.reading_time",
		"article_content".schema_version "latest_article_content.schema_version",
		"article_content".def_time "latest_article_content.def_time",
		"article_content".mod_time "latest_article_content.mod_time",
		"article_content".article_id "latest_article_content.article_id",
		"article_content".language_id "latest_article_content.language_id",
		"article_content".user_id "latest_article_content.user_id"`
)

type Article struct {
	db.BaseEntity

	OrganizationId hide.ID     `db:"organization_id" json:"organization_id"`
	Views          uint        `json:"views"`
	WIP            int         `json:"wip"`
	ReadEveryone   bool        `db:"read_everyone" json:"read_everyone"`
	WriteEveryone  bool        `db:"write_everyone" json:"write_everyone"`
	Private        bool        `json:"private"`
	ClientAccess   bool        `db:"client_access" json:"client_access"`
	Archived       null.String `json:"archived"`
	Published      null.Time   `json:"published"`
	Pinned         bool        `json:"pinned"`

	LatestArticleContent *ArticleContent `db:"latest_article_content" json:"latest_article_content"`
	Access               []ArticleAccess `db:"-" json:"access"`
	Tags                 []Tag           `db:"-" json:"tags"`
	PreviewImage         string          `json:"preview_image"`

	Rank float32 `db:"rank" json:"-"`
}

func GetArticleByOrganizationIdAndIdAndPinned(orgaId, id hide.ID) *Article {
	entity := new(Article)

	if err := connection.Get(entity, `SELECT * FROM "article" WHERE organization_id = $1 AND id = $2 AND pinned IS TRUE`, orgaId, id); err != nil {
		logbuch.Debug("Article by organization id and id and pinned not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetArticleByOrganizationIdAndId(orgaId, id hide.ID) *Article {
	return GetArticleByOrganizationIdAndIdTx(nil, orgaId, id)
}

func GetArticleByOrganizationIdAndIdTx(tx *sqlx.Tx, orgaId, id hide.ID) *Article {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(Article)

	if err := tx.Get(entity, `SELECT * FROM "article" WHERE organization_id = $1 AND id = $2 AND archived IS NULL`, orgaId, id); err != nil {
		logbuch.Debug("Article by organization id and id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func GetArticleByOrganizationIdAndIdIgnoreArchived(orgaId, id hide.ID) *Article {
	return GetArticleByOrganizationIdAndIdIgnoreArchivedTx(nil, orgaId, id)
}

func GetArticleByOrganizationIdAndIdIgnoreArchivedTx(tx *sqlx.Tx, orgaId, id hide.ID) *Article {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(Article)

	if err := tx.Get(entity, `SELECT * FROM "article" WHERE organization_id = $1 AND id = $2`, orgaId, id); err != nil {
		logbuch.Debug("Article by organization id and id not found ignoring archived", logbuch.Fields{"err": err, "orga_id": orgaId, "id": id})
		return nil
	}

	return entity
}

func FindArticleByOrganizationIdAndUserIdAndLanguageIdAndUnpublishedWithLimit(orgaId, userId, langId hide.ID, offset, n int) []Article {
	query := `SELECT DISTINCT ON ("article".id) "article".*,
		` + articleContentFieldsQuery + `
		FROM "article"
		JOIN "article_content" ON "article".id = "article_content".article_id AND "article_content".version = 0
		JOIN "article_access" ON "article".id = "article_access".article_id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id 
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE "article".organization_id = $1
		AND ("article_access".user_id = $2 OR "user_group_member".user_id = $2)
		AND "article".published IS NULL
		` + fmt.Sprintf(articleSelectNameQuery, 5, 5, 5, 5) + `
		LIMIT $4 OFFSET $3`
	var entities []Article

	if err := connection.Select(&entities, query, orgaId, userId, offset, n, langId); err != nil {
		logbuch.Error("Article by organization id and user id and language id and unpublished not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "offset": offset, "n": n})
		return nil
	}

	return joinArticleAuthors(entities)
}

func FindArticleByOrganizationIdAndUserIdAndLanguageIdAndPrivateWithLimit(orgaId, userId, langId hide.ID, offset, n int) []Article {
	// no indirect access through group here!
	query := `SELECT DISTINCT ON ("article".id) "article".*,
		` + articleContentFieldsQuery + `
		FROM "article"
		JOIN "article_content" ON "article".id = "article_content".article_id AND "article_content".version = 0
		JOIN "article_access" ON "article".id = "article_access".article_id
		WHERE "article".organization_id = $1
		AND "article_access".user_id = $2
		AND "article".private IS TRUE
		` + fmt.Sprintf(articleSelectNameQuery, 5, 5, 5, 5) + `
		LIMIT $4 OFFSET $3`
	var entities []Article

	if err := connection.Select(&entities, query, orgaId, userId, offset, n, langId); err != nil {
		logbuch.Error("Article by organization id and user id and language id and private not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "offset": offset, "n": n})
		return nil
	}

	return joinArticleAuthors(entities)
}

func FindArticleByOrganizationIdAndUserIdAndLanguageIdAndObservedWithLimit(orgaId, userId, langId hide.ID, offset, n int) []Article {
	query := `SELECT DISTINCT ON ("article".id) "article".*,
		` + articleContentFieldsQuery + `
		FROM "article"
		JOIN "article_content" ON "article".id = "article_content".article_id AND "article_content".version = 0
		JOIN "article_access" ON "article".id = "article_access".article_id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id 
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		JOIN "observed_object" ON "article".id = "observed_object".article_id AND "observed_object".user_id = $2
		WHERE "article".organization_id = $1
		AND ("article_access".user_id = $2 OR "user_group_member".user_id = $2)
		` + fmt.Sprintf(articleSelectNameQuery, 5, 5, 5, 5) + `
		LIMIT $4 OFFSET $3`
	var entities []Article

	if err := connection.Select(&entities, query, orgaId, userId, offset, n, langId); err != nil {
		logbuch.Error("Article by organization id and user id and language id and observed not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "offset": offset, "n": n})
		return nil
	}

	return joinArticleAuthors(entities)
}

// FindArticleByOrganizationIdAndAccessByUserIdOnlyTx returns articles where only the given user has access to and no one else.
func FindArticleByOrganizationIdAndAccessByUserIdOnlyTx(tx *sqlx.Tx, orgaId, userId hide.ID) []Article {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT "article".* FROM "article"
		JOIN "article_access" ON "article".id = "article_access".article_id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id 
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id
		WHERE "article".organization_id = $1
		AND "article".write_everyone IS FALSE
		AND ("article_access".user_id = $2 OR "user_group_member".user_id = $2)
		AND 1 = (SELECT COUNT(1) FROM "article_access" WHERE article_id = "article".id AND "write" IS TRUE)`
	var entities []Article

	if err := tx.Select(&entities, query, orgaId, userId); err != nil {
		logbuch.Error("Article by organization id and user id and language id and pinned not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return nil
	}

	return entities
}

func FindArticleByOrganizationIdAndUserIdAndLanguageIdAndClientAccessAndPinnedLimit(orgaId, userId, langId hide.ID, clientAccess bool, offset, n int) []Article {
	query := `SELECT DISTINCT ON ("article".id) "article".*,
		` + articleContentFieldsQuery + `
		FROM "article"
		JOIN "article_content" ON "article".id = "article_content".article_id AND "article_content".version = 0
		LEFT JOIN "article_access" ON "article".id = "article_access".article_id
		WHERE organization_id = $1
		AND (read_everyone IS TRUE OR write_everyone IS TRUE OR "article_access".user_id = $2)
		AND pinned IS TRUE
		AND archived IS NULL ` + fmt.Sprintf(articleSelectNameQuery, 5, 5, 5, 5)

	if clientAccess {
		query += `AND client_access IS TRUE `
	}

	query += `LIMIT $4 OFFSET $3`
	var entities []Article

	if err := connection.Select(&entities, query, orgaId, userId, offset, n, langId); err != nil {
		logbuch.Error("Article by organization id and user id and language id and pinned not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId})
		return nil
	}

	return joinArticleAuthors(entities)
}

func FindArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimit(orgaId, userId, langId hide.ID, keywords string, filter *SearchArticleFilter) []Article {
	query, params := buildArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimitQuery(orgaId, userId, langId, keywords, filter, false)
	var entities []Article

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Article by organization id and user id and query and filter with limit not found", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "keywords": keywords, "filter": filter})
		return nil
	}

	return joinArticleAuthors(entities)
}

func CountArticleByOrganizationIdAndUserIdAndClientAccessAndPinned(orgaId, userId hide.ID, clientAccess bool) int {
	query := `SELECT COUNT(DISTINCT "article".id) FROM "article"
		LEFT JOIN "article_access" ON "article".id = "article_access".article_id
		WHERE organization_id = $1
		AND (read_everyone IS TRUE OR write_everyone IS TRUE OR "article_access".user_id = $2)
		AND pinned IS TRUE
		AND archived IS NULL `

	if clientAccess {
		query += `AND client_access IS TRUE`
	}

	var count int

	if err := connection.Get(&count, query, orgaId, userId); err != nil {
		logbuch.Error("Error counting articles by organization id and user id", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
		return 0
	}

	return count
}

func CountArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimit(orgaId, userId hide.ID, keywords string, filter *SearchArticleFilter) int {
	query, params := buildArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimitQuery(orgaId, userId, 0, keywords, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting articles by organization id and user id and query and filter with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "keywords": keywords, "filter": filter})
		return 0
	}

	return count
}

func CountArticleByOrganizationId(orgaId hide.ID) int {
	query := `SELECT COUNT(1) FROM "article" WHERE organization_id = $1`
	var count int

	if err := connection.Get(&count, query, orgaId); err != nil {
		logbuch.Error("Error counting articles by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return 0
	}

	return count
}

func buildArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimitQuery(orgaId, userId, langId hide.ID, keywords string, filter *SearchArticleFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 2)
	params[0] = orgaId
	params[1] = userId
	var sb strings.Builder

	if count {
		sb.WriteString(`SELECT COUNT(DISTINCT(id)) FROM (SELECT "article".id `)
	} else {
		sb.WriteString(`SELECT id, organization_id, views, wip, read_everyone, write_everyone, private, client_access, archived, published, def_time, mod_time,
		article_content_id "latest_article_content.id",
		title "latest_article_content.title",
		version "latest_article_content.version",
		commit "latest_article_content.commit",
		article_content_wip "latest_article_content.wip",
		article_content_def_time "latest_article_content.def_time",
		article_content_mod_time "latest_article_content.mod_time",
		article_content_article_id "latest_article_content.article_id",
		article_content_language_id "latest_article_content.language_id",
		article_content_user_id "latest_article_content.user_id",
		title_rank + content_rank*0.2 "rank"
		FROM (`)

		sb.WriteString(`SELECT DISTINCT ON (id) "article".*,
			"article_content".id "article_content_id",
			"article_content".title,
			"article_content".version,
			"article_content".commit,
			"article_content".wip AS "article_content_wip",
			"article_content".def_time AS "article_content_def_time",
			"article_content".mod_time AS "article_content_mod_time",
			"article_content".article_id AS "article_content_article_id",
			"article_content".language_id AS "article_content_language_id",
			"article_content".user_id AS "article_content_user_id", `)

		if keywords != "" {
			sb.WriteString(`ts_rank_cd("article_content".title_tsvector, to_tsquery($4)) AS title_rank,
				ts_rank_cd("article_content".content_tsvector, to_tsquery($4), 1) AS content_rank `)
		} else {
			sb.WriteString(`1 title_rank, 1 content_rank `)
		}
	}

	sb.WriteString(`FROM "article"
		JOIN "article_content" ON "article".id = "article_content".article_id AND "article_content".version = 0
		JOIN "article_access" ON "article".id = "article_access".article_id
		LEFT JOIN "user_group" ON "article_access".user_group_id = "user_group".id 
		LEFT JOIN "user_group_member" ON "user_group".id = "user_group_member".user_group_id `)

	if len(filter.AuthorUserIds) > 0 {
		sb.WriteString(`LEFT JOIN "article_content_author" ON "article_content_author".article_content_id IN (SELECT id FROM "article_content" ac WHERE ac.article_id = "article".id)
			LEFT JOIN "user" ON "article_content_author".user_id = "user".id
			LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1 `)
	}

	if len(filter.TagIds) > 0 {
		sb.WriteString(`LEFT JOIN "article_tag" ON "article".id = "article_tag".article_id `)
	}

	sb.WriteString(`WHERE "article".organization_id = $1 `)

	if keywords != "" {
		params = append(params, keywords)
		params = append(params, db.ToTSVector(keywords))
		sb.WriteString(`AND ("article_content".title_tsvector @@ to_tsquery($4)
			OR SIMILARITY("article_content".title, $3) > 0.2
			OR LOWER("article_content".title) LIKE LOWER('%'||$3||'%')
			OR "article_content".content_tsvector @@ to_tsquery($4)
			OR EXISTS (
				SELECT 1 FROM "article_tag"
				JOIN "tag" ON "article_tag".tag_id = "tag".id
				WHERE "article_tag".article_id = "article".id
				AND (SIMILARITY("tag"."name", $3) > 0.2 OR LOWER("tag"."name") LIKE LOWER('%'||$3||'%'))
			)
			OR SIMILARITY("article_content"."commit", $3) > 0.2
			OR LOWER("article_content"."commit") LIKE LOWER('%'||$3||'%')) `)
	}

	// add field filter (joined with "AND")
	if filter.Archived {
		sb.WriteString(`AND "article".archived IS NOT NULL `)
	}

	if filter.ClientAccess {
		sb.WriteString(`AND "article".client_access IS TRUE `)
	}

	if filter.WIP {
		sb.WriteString(`AND "article".wip != -1 `)
	} // don't filter for the other case, because we don't exclude wip articles in general

	index := len(params) + 1
	fieldFilter, index, params := filter.addFieldFilter("article_content", index, params, []string{filter.Commits}, "commit")
	sb.WriteString(fieldFilter)
	tsvectorFilter, index, params := filter.addTSVectorFieldFilter("article_content", index, params, []string{filter.Title, filter.Content}, "title_tsvector", "content_tsvector")
	sb.WriteString(tsvectorFilter)
	tags := strings.TrimSpace(filter.Tags)

	if tags != "" {
		sb.WriteString(fmt.Sprintf(`AND EXISTS (SELECT 1 FROM "article_tag" JOIN "tag" ON "article_tag".tag_id = "tag".id WHERE "article_tag".article_id = "article".id AND "tag"."name" %% $%v) `, index))
		params = append(params, tags)
		index++
	}

	if len(filter.TagIds) > 0 {
		var tagFilter string
		tagFilter, index, params = filter.filterInIds(filter.TagIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "article_tag".tag_id %v`, tagFilter))
	}

	if len(filter.AuthorUserIds) > 0 {
		var authorFilter string
		authorFilter, index, params = filter.filterInIds(filter.AuthorUserIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "organization_member".user_id %v`, authorFilter))
	}

	if len(filter.UserGroupIds) > 0 {
		var userGroupFilter string
		userGroupFilter, index, params = filter.filterInIds(filter.UserGroupIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "article_access".user_group_id %v`, userGroupFilter))
	}

	// add date filter
	dateFilter, index, params := filter.addDateFilter("article", index, params)
	sb.WriteString(dateFilter)

	if !filter.PublishedStart.IsZero() {
		sb.WriteString(fmt.Sprintf(`AND "article".published > $%v `, index))
		params = append(params, filter.PublishedStart)
		index++
	}

	if !filter.PublishedEnd.IsZero() {
		sb.WriteString(fmt.Sprintf(`AND "article".published < $%v `, index))
		params = append(params, filter.PublishedEnd)
		index++
	}

	// check access
	sb.WriteString(`AND ("article".read_everyone IS TRUE
		OR "article".write_everyone IS TRUE
		OR "article_access".user_id = $2
		OR "user_group_member".user_id = $2) `)

	if !count {
		// select name based on user preference and availability
		params = append(params, langId)
		sb.WriteString(fmt.Sprintf(articleSelectNameQuery, index, index, index, index))
		index++

		// close distinct select
		sb.WriteString(") AS result_set ")

		// sorting
		if filter.SortPublished != "" ||
			filter.SortCreated != "" ||
			filter.SortUpdated != "" ||
			filter.SortTitle != "" {
			sb.WriteString(filter.addSorting("result_set", nil, SortValue{`"result_set".title`, filter.SortTitle}, SortValue{`"result_set".published`, filter.SortPublished}, SortValue{`"result_set"."article_content_mod_time"`, filter.SortUpdated}))
		} else {
			rankDirection := sortDirectionDESC

			if strings.ToUpper(filter.SortRelevance) == sortDirectionASC {
				rankDirection = sortDirectionASC
			}

			defaultFields := []SortValue{{"rank", rankDirection}}
			sb.WriteString(filter.addSorting("result_set", defaultFields, SortValue{`"result_set".title`, filter.SortTitle}, SortValue{`"result_set".published`, filter.SortPublished}, SortValue{`"result_set"."article_content_mod_time"`, filter.SortUpdated}))
		}

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	} else {
		sb.WriteString(") AS article_count")
	}

	return sb.String(), params
}

func DeleteArticleById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "file" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article files by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_recommendation" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article recommendations by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_visit" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article visit by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_access" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article access by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "bookmark" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting bookmark by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "observed_object" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting observed object by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_content_author"
		WHERE article_content_id IN (SELECT id FROM "article_content"
		WHERE "article_content".id = "article_content_author".article_content_id
		AND "article_content".article_id = $1)`, id)

	if err != nil {
		logbuch.Error("Error deleting article tag by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_content" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article content by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_list_entry" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article list entry by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article_tag" WHERE article_id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article tag by article id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	_, err = tx.Exec(`DELETE FROM "article" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticle(tx *sqlx.Tx, entity *Article) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article" (organization_id,
			views,
			wip,
			read_everyone,
			write_everyone,
			private,
			client_access,
			archived,
			published,
			pinned)
			VALUES (:organization_id,
			:views,
			:wip,
			:read_everyone,
			:write_everyone,
			:private,
			:client_access,
			:archived,
			:published,
			:pinned) RETURNING id`,
		`UPDATE "article" SET organization_id = :organization_id,
			views = :views,
			wip = :wip,
			read_everyone = :read_everyone,
			write_everyone = :write_everyone,
			private = :private,
			client_access = :client_access,
			archived = :archived,
			published = :published,
			pinned = :pinned
			WHERE id = :id`)
}

func joinArticleAuthors(entities []Article) []Article {
	for i := range entities {
		entities[i].LatestArticleContent.Authors = FindArticleContentAuthorUserByArticleIdTx(nil, entities[i].ID)
	}

	return entities
}
