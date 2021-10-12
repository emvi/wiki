package member

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/constants"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
)

func ToggleAdmin(orga *model.Organization, userId, memberUserId hide.ID) error {
	if orga.OwnerUserId == memberUserId {
		return errs.OrganizationOwner
	}

	admin, member, err := checkIsAdminAndGetMember(orga.ID, userId, memberUserId)

	if err != nil {
		return err
	}

	if admin.UserId == memberUserId {
		return errs.ChangeAdminYourself
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when toggling admin", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	member.IsAdmin = !member.IsAdmin
	member.ReadOnly = false

	if member.IsAdmin {
		member.IsModerator = true
	}

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		logbuch.Error("Error toggling organization member admin", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return errs.Saving
	}

	if err := updateDefaultGroupMembers(tx, orga.ID, memberUserId, member); err != nil {
		logbuch.Error("Error updating default group when toggling organization member admin", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return err
	}

	if err := createToggleAdminFeed(tx, orga, userId, member); err != nil {
		logbuch.Error("Error creating feed when toggling organization member admin", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_user_id": memberUserId})
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when toggling admin", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func updateDefaultGroupMembers(tx *sqlx.Tx, orgaId, memberUserId hide.ID, member *model.OrganizationMember) error {
	if err := updateDefaultGroupMember(tx, orgaId, memberUserId, constants.GroupAdminName, member.IsAdmin); err != nil {
		return err
	}

	if err := updateDefaultGroupMember(tx, orgaId, memberUserId, constants.GroupModName, member.IsModerator); err != nil {
		return err
	}

	if err := updateDefaultGroupMember(tx, orgaId, memberUserId, constants.GroupReadOnlyName, member.ReadOnly); err != nil {
		return err
	}

	return nil
}

func updateDefaultGroupMember(tx *sqlx.Tx, orgaId, memberUserId hide.ID, name string, add bool) error {
	group := model.GetUserGroupByOrganizationIdAndNameTx(tx, orgaId, name)

	if group == nil {
		logbuch.Error("Default user group not found when toggling member role", logbuch.Fields{"orga_id": orgaId, "user_id": memberUserId})
		db.Rollback(tx)
		return errs.Saving
	}

	member := model.GetUserGroupMemberByGroupIdAndUserIdTx(tx, group.ID, memberUserId)

	if add {
		if member == nil {
			member := &model.UserGroupMember{UserGroupId: group.ID, UserId: memberUserId}

			if err := model.SaveUserGroupMember(tx, member); err != nil {
				return errs.Saving
			}
		}
	} else {
		if member == nil {
			return nil
		}

		if err := model.DeleteUserGroupMemberByUserGroupIdAndId(tx, group.ID, member.ID); err != nil {
			logbuch.Error("Error deleting default user group member when toggling member role", logbuch.Fields{"orga_id": orgaId, "user_id": memberUserId})
			return errs.Saving
		}
	}

	return nil
}

func createToggleAdminFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, member *model.OrganizationMember) error {
	reason := "set_organization_admin"

	if !member.IsAdmin {
		reason = "remove_organization_admin"
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
