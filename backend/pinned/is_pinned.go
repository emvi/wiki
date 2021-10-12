package pinned

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func IsPinned(orgaId, articleId, listId hide.ID) bool {
	if articleId != 0 {
		return model.GetArticleByOrganizationIdAndIdAndPinned(orgaId, articleId) != nil
	} else if listId != 0 {
		return model.GetArticleListIdByOrganizationIdAndUserIdAndIdAndPinned(orgaId, listId) != nil
	}

	return false
}
