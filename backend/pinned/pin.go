package pinned

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

// Pins an article or list.
func PinObject(orga *model.Organization, userId, articleId, listId hide.ID) error {
	if articleId == 0 && listId == 0 {
		return errs.NoObjectToPin
	}

	if _, err := perm.CheckUserIsAdminOrMod(orga.ID, userId); err != nil {
		return errs.PermissionDenied
	}

	if articleId != 0 {
		article := model.GetArticleByOrganizationIdAndId(orga.ID, articleId)

		if article == nil {
			return errs.ArticleNotFound
		}

		article.Pinned = !article.Pinned

		if err := model.SaveArticle(nil, article); err != nil {
			logbuch.Error("Error saving article when pinning", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId})
			return errs.Saving
		}
	} else {
		list := model.GetArticleListByOrganizationIdAndId(orga.ID, listId)

		if list == nil {
			return errs.ArticleListNotFound
		}

		list.Pinned = !list.Pinned

		if err := model.SaveArticleList(nil, list); err != nil {
			logbuch.Error("Error saving article list when pinning", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "list_id": listId})
			return errs.Saving
		}
	}

	return nil
}
