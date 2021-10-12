package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func AddArticleListMember(organization *model.Organization, userId, listId hide.ID, userIds, memberGroupIds []hide.ID) ([]model.ArticleListMember, error) {
	if _, err := checkListExists(organization, listId); err != nil {
		return nil, err
	}

	if err := checkUserModAccess(listId, userId); err != nil {
		return nil, err
	}

	existingMemberIds := model.FindArticleListMemberUserIdByArticleListId(listId)
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to add article list members", logbuch.Fields{"err": err})
		return nil, errs.TxBegin
	}

	membersAdded, err := addMember(tx, organization, listId, userIds)

	if err != nil {
		return nil, err
	}

	var groupMembersAdded []model.ArticleListMember

	if organization.Expert {
		groupMembersAdded, err = addGroupsAsMembers(tx, organization, listId, memberGroupIds)

		if err != nil {
			return nil, err
		}
	}

	membersAdded = append(membersAdded, groupMembersAdded...)

	if err := createAddMemberFeed(tx, organization, userId, listId, existingMemberIds, membersAdded); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when adding article list members", logbuch.Fields{"err": err})
		return nil, errs.TxCommit
	}

	return membersAdded, nil
}

func addMember(tx *sqlx.Tx, organization *model.Organization, listId hide.ID, userIds []hide.ID) ([]model.ArticleListMember, error) {
	addedMembers := make([]model.ArticleListMember, 0, len(userIds))

	for _, addUserId := range userIds {
		// continue if member already
		if model.GetArticleListMemberByArticleListIdAndUserIdTx(tx, listId, addUserId) != nil {
			continue
		}

		// find user and create member
		user := model.GetUserWithOrganizationMemberByOrganizationIdAndIdTx(tx, organization.ID, addUserId)

		if user == nil {
			db.Rollback(tx)
			return nil, errs.UserNotFound
		}

		member := &model.ArticleListMember{ArticleListId: listId,
			UserId: addUserId,
			User:   user}

		if err := model.SaveArticleListMember(tx, member); err != nil {
			return nil, errs.Saving
		}

		// append to list for feed
		addedMembers = append(addedMembers, *member)
	}

	return addedMembers, nil
}

func addGroupsAsMembers(tx *sqlx.Tx, organization *model.Organization, listId hide.ID, groupIds []hide.ID) ([]model.ArticleListMember, error) {
	addedMembers := make([]model.ArticleListMember, 0, len(groupIds))

	for _, addGroupId := range groupIds {
		// continue if member already
		if model.GetArticleListMemberByArticleListIdAndUserGroupIdTx(tx, listId, addGroupId) != nil {
			continue
		}

		// find an create member
		group := model.GetUserGroupByOrganizationIdAndIdTx(tx, organization.ID, addGroupId)

		if group == nil {
			db.Rollback(tx)
			return nil, errs.GroupNotFound
		}

		member := &model.ArticleListMember{ArticleListId: listId,
			UserGroupId: addGroupId,
			UserGroup:   group}

		if err := model.SaveArticleListMember(tx, member); err != nil {
			return nil, errs.Saving
		}

		// append to list for feed
		addedMembers = append(addedMembers, *member)
	}

	return addedMembers, nil
}

func createAddMemberFeed(tx *sqlx.Tx, orga *model.Organization, userId, listId hide.ID, existingMemberIds []hide.ID, addedMembers []model.ArticleListMember) error {
	list := model.GetArticleListByOrganizationIdAndIdTx(tx, orga.ID, listId)

	if list == nil {
		db.Rollback(tx)
		return errs.ArticleListNotFound
	}

	refs := make([]interface{}, 0, len(addedMembers)+1)
	refs = append(refs, list)

	for _, member := range addedMembers {
		if member.UserId != 0 {
			refs = append(refs, member.User)
		} else {
			refs = append(refs, member.UserGroup)
		}
	}

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "add_article_list_member",
		Public:       list.Public,
		Access:       existingMemberIds,
		Notify:       util.RemoveIds(getMemberUserIds(tx, orga.ID, addedMembers), existingMemberIds),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when adding member to article list", logbuch.Fields{"err": err})
		return err
	}

	return nil
}

func getMemberUserIds(tx *sqlx.Tx, orgaId hide.ID, member []model.ArticleListMember) []hide.ID {
	ids := make([]hide.ID, 0, len(member))

	for _, m := range member {
		if m.UserId != 0 {
			ids = append(ids, m.UserId)
		} else {
			groupMember := model.FindUserGroupMemberUserByOrganizationIdAndUserGroupIdTx(tx, orgaId, m.UserGroupId)

			for i := range groupMember {
				ids = append(ids, groupMember[i].ID)
			}
		}
	}

	return ids
}
