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

func RemoveUserGroupMember(organization *model.Organization, userId, groupId hide.ID, removeMemberIds []hide.ID) error {
	_, err := checkGroupExists(organization, groupId)

	if err != nil {
		return err
	}

	if err := checkUserAccess(organization, userId, groupId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to remove group members", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	removedUsers, err := removeMember(tx, userId, groupId, removeMemberIds)

	if err != nil {
		return err
	}

	if err := createRemoveMemberFeed(tx, organization, userId, groupId, removedUsers); err != nil {
		logbuch.Error("Error creating feed when removing member from user group", logbuch.Fields{"err": err})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when removing group members", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func removeMember(tx *sqlx.Tx, userId, groupId hide.ID, removeMemberIds []hide.ID) ([]model.User, error) {
	removedUsers := make([]model.User, 0, len(removeMemberIds))

	for _, removeMemberId := range removeMemberIds {
		if _, err := checkGroupMemberExists(tx, groupId, removeMemberId); err != nil {
			continue
		}

		removeMember := model.GetUserGroupMemberByIdTx(tx, removeMemberId)

		if removeMember == nil {
			db.Rollback(tx)
			return nil, errs.MemberNotFound
		}

		if userId == removeMember.UserId {
			db.Rollback(tx)
			return nil, errs.RemoveYourself
		}

		if err := model.DeleteUserGroupMemberByUserGroupIdAndId(tx, groupId, removeMemberId); err != nil {
			db.Rollback(tx)
			return nil, errs.Saving
		}

		user := model.GetUserByIdTx(tx, removeMember.UserId)

		if user == nil {
			db.Rollback(tx)
			return nil, errs.UserNotFound
		}

		removedUsers = append(removedUsers, *user)
	}

	return removedUsers, nil
}

func checkGroupMemberExists(tx *sqlx.Tx, groupId, memberId hide.ID) (*model.UserGroupMember, error) {
	member := model.GetUserGroupMemberByGroupIdAndIdTx(tx, groupId, memberId)

	if member == nil {
		return nil, errs.MemberNotFound
	}

	return member, nil
}

func createRemoveMemberFeed(tx *sqlx.Tx, orga *model.Organization, userId, groupId hide.ID, removedUsers []model.User) error {
	group := model.GetUserGroupByOrganizationIdAndIdTx(tx, orga.ID, groupId)

	if group == nil {
		db.Rollback(tx)
		return errs.GroupNotFound
	}

	removedUserIds := make([]hide.ID, len(removedUsers))
	refs := make([]interface{}, len(removedUsers)+1)
	refs[0] = group

	for i, user := range removedUsers {
		removedUserIds[i] = user.ID
		refs[i+1] = &user
	}

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "remove_user_group_member",
		Public:       true,
		Notify:       removedUserIds,
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
