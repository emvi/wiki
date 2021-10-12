package observe

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func IsObserved(userId, articleId, articleListId, userGroupId hide.ID) bool {
	return (articleId != 0 && model.GetObservedObjectByUserIdAndArticleId(userId, articleId) != nil) ||
		(articleListId != 0 && model.GetObservedObjectByUserIdAndArticleListId(userId, articleListId) != nil) ||
		(userGroupId != 0 && model.GetObservedObjectByUserIdAndUserGroupId(userId, userGroupId) != nil)
}
