package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"time"
)

type File struct {
	db.BaseEntity

	OrganizationId hide.ID     `db:"organization_id" json:"-"` // nullable
	UserId         hide.ID     `db:"user_id" json:"-"`         // uploader
	ArticleId      hide.ID     `db:"article_id" json:"-"`      // optional article this file is used in
	RoomId         null.String `db:"room_id" json:"-"`         // optional article without ID this file is used in
	LanguageId     hide.ID     `db:"language_id" json:"-"`     // optional language of content this file is used in
	OriginalName   string      `db:"original_name" json:"original_name"`
	UniqueName     string      `db:"unique_name" json:"unique_name"`
	Path           string      `json:"-"`
	Type           string      `json:"type"`
	MimeType       string      `db:"mime_type" json:"mime_type"`
	Size           int64       `json:"size"`
	MD5            string      `json:"md5"`
}

func GetFileStorageUsageByOrganizationId(orgaId hide.ID) int64 {
	query := `SELECT CASE WHEN SUM("size") IS NULL THEN 0 ELSE SUM("size") END
		FROM (SELECT DISTINCT ON (unique_name) "size" FROM "file" WHERE organization_id = $1 AND article_id IS NOT NULL) AS files`
	var size int64

	if err := connection.Get(&size, query, orgaId); err != nil {
		logbuch.Debug("File storage usage by organization id not found", logbuch.Fields{"err": err, "orga_id": orgaId})
		return -1
	}

	return size
}

func GetFileByOrganizationIdAndUniqueName(orgaId hide.ID, name string) *File {
	entity := new(File)

	if err := connection.Get(entity, `SELECT * FROM "file" WHERE organization_id = $1 AND unique_name = $2`, orgaId, name); err != nil {
		logbuch.Debug("File by organization id and unique name not found", logbuch.Fields{"err": err, "orga_id": orgaId, "name": name})
		return nil
	}

	return entity
}

func GetFileByUniqueName(name string) *File {
	entity := new(File)

	if err := connection.Get(entity, `SELECT * FROM "file" WHERE unique_name = $1`, name); err != nil {
		logbuch.Debug("File by unique name not found", logbuch.Fields{"err": err, "name": name})
		return nil
	}

	return entity
}

func GetFileByOrganizationIdAndMD5AndArticleIdOrRoomIdNotNullTx(tx *sqlx.Tx, orgaId hide.ID, md5 string) *File {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT * FROM "file"
		WHERE organization_id = $1
		AND md5 = $2
		AND (article_id IS NOT NULL OR room_id IS NOT NULL)`
	entity := new(File)

	if err := tx.Get(entity, query, orgaId, md5); err != nil {
		logbuch.Debug("File by organization id and md5 and article id not null not found", logbuch.Fields{"err": err, "orga_Id": orgaId, "md5": md5})
		return nil
	}

	return entity
}

func FindFileByOrganizationId(orgaId hide.ID) []File {
	query := `SELECT * FROM "file" WHERE organization_id = $1 AND (article_id IS NOT NULL OR room_id IS NOT NULL)`
	var entities []File

	if err := connection.Select(&entities, query, orgaId); err != nil {
		logbuch.Error("Error reading files by organization id", logbuch.Fields{"err": err, "orga_id": orgaId})
		return nil
	}

	return entities
}

func FindFileByOrganizationIdAndUniqueNameAndNotId(orgaId hide.ID, uniqueName string, id hide.ID) []File {
	query := `SELECT * FROM "file" WHERE organization_id = $1 AND unique_name = $2 AND id != $3`
	var entities []File

	if err := connection.Select(&entities, query, orgaId, uniqueName, id); err != nil {
		logbuch.Error("Error reading files by organization id unique name and not id", logbuch.Fields{"err": err, "orga_id": orgaId, "unique_name": uniqueName, "id": id})
		return nil
	}

	return entities
}

func FindFileByOrganizationIdAndArticleIdAndUniqueInOrganization(tx *sqlx.Tx, orgaId, articleId hide.ID) []File {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT a.* FROM "file" a
		WHERE a.organization_id = $1
		AND a.article_id = $2
		AND NOT EXISTS (
			SELECT 1 FROM "file" b
			WHERE b.organization_id = a.organization_id
			AND b.md5 = a.md5
			AND b.article_id != a.article_id
		)`
	var entities []File

	if err := tx.Select(&entities, query, orgaId, articleId); err != nil {
		logbuch.Error("Error reading files by organization id and article id", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId})
		return nil
	}

	return entities
}

func FindFileByOrganizationIdAndRoomIdTx(tx *sqlx.Tx, orgaId hide.ID, roomId string) []File {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	query := `SELECT * FROM "file" WHERE organization_id = $1 AND room_id = $2`
	var entities []File

	if err := tx.Select(&entities, query, orgaId, roomId); err != nil {
		logbuch.Error("Error reading files by organization id and room id", logbuch.Fields{"err": err, "orga_id": orgaId, "room_id": roomId})
		return nil
	}

	return entities
}

func FindFileByOrganizationIdAndArticleIdAndLanguageIdAndDefTimeAfter(orgaId, articleId, langId hide.ID, defTime time.Time) []File {
	query := `SELECT * FROM "file" WHERE organization_id = $1 AND article_id = $2 AND language_id = $3 AND def_time > $4`
	var entities []File

	if err := connection.Select(&entities, query, orgaId, articleId, langId, defTime); err != nil {
		logbuch.Error("Error reading files by organization id and article id and language id and after def time", logbuch.Fields{"err": err, "orga_id": orgaId, "article_id": articleId, "lang_id": langId, "def_time": defTime})
		return nil
	}

	return entities
}

func CountFileByUniqueNameTx(tx *sqlx.Tx, name string) int {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	var count int

	if err := tx.Get(&count, `SELECT COUNT(1) FROM "file" WHERE unique_name = $1`, name); err != nil {
		logbuch.Error("Error counting file by unique name", logbuch.Fields{"err": err, "name": name})
		return 0
	}

	return count
}

func DeleteFileById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		tx, _ = connection.Beginx()
		defer db.Commit(tx)
	}

	_, err := tx.Exec(`DELETE FROM "file" WHERE id = $1`, id)

	if err != nil {
		logbuch.Error("Error deleting file by id", logbuch.Fields{"err": err, "id": id})
		db.Rollback(tx)
		return err
	}

	return nil
}

func SaveFile(tx *sqlx.Tx, entity *File) error {
	return connection.SaveEntity(tx, entity,
		`INSERT INTO "file" (organization_id,
			user_id,
			article_id,
			room_id,
			language_id,
			original_name,
			unique_name,
			path,
			type,
			mime_type,
			size,
			md5)
			VALUES (:organization_id,
			:user_id,
			:article_id,
			:room_id,
			:language_id,
			:original_name,
			:unique_name,
			:path,
			:type,
			:mime_type,
			:size,
			:md5) RETURNING id`,
		`UPDATE "file" SET organization_id = :organization_id,
			user_id = :user_id,
			article_id = :article_id,
			room_id = :room_id,
			language_id = :language_id,
			original_name = :original_name,
			unique_name = :unique_name,
			path = :path,
			type = :type,
			mime_type = :mime_type,
			size = :size,
			md5 = :md5
			WHERE id = :id`)
}
