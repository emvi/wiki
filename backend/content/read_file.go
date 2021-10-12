package content

import (
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"io"
	"path/filepath"

	"emviwiki/backend/errs"
	"emviwiki/shared/model"
)

// TODO protect files with signed URLs
func ReadFile(uniqueName string) (*model.File, io.ReadCloser, error) {
	file := model.GetFileByUniqueName(uniqueName)

	if file == nil {
		// only log in debug, because this might happen very frequently
		logbuch.Debug("File not found in database by unique name", logbuch.Fields{"unique_name": uniqueName})
		return nil, nil, errs.FileNotFound
	}

	reader, err := store.Read(filepath.Join(file.Path, file.UniqueName))

	if err != nil {
		// only log in debug, because this might happen very frequently
		logbuch.Debug("Error reading file from store", logbuch.Fields{"err": err})
		return nil, nil, errs.FileNotFound
	}

	return file, reader, nil
}

func checkArticleAccess(orgaId, userId, articleId hide.ID, write bool) error {
	article := model.GetArticleByOrganizationIdAndIdIgnoreArchived(orgaId, articleId)

	if article == nil {
		return errs.ArticleNotFound
	}

	access := model.FindArticleAccessByArticleIdAndUserIdAndWrite(articleId, userId, write)

	if (!write && !article.ReadEveryone || write && !article.WriteEveryone) && len(access) == 0 {
		return errs.PermissionDenied
	}

	return nil
}
