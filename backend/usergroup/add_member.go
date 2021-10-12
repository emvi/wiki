package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func AddUserGroupMember(organization *model.Organization, userId, groupId hide.ID, userIds, groupIds []hide.ID) ([]model.UserGroupMember, error) {
	_, err := checkGroupExists(organization, groupId)

	if err != nil {
		return nil, err
	}

	if err := checkUserAccess(organization, userId, groupId); err != nil {
		return nil, err
	}

	userIds, err = AppendUserFromGroups(nil, organization, groupId, userIds, groupIds)

	if err != nil {
		return nil, err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to add group members", logbuch.Fields{"err": err})
		return nil, errs.TxBegin
	}

	usersAdded := make([]model.User, 0)
	membersAdded := make([]model.UserGroupMember, 0)

	for _, addUserId := range userIds {
		userEntity := model.GetUserWithOrganizationMemberByOrganizationIdAndIdTx(tx, organization.ID, addUserId)

		if userEntity == nil {
			db.Rollback(tx)
			return nil, errs.UserNotFound
		}

		if model.GetUserGroupMemberByGroupIdAndUserIdTx(tx, groupId, addUserId) != nil {
			continue
		}

		member := &model.UserGroupMember{UserGroupId: groupId,
			UserId: addUserId,
			User:   userEntity}

		if err := model.SaveUserGroupMember(tx, member); err != nil {
			return nil, errs.Saving
		}

		usersAdded = append(usersAdded, *userEntity)
		membersAdded = append(membersAdded, *member)
	}

	if err := createAddMemberFeed(tx, organization, userId, groupId, usersAdded); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when adding group members", logbuch.Fields{"err": err})
		return nil, errs.TxCommit
	}

	return membersAdded, nil
}

// AppendUserFromGroups reads all member user IDs for given group IDs and appends the user IDs to the given list of IDs.
// This does NOT remove duplicates!
// Returns an error if one of the groups cannot be found.
func AppendUserFromGroups(tx *sqlx.Tx, organization *model.Organization, groupId hide.ID, userIds, groupIds []hide.ID) ([]hide.ID, error) {
	for _, addGroupId := range groupIds {
		// adding members from group A to group A is allowed, but there is nothing to do
		if groupId == addGroupId {
			continue
		}

		// find group to add
		groupToAdd := model.GetUserGroupByOrganizationIdAndIdTx(tx, organization.ID, addGroupId)

		if groupToAdd == nil {
			return nil, errs.GroupToAddNotFound
		}

		// add all members to list of user IDs
		member := model.FindUserGroupMemberOnlyByUserGroupIdTx(tx, addGroupId)

		for _, m := range member {
			userIds = append(userIds, m.UserId)
		}
	}

	return userIds, nil
}

func createAddMemberFeed(tx *sqlx.Tx, orga *model.Organization, userId, groupId hide.ID, usersAdded []model.User) error {
	group := model.GetUserGroupByOrganizationIdAndIdTx(tx, orga.ID, groupId)

	if group == nil {
		db.Rollback(tx)
		return errs.GroupNotFound
	}

	userAddedIds := make([]hide.ID, 0)
	refs := make([]interface{}, 0)
	refs = append(refs, group)

	for i := range usersAdded {
		refs = append(refs, &usersAdded[i])
		userAddedIds = append(userAddedIds, usersAdded[i].ID)
	}

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "add_user_group_member",
		Public:       true,
		Notify:       userAddedIds,
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when adding member to user group", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
