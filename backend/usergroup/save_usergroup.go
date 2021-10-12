package usergroup

import (
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"unicode/utf8"

	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
)

const (
	maxNameLen = 40
	maxInfoLen = 100
)

type SaveUserGroupData struct {
	Id   hide.ID `json:"id"`
	Name string  `json:"name"`
	Info string  `json:"info"`
}

func (data *SaveUserGroupData) validate(organization *model.Organization) []error {
	data.Name = strings.TrimSpace(data.Name)
	data.Info = strings.TrimSpace(data.Info)

	err := make([]error, 0)

	if data.Name == "" {
		err = append(err, errs.NameTooShort)
	} else if utf8.RuneCountInString(data.Name) > maxNameLen {
		err = append(err, errs.NameTooLong)
	}

	if utf8.RuneCountInString(data.Info) > maxInfoLen {
		err = append(err, errs.InfoTooLong)
	}

	group := model.GetUserGroupByOrganizationIdAndName(organization.ID, data.Name)

	if group != nil && group.ID != data.Id {
		err = append(err, errs.NameInUse)
	}

	if len(err) == 0 {
		return nil
	}

	return err
}

func SaveUserGroup(organization *model.Organization, userId hide.ID, data SaveUserGroupData) (hide.ID, []error) {
	if err := checkUserAccess(organization, userId, data.Id); err != nil {
		return 0, []error{err}
	}

	if err := data.validate(organization); err != nil {
		return 0, err
	}

	// update or create user group
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to save user group", logbuch.Fields{"err": err})
		return 0, []error{errs.TxBegin}
	}

	var group *model.UserGroup
	newGroup := false

	if data.Id != 0 {
		group = model.GetUserGroupByOrganizationIdAndIdTx(tx, organization.ID, data.Id)

		if group == nil {
			db.Rollback(tx)
			return 0, []error{errs.GroupNotFound}
		}
	} else {
		group = &model.UserGroup{OrganizationId: organization.ID}
		newGroup = true
	}

	group.Name = data.Name
	group.Info = null.NewString(data.Info, data.Info != "")

	if err := model.SaveUserGroup(tx, group); err != nil {
		return 0, []error{errs.Saving}
	}

	// create first user if group is new
	if newGroup {
		member := &model.UserGroupMember{UserGroupId: group.ID,
			UserId:      userId,
			IsModerator: true}

		if err := model.SaveUserGroupMember(tx, member); err != nil {
			return 0, []error{errs.Saving}
		}
	}

	if err := createSaveUserGroupFeed(tx, organization, userId, group, data.Id == 0); err != nil {
		logbuch.Error("Error creating feed when saving user group", logbuch.Fields{"err": err})
		return 0, []error{err}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when saving user group", logbuch.Fields{"err": err})
		return 0, []error{errs.TxCommit}
	}

	return group.ID, nil
}

func checkUserAccess(orga *model.Organization, userId, groupId hide.ID) error {
	if orga.CreateGroupAdmin || orga.CreateGroupMod {
		orgaMember := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, userId)

		if orgaMember == nil ||
			!(orga.CreateGroupMod && orgaMember.IsModerator || orga.CreateGroupAdmin && orgaMember.IsAdmin) {
			return errs.PermissionDenied
		}
	}

	if groupId == 0 {
		return nil
	}

	member := model.GetUserGroupMemberByGroupIdAndUserId(groupId, userId)

	if member == nil {
		logbuch.Warn("User group member not found", logbuch.Fields{"organization_id": orga.ID, "user_id": userId, "group_id": groupId})
		return errs.PermissionDenied
	}

	if !member.IsModerator {
		return errs.PermissionDenied
	}

	return nil
}

func createSaveUserGroupFeed(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, group *model.UserGroup, isNew bool) error {
	var notify []hide.ID
	reason := "create_user_group"

	if !isNew {
		reason = "update_user_group"
		notify = model.FindObservedObjectUserIdByUserGroupIdTx(tx, group.ID)
	}

	refs := make([]interface{}, 1)
	refs[0] = group
	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       userId,
		Reason:       reason,
		Public:       true,
		Notify:       notify,
		Refs:         refs}

	if err := feed.CreateFeed(feedData); err != nil {
		return err
	}

	return nil
}
