package member

import (
	"emviwiki/backend/billing"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func ToggleReadOnly(orga *model.Organization, userId, memberUserId hide.ID) error {
	if orga.OwnerUserId == memberUserId {
		return errs.OrganizationOwner
	}

	_, member, err := checkIsAdminAndGetMember(orga.ID, userId, memberUserId)

	if err != nil {
		return err
	}

	if userId == member.UserId {
		return errs.ChangeReadOnlyYourself
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when toggling read only member", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member.ReadOnly = !member.ReadOnly
	member.IsAdmin = false
	member.IsModerator = false

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		logbuch.Error("Error saving when toggling read only for member", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return errs.Saving
	}

	if err := updateDefaultGroupMembers(tx, orga.ID, memberUserId, member); err != nil {
		logbuch.Error("Error updating default group when toggling organization member read only", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return err
	}

	if err := createToggleReadOnlyFeed(tx, orga, userId, member); err != nil {
		logbuch.Error("Error creating feed when toggling read only for member", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when toggling read only member", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	go func() {
		if err := billing.UpdateSubscription(orga); err != nil {
			logbuch.Error("Error updating subscription while toggling read only member", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}
	}()

	return nil
}

func createToggleReadOnlyFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, member *model.OrganizationMember) error {
	reason := "set_member_read_only"

	if !member.ReadOnly {
		reason = "remove_member_read_only"
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
