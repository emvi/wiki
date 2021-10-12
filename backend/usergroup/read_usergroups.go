package usergroup

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/observe"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

const (
	usergroupMemberLimit = 20
)

func ReadUserGroup(organization *model.Organization, userId, groupId hide.ID) (*model.UserGroup, bool, bool, error) {
	group := model.GetUserGroupByOrganizationIdAndId(organization.ID, groupId)

	if group == nil {
		return nil, false, false, errs.GroupNotFound
	}

	isObserved := observe.IsObserved(userId, 0, 0, groupId)
	isMod := checkUserAccess(organization, userId, groupId) == nil
	return group, isMod, isObserved, nil
}

func ReadUserGroupMember(organization *model.Organization, groupId hide.ID, filter *model.SearchUserGroupMemberFilter) ([]model.UserGroupMember, int, error) {
	group := model.GetUserGroupByOrganizationIdAndId(organization.ID, groupId)

	if group == nil {
		return nil, 0, errs.GroupNotFound
	}

	if filter == nil {
		filter = new(model.SearchUserGroupMemberFilter)
	}

	filter.Limit = usergroupMemberLimit
	return model.FindUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimit(organization.ID, groupId, filter),
		model.CountUserGroupMemberByOrganizationIdAndUserGroupIdAndFilterLimit(organization.ID, groupId, filter), nil
}
