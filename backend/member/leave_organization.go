package member

import (
	"emviwiki/backend/billing"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func LeaveOrganization(orga *model.Organization, userId hide.ID, orgaName string) error {
	if orga.OwnerUserId == userId {
		return errs.OrganizationOwnerLeave
	}

	if orga.Name != orgaName {
		return errs.NameDoesNotMatch
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to leave organization", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserIdTx(tx, orga.ID, userId)

	if member == nil {
		logbuch.Error("Member not found when leaving organization", logbuch.Fields{"orga_id": orga.ID, "user_id": userId})
		return errs.UserNotFound
	}

	member.Active = false

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		logbuch.Error("Error deactivating organization member when leaving organization", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId})
		return errs.Saving
	}

	if err := updateObjectPermissionsForInactiveUser(tx, orga, userId, false); err != nil {
		return err
	}

	if err := createLeaveOrganizationFeed(tx, orga, member.User); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when leaving organization", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	go func() {
		if err := billing.UpdateSubscription(orga); err != nil {
			logbuch.Error("Error updating subscription while leaving organization", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}
	}()

	return nil
}

func createLeaveOrganizationFeed(tx *sqlx.Tx, orga *model.Organization, user *model.User) error {
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       user.ID,
		Reason:       "left_organization",
		Public:       true}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when leaving organization", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
