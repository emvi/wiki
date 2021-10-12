package bookmark

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
)

// Toggles a bookmarked article or list.
func BookmarkObject(orga *model.Organization, userId, articleId, listId hide.ID) error {
	if articleId == 0 && listId == 0 {
		return errs.NoObjectToBookmark
	}

	if articleId != 0 && model.GetArticleByOrganizationIdAndIdIgnoreArchived(orga.ID, articleId) == nil {
		return errs.ArticleNotFound
	}

	if listId != 0 && model.GetArticleListByOrganizationIdAndId(orga.ID, listId) == nil {
		return errs.ArticleListNotFound
	}

	bookmark := model.GetBookmarkByOrganizationIdAndUserIdAndArticleIdOrArticleListid(orga.ID, userId, articleId, listId)

	if bookmark == nil {
		bookmark = &model.Bookmark{OrganizationId: orga.ID,
			UserId:        userId,
			ArticleId:     null.NewInt64(int64(articleId), articleId != 0),
			ArticleListId: null.NewInt64(int64(listId), listId != 0)}

		if err := model.SaveBookmark(nil, bookmark); err != nil {
			logbuch.Error("Error saving bookmark when toggling", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "article_id": articleId})
			return errs.Saving
		}
	} else {
		if err := model.DeleteBookmarkById(nil, bookmark.ID); err != nil {
			logbuch.Error("Error deleting bookmark when toggling")
			return errs.Saving
		}
	}

	return nil
}
