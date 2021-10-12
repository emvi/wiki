package model

import (
	"emviwiki/shared/db"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ArticleListEntry struct {
	db.BaseEntity

	ArticleListId hide.ID `db:"article_list_id" json:"article_list_id"`
	ArticleId     hide.ID `db:"article_id" json:"article_id"`
	Position      uint    `json:"position"`

	Article *Article `db:"article" json:"article"`
}

func GetArticleListEntryByArticleListIdAndArticleId(listId, articleId hide.ID) *ArticleListEntry {
	return GetArticleListEntryByArticleListIdAndArticleIdTx(nil, listId, articleId)
}

func GetArticleListEntryByArticleListIdAndArticleIdTx(tx *sqlx.Tx, listId, articleId hide.ID) *ArticleListEntry {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	entity := new(ArticleListEntry)

	if err := tx.Get(entity, `SELECT * FROM "article_list_entry"
		WHERE article_list_id = $1
		AND article_id = $2`, listId, articleId); err != nil {
		logbuch.Debug("Article list entry by article list id and article id not found", logbuch.Fields{"err": err, "list_id": listId, "article_id": articleId})
		return nil
	}

	return entity
}

func GetArticleListEntryLastPositionByArticleListIdTx(tx *sqlx.Tx, listId hide.ID) uint {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var pos uint

	if err := tx.Get(&pos, `SELECT position FROM "article_list_entry"
		WHERE article_list_id = $1
		ORDER BY position DESC
		LIMIT 1`, listId); err != nil {
		logbuch.Debug("Article list entry last position by article list id not found", logbuch.Fields{"err": err, "list_id": listId})
		return 0
	}

	return pos
}

func GetArticleListEntryByArticleListIdAndPosition(listId hide.ID, pos uint) *ArticleListEntry {
	entity := new(ArticleListEntry)

	if err := connection.Get(entity, `SELECT * FROM "article_list_entry"
		WHERE article_list_id = $1
		AND position = $2`, listId, pos); err != nil {
		logbuch.Debug("Article list entry by article list id and position not found", logbuch.Fields{"err": err, "list_id": listId, "pos": pos})
		return nil
	}

	return entity
}

func CountArticleListEntryByArticleListId(listId hide.ID) int {
	query := `SELECT COUNT(1) FROM "article_list_entry" WHERE article_list_id = $1`
	var count int

	if err := connection.Get(&count, query, listId); err != nil {
		logbuch.Error("Error counting article list entries by article list id", logbuch.Fields{"err": err, "list_id": listId})
		return 0
	}

	return count
}

func CountArticleListEntryByArticleListIdAndPositionBefore(listId hide.ID, pos uint) int {
	query := `SELECT COUNT(1) FROM (
			SELECT 1
			FROM "article_list_entry"
			WHERE article_list_id = $1
			AND "position" < $2
			ORDER BY "position" ASC
		) AS result_count`
	var count int

	if err := connection.Get(&count, query, listId, pos); err != nil {
		logbuch.Error("Error counting article list entries by article list id and position before", logbuch.Fields{"err": err, "list_id": listId, "pos": pos})
		return 0
	}

	return count
}

func FindArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orgaId, userId, langId, listId hide.ID, filter *SearchArticleListEntryFilter) []Article {
	query, params := buildArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdAndFilterLimitQuery(orgaId, userId, listId, filter, false)
	var entities []Article

	if err := connection.Select(&entities, query, params...); err != nil {
		logbuch.Error("Error reading article list entries by article list id and user id and filter with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "list_id": listId, "filter": filter})
		return nil
	}

	for i := range entities {
		entities[i].LatestArticleContent = GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageId(orgaId, entities[i].ID, langId, true)
	}

	return entities
}

func CountArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdLimit(orgaId, userId, langId, listId hide.ID, filter *SearchArticleListEntryFilter) int {
	query, params := buildArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdAndFilterLimitQuery(orgaId, userId, listId, filter, true)
	var count int

	if err := connection.Get(&count, query, params...); err != nil {
		logbuch.Error("Error counting article list entries by article list id and user id and filter with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "lang_id": langId, "list_id": listId, "filter": filter})
		return 0
	}

	return count
}

func buildArticleListEntryArticlesByOrganizationIdAndUserIdAndLanguageIdArticleListIdAndFilterLimitQuery(orgaId, userId, listId hide.ID, filter *SearchArticleListEntryFilter, count bool) (string, []interface{}) {
	params := make([]interface{}, 3)
	params[0] = listId
	params[1] = userId
	params[2] = orgaId

	defaultFields := []SortValue{{`"article_list_entry".position`, sortDirectionASC}}
	customFields := []SortValue{{`"article_content".title`, filter.SortTitle}, {`"article_list_entry"."position"`, filter.SortPosition}}
	selectOrCount := `COUNT(DISTINCT("article".id))`

	if !count {
		sortColumns := filter.getSortColumns("article", defaultFields, customFields...)
		selectOrCount = `DISTINCT ON ("article".id,` + sortColumns + `) "article".*`
	}

	var sb strings.Builder
	sb.WriteString(`SELECT `)
	sb.WriteString(selectOrCount)
	sb.WriteString(` FROM "article_list_entry"
		JOIN "article" ON "article_list_entry".article_id = "article".id
		JOIN "article_content" ON "article".id = "article_content".article_id AND "article_content".version = 0
		LEFT JOIN "article_content_author" ON "article_content".id = "article_content_author".article_content_id
		LEFT JOIN "user" ON "article_content_author".user_id = "user".id
		LEFT JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $3
		WHERE article_list_id = $1`)

	// add archived filter
	if filter.Archived {
		sb.WriteString(`AND "article".archived IS NOT NULL`)
	}

	sb.WriteString(` AND ("article".write_everyone IS TRUE
			OR "article".read_everyone IS TRUE
			OR EXISTS (SELECT 1 FROM "article_access"
			LEFT JOIN "user_group_member" ON "article_access".user_group_id = "user_group_member".user_group_id
			WHERE article_id = "article".id
			AND ("article_access".user_id = $2 OR "user_group_member".user_id = $2))) `)

	// add field filter (joined with "AND")
	if filter.ClientAccess {
		sb.WriteString(`AND "article".client_access IS TRUE `)
	}

	index := 4
	fieldFilter, index, params := filter.addFieldFilter("article_content", index, params, []string{filter.Title, filter.Content, filter.Commits}, "title", "content", "commit")
	sb.WriteString(fieldFilter)
	tags := strings.TrimSpace(filter.Tags)

	if tags != "" {
		sb.WriteString(fmt.Sprintf(`AND $%v = ANY("article".tags) `, index))
		params = append(params, tags)
		index++
	}

	if len(filter.AuthorUserIds) > 0 {
		var authorFilter string
		authorFilter, index, params = filter.filterInIds(filter.AuthorUserIds, index, params)
		sb.WriteString(fmt.Sprintf(`AND "organization_member".user_id %v`, authorFilter))
	}

	// add date filter
	dateFilter, index, params := filter.addDateFilter("article", index, params)
	sb.WriteString(dateFilter)

	if !count {
		// sorting
		sb.WriteString(filter.addSorting("article", defaultFields, customFields...))

		// set limit
		var limit string
		limit, _, params = filter.addLimit(index, params)
		sb.WriteString(limit)
	}

	return sb.String(), params
}

func UpdateArticleListEntryPositionByArticleListIdTx(tx *sqlx.Tx, listId hide.ID) error {
	query := `UPDATE article_list_entry e
		SET "position" = rn.pos
		FROM (
			SELECT id, ROW_NUMBER() OVER (ORDER BY "position") AS pos
			FROM article_list_entry
			WHERE article_list_id = $1
		) rn
		WHERE e.id = rn.id`

	if _, err := tx.Exec(query, listId); err != nil {
		logbuch.Error("Error updating article list entries by article list id", logbuch.Fields{"err": err, "list_id": listId})
		return err
	}

	return nil
}

func DeleteArticleListEntryById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "article_list_entry" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting article list entry by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleListEntry(tx *sqlx.Tx, entity *ArticleListEntry) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_list_entry" (article_list_id,
			article_id,
			position)
			VALUES (:article_list_id,
			:article_id,
			:position) RETURNING id`,
		`UPDATE "article_list_entry" SET article_list_id = :article_list_id,
			article_id = :article_id,
			position = :position
			WHERE id = :id`)

}
