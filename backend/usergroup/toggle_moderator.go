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

func ToggleUserGroupModerator(organization *model.Organization, userId, groupId, memberId hide.ID) error {
	member, err := checkGroupMemberExists(nil, groupId, memberId)

	if err != nil {
		return err
	}

	if userId == member.UserId {
		return errs.ModeratorYourself
	}

	if _, err := checkGroupExists(organization, member.UserGroupId); err != nil {
		return err
	}

	if err := checkUserAccess(organization, userId, member.UserGroupId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when toggling user group moderator", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member.IsModerator = !member.IsModerator

	if err := model.SaveUserGroupMember(tx, member); err != nil {
		return errs.Saving
	}

	if err := createToggleModeratorFeed(tx, organization, userId, member); err != nil {
		logbuch.Error("Error creating feed when toggling user group moderator", logbuch.Fields{"err": err})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when toggling user group moderator", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func createToggleModeratorFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, member *model.UserGroupMember) error {
	reason := "set_user_group_moderator"

	if !member.IsModerator {
		reason = "remove_user_group_moderator"
	}

	group := model.GetUserGroupByOrganizationIdAndIdTx(tx, orga.ID, member.UserGroupId)

	if group == nil {
		db.Rollback(tx)
		return errs.GroupNotFound
	}

	memberUser := model.GetUserByOrganizationIdAndIdTx(tx, orga.ID, member.UserId)

	if memberUser == nil {
		db.Rollback(tx)
		return errs.UserNotFound
	}

	refs := make([]interface{}, 2)
	refs[0] = group
	refs[1] = memberUser
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       false,
		Notify:       []hide.ID{member.UserId},
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
