package bookmark

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func IsBookmarked(userId, articleId, listId hide.ID) bool {
	if articleId != 0 {
		return model.GetBookmarkByUserIdAndArticleId(userId, articleId) != nil
	} else if listId != 0 {
		return model.GetBookmarkByUserIdAndArticleListId(userId, listId) != nil
	}

	return false
}
