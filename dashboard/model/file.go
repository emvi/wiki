package model

import (
	"emviwiki/shared/db"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

type File struct {
	db.BaseEntity

	Filename         string `json:"filename"`
	OriginalFilename string `db:"original_filename" json:"original_filename"`
	MimeType         string `db:"mime_type" json:"mime_type"`
	Size             int64  `json:"size"`
	MD5              string `json:"md5"`
}

func GetFileById(id hide.ID) *File {
	query := `SELECT * FROM "file" WHERE id = $1`
	entity := new(File)

	if err := dashboardDB.Get(entity, query, id); err != nil {
		logbuch.Debug("File by id not found", logbuch.Fields{"err": err, "id": id})
		return nil
	}

	return entity
}

func GetFileByFilename(filename string) *File {
	query := `SELECT * FROM "file" WHERE filename = $1`
	entity := new(File)

	if err := dashboardDB.Get(entity, query, filename); err != nil {
		logbuch.Debug("File by filename not found", logbuch.Fields{"err": err, "filename": filename})
		return nil
	}

	return entity
}

func FindFile() []File {
	query := `SELECT * FROM "file" ORDER BY def_time DESC`
	var entities []File

	if err := dashboardDB.Select(&entities, query); err != nil {
		logbuch.Error("Error finding file", logbuch.Fields{"err": err})
		return nil
	}

	return entities
}

func SaveFile(tx *sqlx.Tx, entity *File) error {
	return dashboardDB.SaveEntity(tx, entity,
		`INSERT INTO "file" (filename,
			original_filename,
			mime_type,
			size,
			md5) VALUES (:filename,
			:original_filename,
			:mime_type,
			:size,
			:md5) RETURNING id`,
		`UPDATE "file" SET filename = :filename,
			original_filename = :original_filename,
			mime_type = :mime_type,
			size = :size,
			md5 = :md5
			WHERE id = :id`)
}

func DeleteFileById(tx *sqlx.Tx, id hide.ID) error {
	if tx == nil {
		var err error
		tx, err = dashboardDB.Beginx()

		if err != nil {
			return err
		}

		defer func() {
			if err := tx.Commit(); err != nil {
				logbuch.Error("Error committing transaction to delete newsletter", logbuch.Fields{"err": err, "id": id})
			}
		}()
	}

	if _, err := tx.Exec(`DELETE FROM "file" WHERE id = $1`, id); err != nil {
		logbuch.Error("Error deleting file by id", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
