package content

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"path/filepath"
)

// DeleteFile checks user permissions to delete a file and deletes it if allowed.
// This does not check the file is used anywhere else.
func DeleteFile(orga *model.Organization, userId hide.ID, uniqueName string) error {
	var orgaId hide.ID

	if orga != nil {
		orgaId = orga.ID
	}

	file := model.GetFileByUniqueName(uniqueName)

	if file == nil || (orga == nil && file.UserId != userId) {
		return errs.PermissionDenied
	}

	if err := model.DeleteFileById(nil, file.ID); err != nil {
		logbuch.Error("Error deleting file", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "path": file.Path})
	}

	go DeleteFileInStore(orgaId, userId, filepath.Join(file.Path, file.UniqueName))

	return nil
}

// DeleteFileForArticle deletes all files for given article.
// This function must not be triggered by users directly, since it does not check permissions.
// The file on disk is only removed if it is not referenced by another article.
// The caller is notified through the (optional) channel when all files are deleted.
func DeleteFileForArticle(tx *sqlx.Tx, orga *model.Organization, userId, articleId hide.ID, done chan bool) {
	files := model.FindFileByOrganizationIdAndArticleIdAndUniqueInOrganization(tx, orga.ID, articleId)

	go func() {
		for _, file := range files {
			DeleteFileInStore(orga.ID, userId, filepath.Join(file.Path, file.UniqueName))
		}

		done <- true
	}()
}

// DeleteFileInStore deletes given file in store by path.
// This is just a fancy wrapper for store.Delete and logs the error should it occur.
func DeleteFileInStore(orgaId, userId hide.ID, path string) {
	if err := store.Delete(path); err != nil {
		logbuch.Error("Error deleting file in store", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId, "path": path})
	}
}
