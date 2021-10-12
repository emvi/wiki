package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
)

const (
	articleContentWithoutContentQuery = `SELECT id, article_id, language_id, user_id, title, version, commit, wip, reading_time, schema_version, def_time, mod_time FROM "article_content" `
)

type ArticleContent struct {
	db.BaseEntity

	Title           string      `json:"title"`
	Content         string      `json:"content"`
	Version         int         `json:"version"`
	Commit          null.String `json:"commit"`
	WIP             bool        `json:"wip"`
	ContentTsvector string      `db:"content_tsvector" json:"-"`
	TitleTsvector   string      `db:"title_tsvector" json:"-"`
	ReadingTime     int         `db:"reading_time" json:"reading_time"` // seconds
	SchemaVersion   int         `db:"schema_version" json:"-"`
	RTL             bool        `json:"rtl"` // right to left

	ArticleId  hide.ID `db:"article_id" json:"article_id"`
	LanguageId hide.ID `db:"language_id" json:"language_id"`
	UserId     hide.ID `db:"user_id" json:"user_id"` // user who created this commit

	User    *User  `db:"user" json:"user"`
	Authors []User `db:"-" json:"authors"`
}

func GetArticleContentById(id hide.ID) *ArticleContent {
	entity := new(ArticleContent)

	if err := connection.Get(entity, `SELECT * FROM "article_content" WHERE id = $1`, id); err != nil {
		logbuch.Debug("Article content by id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func GetArticleContentLastByArticleIdAndLanguageIdAndWIP(articleId, langId hide.ID, includeWIP bool) *ArticleContent {
	return GetArticleContentLastByArticleIdAndLanguageIdAndWIPTx(nil, articleId, langId, includeWIP)
}

func GetArticleContentLastByArticleIdAndLanguageIdAndWIPTx(tx *sqlx.Tx, articleId, langId hide.ID, includeWIP bool) *ArticleContent {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT * FROM "article_content" WHERE article_id = $1 AND language_id = $2 AND version != 0`

	if !includeWIP {
		query += ` AND wip IS FALSE`
	}

	query += ` ORDER BY version DESC LIMIT 1`

	entity := new(ArticleContent)

	if err := tx.Get(entity, query, articleId, langId); err != nil {
		logbuch.Debug("Last article content commit by article id and language id not found", logbuch.Fields{"err": err, "article_id": articleId, "language_id": langId})
		return nil
	}

	return entity
}

func GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageId(orgaId, articleId, langId hide.ID, withContent bool) *ArticleContent {
	return GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageIdTx(nil, orgaId, articleId, langId, withContent)
}

// Returns the latest article content for a given organization, article and language.
// If the provided language is not available, the selection will fall back to the default language.
// If the content for the default language is not available, it will use any language that's available.
func GetArticleContentLatestByOrganizationIdAndArticleIdAndLanguageIdTx(tx *sqlx.Tx, orgaId, articleId, langId hide.ID, withContent bool) *ArticleContent {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var entity *ArticleContent

	if langId != 0 {
		entity = GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, articleId, langId, withContent)
	}

	if entity == nil {
		defaultLang := GetDefaultLanguageByOrganizationIdTx(tx, orgaId)
		entity = GetArticleContentLatestByArticleIdAndLanguageIdTx(tx, articleId, defaultLang.ID, withContent)
	}

	if entity == nil {
		query := `SELECT * FROM "article_content" WHERE article_id = $1 AND version = 0 LIMIT 1`

		if !withContent {
			query = articleContentWithoutContentQuery + `WHERE article_id = $1 AND version = 0 LIMIT 1`
		}

		entity = new(ArticleContent)

		if err := tx.Get(entity, query, articleId); err != nil {
			logbuch.Debug("Latest article content by organization id and article id and language id not found", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId, "language_id": langId})
			return nil
		}
	}

	entity.Authors = FindArticleContentAuthorUserByArticleIdTx(tx, articleId)
	return entity
}

func GetArticleContentLatestByArticleIdAndLanguageId(articleId, langId hide.ID, withContent bool) *ArticleContent {
	return GetArticleContentLatestByArticleIdAndLanguageIdTx(nil, articleId, langId, withContent)
}

func GetArticleContentLatestByArticleIdAndLanguageIdTx(tx *sqlx.Tx, articleId, langId hide.ID, withContent bool) *ArticleContent {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var query string

	if withContent {
		query = `SELECT * FROM "article_content"
			WHERE article_id = $1
			AND language_id = $2
			AND version = 0`
	} else {
		query = articleContentWithoutContentQuery + `WHERE article_id = $1
			AND language_id = $2
			AND version = 0`
	}

	entity := new(ArticleContent)

	if err := tx.Get(entity, query, articleId, langId); err != nil {
		// do not log, because this is getting out of hand when language is not available
		return nil
	}

	entity.Authors = FindArticleContentAuthorUserByArticleIdTx(tx, articleId)
	return entity
}

func GetArticleContentByArticleIdAndLanguageIdAndMaxVersion(articleId, langId hide.ID, version int) *ArticleContent {
	entity := new(ArticleContent)

	if err := connection.Get(entity, `SELECT * FROM "article_content"
		WHERE article_id = $1
		AND language_id = $2
		AND version <= $3
		AND version != 0
		ORDER BY version DESC
		LIMIT 1`, articleId, langId, version); err != nil {
		logbuch.Debug("Article content by organization id and article id and language id and max version not found", logbuch.Fields{"err": err, "article_id": articleId, "language_id": langId, "version": version})
		return nil
	}

	return entity
}

func GetArticleContentByArticleIdAndLanguageIdAndVersion(articleId, langId hide.ID, version int) *ArticleContent {
	entity := new(ArticleContent)

	if err := connection.Get(entity, `SELECT * FROM "article_content"
		WHERE article_id = $1
		AND language_id = $2
		AND version = $3
		AND version != 0
		ORDER BY version DESC
		LIMIT 1`, articleId, langId, version); err != nil {
		logbuch.Debug("Article content by organization id and article id and language id and version not found", logbuch.Fields{"err": err, "article_id": articleId, "language_id": langId, "version": version})
		return nil
	}

	return entity
}

func FindArticleContentByArticleIdAndLanguageId(articleId, langId hide.ID) []ArticleContent {
	var entities []ArticleContent

	if err := connection.Select(&entities, `SELECT * FROM "article_content" WHERE article_id = $1 AND language_id = $2 ORDER BY version ASC`, articleId, langId); err != nil {
		logbuch.Error("Error finding article content by article id and language id", logbuch.Fields{"err": err, "article_id": articleId, "language_id": langId})
		return nil
	}

	return entities
}

func FindArticleContentByArticleId(articleId hide.ID) []ArticleContent {
	return FindArticleContentByArticleIdTx(nil, articleId)
}

func FindArticleContentByArticleIdTx(tx *sqlx.Tx, articleId hide.ID) []ArticleContent {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var entities []ArticleContent

	if err := tx.Select(&entities, `SELECT * FROM "article_content" WHERE article_id = $1 ORDER BY version ASC`, articleId); err != nil {
		logbuch.Error("Error finding article content by article id", logbuch.Fields{"err": err, "article_id": articleId})
		return nil
	}

	return entities
}

func FindArticleContentIdByArticleIdAndWIPTx(tx *sqlx.Tx, articleId hide.ID) []hide.ID {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var ids []hide.ID

	if err := tx.Select(&ids, `SELECT id FROM "article_content" WHERE article_id = $1 AND wip IS TRUE`, articleId); err != nil {
		logbuch.Error("Error finding article content id by article id and wip", logbuch.Fields{"err": err, "article_id": articleId})
		return nil
	}

	return ids
}

func CountArticleContentVersionByArticleIdAndLanguageIdAndNotWIP(articleId, langId hide.ID) int {
	query := `SELECT COUNT(1) FROM "article_content"
		WHERE article_id = $1
		AND wip IS FALSE
		AND language_id = $2
		AND version != 0`
	var count int

	if err := connection.Get(&count, query, articleId, langId); err != nil {
		logbuch.Error("Error counting article content by article id and language id", logbuch.Fields{"err": err, "article_id": articleId, "language_id": langId})
		return 0
	}

	return count
}

func FindArticleContentVersionCommitByOrganizationIdAndArticleIdAndLanguageIdAndNotWIPLimit(orgaId, articleId, langId hide.ID, offset, n int) []ArticleContent {
	query := `SELECT "article_content".id,
		"article_content".version,
		"article_content".commit,
		"article_content".article_id,
		"article_content".language_id,
		"article_content".def_time,
		"article_content".mod_time,
		"user".id "user.id",
		"user".email "user.email",
		"user".firstname "user.firstname",
		"user".lastname "user.lastname",
		"user".language "user.language",
		"user".info "user.info",
		"user".picture "user.picture",
		"organization_member".id "user.organization_member.id",
		"organization_member".username "user.organization_member.username",
		"organization_member".info "user.organization_member.info",
		"organization_member".phone "user.organization_member.phone",
		"organization_member".mobile "user.organization_member.mobile",
		"organization_member".is_moderator "user.organization_member.is_moderator",
		"organization_member".is_admin "user.organization_member.is_admin"
		FROM "article_content"
		JOIN "article" ON "article_content".article_id = "article".id
		JOIN "user" ON "article_content".user_id = "user".id
		JOIN "organization_member" ON "user".id = "organization_member".user_id AND "organization_member".organization_id = $1
		WHERE article_id = $2
		AND "article_content".wip IS FALSE
		AND "article_content".language_id = $3
		AND version != 0
		ORDER BY version DESC
		LIMIT $4 OFFSET $5`
	var entities []ArticleContent

	if err := connection.Select(&entities, query, orgaId, articleId, langId, n, offset); err != nil {
		logbuch.Error("Error finding article content by organisation id and article id and language id with limit", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId, "language_id": langId, "offset": offset, "n": n})
		return nil
	}

	return entities
}

// Careful! This does require deleting the feed references beforehand.
func DeleteArticleContentById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := connection.Exec(tx, `DELETE FROM "article_content_author" WHERE article_content_id = $1`, id); err != nil {
		logbuch.Error("Error deleting article content author by article content id", logbuch.Fields{"err": err, "content_id": id})
		db.Rollback(tx)
		return err
	}

	if _, err := connection.Exec(tx, `DELETE FROM "article_content" WHERE id = $1`, id); err != nil {
		logbuch.Error("Error deleting article content by id", logbuch.Fields{"err": err, "content_id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

// Careful! This does require deleting the feed references beforehand.
func DeleteArticleContentByArticleIdAndWIP(tx *sqlx.Tx, articleId hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	if _, err := connection.Exec(tx, `DELETE FROM "article_content_author"
		WHERE article_content_id IN (SELECT id FROM "article_content" WHERE article_id = $1 AND wip IS TRUE)`, articleId); err != nil {
		logbuch.Error("Error deleting article content author by article id and wip", logbuch.Fields{"err": err, "article_id": articleId})
		db.Rollback(tx)
		return err
	}

	if _, err := connection.Exec(tx, `DELETE FROM "article_content" WHERE article_id = $1 AND wip IS TRUE`, articleId); err != nil {
		logbuch.Error("Error deleting article content by article id and wip", logbuch.Fields{"err": err, "article_id": articleId})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveArticleContent(tx *sqlx.Tx, entity *ArticleContent) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "article_content" (title,
			content,
			version,
			commit,
			wip,
			content_tsvector,
			title_tsvector,
			article_id,
			language_id,
			user_id,
			reading_time,
			schema_version,
			rtl)
			VALUES (:title,
			:content,
			:version,
			:commit,
			:wip,
			to_tsvector(:content_tsvector),
			to_tsvector(:title_tsvector),
			:article_id,
			:language_id,
			:user_id,
			:reading_time,
			:schema_version,
			:rtl)
			RETURNING id`,
		`UPDATE "article_content" SET title = :title,
			content = :content,
			version = :version,
			commit = :commit,
			wip = :wip,
			content_tsvector = to_tsvector(:content_tsvector),
			title_tsvector = to_tsvector(:title_tsvector),
			article_id = :article_id,
			language_id = :language_id,
			user_id = :user_id,
			reading_time = :reading_time,
			schema_version = :schema_version,
			rtl = :rtl
			WHERE id = :id`)
}
