package member

import (
	"emviwiki/backend/billing"
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func RemoveMember(orga *model.Organization, userId, memberUserId hide.ID, removePermissions bool) error {
	admin, member, err := checkIsAdminAndGetMember(orga.ID, userId, memberUserId)

	if err != nil {
		return err
	}

	if admin.UserId == memberUserId {
		return errs.RemoveYourself
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction when removing member", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member.Active = false

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		logbuch.Error("Error deactivating organization member", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return errs.Saving
	}

	if err := updateObjectPermissionsForInactiveUser(tx, orga, memberUserId, removePermissions); err != nil {
		return err
	}

	if err := createLeaveOrganizationFeed(tx, orga, member.User); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when removing member", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	go func() {
		if err := billing.UpdateSubscription(orga); err != nil {
			logbuch.Error("Error updating subscription while removing member", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}
	}()

	return nil
}

func checkIsAdminAndGetMember(orgaId, userId, memberUserId hide.ID) (*model.OrganizationMember, *model.OrganizationMember, error) {
	admin, err := perm.CheckUserIsAdmin(orgaId, userId)

	if err != nil {
		return nil, nil, err
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orgaId, memberUserId)

	if member == nil {
		return nil, nil, errs.MemberNotFound
	}

	return admin, member, nil
}
