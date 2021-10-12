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

func DeleteUserGroup(organization *model.Organization, userId, groupId hide.ID) error {
	if _, err := checkGroupExists(organization, groupId); err != nil {
		return err
	}

	if err := checkUserAccess(organization, userId, groupId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error beginning transaction when deleting user group", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	// create feed before deleting list
	if err := createDeletedGroupFeed(tx, organization, userId, groupId); err != nil {
		logbuch.Error("Error creating feed when deleting user group", logbuch.Fields{"err": err})
		return err
	}

	if err := feed.DeleteFeed(tx, &feed.DeleteFeedData{GroupId: groupId}); err != nil {
		logbuch.Error("Error deleting user group feed", logbuch.Fields{"err": err, "group_id": groupId, "user_id": userId})
		return errs.Saving
	}

	if err := model.DeleteUserGroupById(tx, groupId); err != nil {
		return errs.Saving
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when deleting user group", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func checkGroupExists(organization *model.Organization, groupId hide.ID) (*model.UserGroup, error) {
	group := model.GetUserGroupByOrganizationIdAndId(organization.ID, groupId)

	if group == nil {
		return nil, errs.GroupNotFound
	}

	return group, nil
}

func createDeletedGroupFeed(tx *sqlx.Tx, orga *model.Organization, userId, groupId hide.ID) error {
	group := model.GetUserGroupByOrganizationIdAndIdTx(tx, orga.ID, groupId)

	if group == nil {
		db.Rollback(tx)
		return errs.GroupNotFound
	}

	refs := make([]interface{}, 1)
	refs[0] = feed.KeyValue{"name", group.Name}
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       "delete_usergroup",
		Public:       true,
		Notify:       model.FindObservedObjectUserIdByUserGroupIdTx(tx, groupId),
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
