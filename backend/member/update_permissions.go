package member

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/usergroup"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func updateObjectPermissionsForInactiveUser(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, removePermissions bool) error {
	adminGroup := usergroup.GetAdminGroup(tx, orga.ID)
	modGroup := usergroup.GetModGroup(tx, orga.ID)
	userIds, err := getAdminModUserIds(tx, orga, userId, adminGroup, modGroup)

	if err != nil {
		db.Rollback(tx)
		return err
	}

	// update object permissions
	// the transaction is rolled back inside the model package functions
	articles, err := updateArticlePermissionsForInactiveUser(tx, orga, userId, adminGroup, modGroup, removePermissions)

	if err != nil {
		return err
	}

	lists, err := updateListPermissionsForInactiveUser(tx, orga, userId, adminGroup, modGroup, removePermissions)

	if err != nil {
		return err
	}

	groups, err := updateGroupPermissionsForInactiveUser(tx, orga, userId, userIds, removePermissions)

	if err != nil {
		return err
	}

	// create feed for all admins/mods
	if len(articles)+len(lists)+len(groups) > 0 {
		if err := createAdminModAccessPermissionFeed(tx, orga, userId, userIds, articles, lists, groups); err != nil {
			return err
		}
	}

	return nil
}

func getAdminModUserIds(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, adminGroup, modGroup *model.UserGroup) ([]hide.ID, error) {
	userIds := make([]hide.ID, 0)
	var err error
	userIds, err = usergroup.AppendUserFromGroups(tx, orga, 0, userIds, []hide.ID{adminGroup.ID, modGroup.ID})

	if err != nil {
		return nil, err
	}

	// remove duplicates and user ID so that admins/mods don't get added to an article/list/group at the same time they're removed
	userIds = util.RemoveDuplicateIds(userIds)
	return util.RemoveId(userIds, userId), nil
}

func updateArticlePermissionsForInactiveUser(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, adminGroup, modGroup *model.UserGroup, removePermissions bool) ([]model.Article, error) {
	articles := model.FindArticleByOrganizationIdAndAccessByUserIdOnlyTx(tx, orga.ID, userId)
	transferredArticles := make([]model.Article, 0, len(articles))

	for _, article := range articles {
		// exclude private articles
		if article.Private {
			continue
		}

		if removePermissions {
			if err := model.DeleteArticleAccessByArticleId(tx, article.ID); err != nil {
				logbuch.Error("Error deleting article access when updating user object permissions", logbuch.Fields{"err": err, "article_id": article.ID})
				return nil, errs.Saving
			}
		}

		adminMember := &model.ArticleAccess{ArticleId: article.ID,
			UserGroupId: adminGroup.ID,
			Write:       true}

		if err := model.SaveArticleAccess(tx, adminMember); err != nil {
			logbuch.Error("Error saving admin group article access when updating user object permissions", logbuch.Fields{"err": err, "article_id": article.ID})
			return nil, errs.Saving
		}

		modMember := &model.ArticleAccess{ArticleId: article.ID,
			UserGroupId: modGroup.ID,
			Write:       true}

		if err := model.SaveArticleAccess(tx, modMember); err != nil {
			logbuch.Error("Error saving moderator group article access when updating user object permissions", logbuch.Fields{"err": err, "article_id": article.ID})
			return nil, errs.Saving
		}

		transferredArticles = append(transferredArticles, article)
	}

	return transferredArticles, nil
}

func updateListPermissionsForInactiveUser(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, adminGroup, modGroup *model.UserGroup, removePermissions bool) ([]model.ArticleList, error) {
	lists := model.FindArticleListsByOrganizationIdAndAccessByUserIdOnlyTx(tx, orga.ID, userId)
	transferredLists := make([]model.ArticleList, 0, len(lists))

	for _, list := range lists {
		if removePermissions {
			if err := model.DeleteArticleListMemberByArticleListId(tx, list.ID); err != nil {
				logbuch.Error("Error deleting list member when updating user object permissions", logbuch.Fields{"err": err, "list_id": list.ID})
				return nil, errs.Saving
			}
		}

		adminMember := &model.ArticleListMember{ArticleListId: list.ID,
			UserGroupId: adminGroup.ID,
			IsModerator: true}

		if err := model.SaveArticleListMember(tx, adminMember); err != nil {
			logbuch.Error("Error saving admin group list member when updating user object permissions", logbuch.Fields{"err": err, "list_id": list.ID})
			return nil, errs.Saving
		}

		modMember := &model.ArticleListMember{ArticleListId: list.ID,
			UserGroupId: modGroup.ID,
			IsModerator: true}

		if err := model.SaveArticleListMember(tx, modMember); err != nil {
			logbuch.Error("Error saving moderator group list member when updating user object permissions", logbuch.Fields{"err": err, "list_id": list.ID})
			return nil, errs.Saving
		}

		transferredLists = append(transferredLists, list)
	}

	return transferredLists, nil
}

func updateGroupPermissionsForInactiveUser(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, newMemberIds []hide.ID, removePermissions bool) ([]model.UserGroup, error) {
	groups := model.FindUserGroupsByOrganizationIdAndAccessByUserIdOnlyTx(tx, orga.ID, userId)
	transferredGroups := make([]model.UserGroup, 0, len(groups))

	for _, group := range groups {
		if removePermissions {
			if err := model.DeleteUserGroupMemberByUserGroupId(tx, group.ID); err != nil {
				logbuch.Error("Error deleting group member when updating user object permissions", logbuch.Fields{"err": err, "group_id": group.ID})
				return nil, errs.Saving
			}
		}

		for _, id := range newMemberIds {
			member := &model.UserGroupMember{UserGroupId: group.ID,
				UserId:      id,
				IsModerator: true}

			if err := model.SaveUserGroupMember(tx, member); err != nil {
				logbuch.Error("Error saving group member when updating user object permissions", logbuch.Fields{"err": err, "group_id": group.ID, "member_user_id": id})
				return nil, errs.Saving
			}
		}

		transferredGroups = append(transferredGroups, group)
	}

	return transferredGroups, nil
}

func createAdminModAccessPermissionFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, adminModIds []hide.ID, articles []model.Article, lists []model.ArticleList, groups []model.UserGroup) error {
	refs := make([]interface{}, 0, len(articles)+len(lists)+len(groups))

	for _, article := range articles {
		refs = append(refs, &article)
	}

	for _, list := range lists {
		refs = append(refs, &list)
	}

	for _, group := range groups {
		refs = append(refs, &group)
	}

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "transfer_ownership",
		Public:       false,
		Notify:       adminModIds,
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when transfering ownership of objects to admin/mods", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
