package observe

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

// Toggles an observed object.
func ObserveObject(organization *model.Organization, userId, articleId, articleListId, userGroupId hide.ID) error {
	if articleId == 0 && articleListId == 0 && userGroupId == 0 {
		return errs.NoObjectToObserve
	}

	if articleId != 0 && model.GetArticleByOrganizationIdAndId(organization.ID, articleId) == nil {
		return errs.ArticleNotFound
	}

	if articleListId != 0 && model.GetArticleListByOrganizationIdAndId(organization.ID, articleListId) == nil {
		return errs.ArticleListNotFound
	}

	if userGroupId != 0 && model.GetUserGroupByOrganizationIdAndId(organization.ID, userGroupId) == nil {
		return errs.GroupNotFound
	}

	observe := model.GetObservedObjectByUserIdAndArticleIdOrArticleListIdOrUserGroupId(userId, articleId, articleListId, userGroupId)

	if observe == nil {
		observe = &model.ObservedObject{UserId: userId,
			ArticleId:     articleId,
			ArticleListId: articleListId,
			UserGroupId:   userGroupId}

		if err := model.SaveObservedObject(nil, observe); err != nil {
			return errs.Saving
		}
	} else {
		if err := model.DeleteObservedObjectById(nil, observe.ID); err != nil {
			return errs.Saving
		}
	}

	return nil
}
