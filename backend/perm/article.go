package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

// SaveArticleAccess is used to set an access permission for an user or group.
// This permission can either be read or write depending on if Write is set to true or not.
type SaveArticleAccess struct {
	UserId      hide.ID `json:"user_id"`
	UserGroupId hide.ID `json:"user_group_id"`
	Write       bool    `json:"write"`
}

// CheckUserReadOrWriteAccess checks a user has read or write access to an article.
func CheckUserReadOrWriteAccess(articleId, userId hide.ID) bool {
	return len(model.FindArticleAccessByArticleIdAndUserId(articleId, userId)) > 0
}

// CheckUserWriteAccess checks a user has write access to an article.
func CheckUserWriteAccess(articleId, userId hide.ID) bool {
	return len(model.FindArticleAccessByArticleIdAndUserIdAndWrite(articleId, userId, true)) > 0
}

// FilterAccessList filters all duplicate entries from an access list and merges write access
// (write access is set when one entry in the list permits it).
func FilterAccessList(access []SaveArticleAccess) []SaveArticleAccess {
	list := make([]SaveArticleAccess, 0)

	for _, a := range access {
		index := findAccessInList(list, a)

		if index == -1 {
			list = append(list, a)
		} else if a.Write {
			// only set to true, to not overwrite with false
			list[index].Write = true
		}
	}

	return list
}

func findAccessInList(list []SaveArticleAccess, access SaveArticleAccess) int {
	for i, a := range list {
		if (a.UserId != 0 && a.UserId == access.UserId) ||
			(a.UserGroupId != 0 && a.UserGroupId == access.UserGroupId) {
			return i
		}
	}

	return -1
}

// CheckAccessList validates all users and groups in list exist.
func CheckAccessList(tx *sqlx.Tx, organization *model.Organization, access []SaveArticleAccess) error {
	for _, a := range access {
		if !checkUserOrGroupExists(tx, organization.ID, a) {
			logbuch.Warn("User or user group not found", logbuch.Fields{"user_id": a.UserId, "user_group_id": a.UserGroupId})
			return errs.UserOrUserGroupNotFound
		}
	}

	return nil
}

func checkUserOrGroupExists(tx *sqlx.Tx, orgaId hide.ID, access SaveArticleAccess) bool {
	if (access.UserId != 0 && model.GetUserByOrganizationIdAndIdTx(tx, orgaId, access.UserId) != nil) ||
		(access.UserGroupId != 0 && model.GetUserGroupByOrganizationIdAndIdTx(tx, orgaId, access.UserGroupId) != nil) {
		return true
	}

	return false
}

// SaveAccessList saves all permissions and sets references to article.
// Skips groups if organization is not expert.
// Performs a rollback on error.
func SaveAccessList(tx *sqlx.Tx, orga *model.Organization, access []SaveArticleAccess, articleId hide.ID) error {
	for _, a := range access {
		if a.UserGroupId != 0 && !orga.Expert {
			continue
		}

		articleAccess := &model.ArticleAccess{UserId: a.UserId,
			UserGroupId: a.UserGroupId,
			ArticleId:   articleId,
			Write:       a.Write}

		if err := model.SaveArticleAccess(tx, articleAccess); err != nil {
			logbuch.Error("Error saving article access when saving article", logbuch.Fields{"err": err, "user_id": a.UserId, "user_group_id": a.UserGroupId})
			return errs.Saving
		}
	}

	return nil
}

// DeleteAccess deletes all permissions to an article. Performs a rollback on error.
func DeleteAccess(tx *sqlx.Tx, articleId hide.ID) error {
	if err := model.DeleteArticleAccessByArticleId(tx, articleId); err != nil {
		logbuch.Error("Error deleting article access for existing article", logbuch.Fields{"err": err, "article_id": articleId})
		return errs.Saving
	}

	return nil
}

// GetUserIdsFromArticleAccess returns all user IDs for given slice of access objects.
func GetUserIdsFromArticleAccess(access []SaveArticleAccess) []hide.ID {
	ids := make([]hide.ID, 0)

	for _, a := range access {
		if a.UserId != 0 {
			ids = append(ids, a.UserId)
		} else {
			ids = append(ids, model.FindUserGroupMemberUserIdByUserGroupId(a.UserGroupId)...)
		}
	}

	return ids
}

// GetUserIdsFromAccess returns all user IDs for given slice of access objects.
func GetUserIdsFromAccess(tx *sqlx.Tx, access []model.ArticleAccess) []hide.ID {
	ids := make([]hide.ID, 0)

	for _, a := range access {
		if a.UserId != 0 {
			ids = append(ids, a.UserId)
		} else {
			ids = append(ids, model.FindUserGroupMemberUserIdByUserGroupIdTx(tx, a.UserGroupId)...)
		}
	}

	return ids
}
