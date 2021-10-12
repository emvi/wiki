package member

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func ToggleModerator(orga *model.Organization, userId, memberUserId hide.ID) error {
	if orga.OwnerUserId == memberUserId {
		return errs.OrganizationOwner
	}

	_, member, err := checkIsAdminAndGetMember(orga.ID, userId, memberUserId)

	if err != nil {
		return err
	}

	if userId == member.UserId {
		return errs.ModeratorYourself
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when toggling moderator", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member.IsModerator = !member.IsModerator
	member.ReadOnly = false

	if !member.IsModerator {
		member.IsAdmin = false
	}

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		logbuch.Error("Error toggling organization member moderator", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return errs.Saving
	}

	if err := updateDefaultGroupMembers(tx, orga.ID, memberUserId, member); err != nil {
		logbuch.Error("Error updating default group when toggling organization member moderator", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return err
	}

	if err := createToggleModeratorFeed(tx, orga, userId, member); err != nil {
		logbuch.Error("Error creating feed when toggling organization member moderator", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when toggling moderator", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func createToggleModeratorFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, member *model.OrganizationMember) error {
	reason := "set_organization_moderator"

	if !member.IsModerator {
		reason = "remove_organization_moderator"
	}

	memberUser := model.GetUserByOrganizationIdAndIdTx(tx, orga.ID, member.UserId)

	if memberUser == nil {
		db.Rollback(tx)
		return errs.UserNotFound
	}

	refs := make([]interface{}, 1)
	refs[0] = memberUser
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       false,
		Notify:       []hide.ID{memberUser.ID},
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
