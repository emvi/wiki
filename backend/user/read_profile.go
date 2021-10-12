package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func ReadUserProfile(organization *model.Organization, userId hide.ID, username string) (*model.User, error) {
	var user *model.User

	if userId != 0 {
		user = model.GetUserWithOrganizationMemberByOrganizationIdAndId(organization.ID, userId)
	} else {
		user = model.GetUserWithOrganizationMemberByOrganizationIdAndUsername(organization.ID, username)
	}

	if user == nil {
		return nil, errs.UserNotFound
	}

	return user, nil
}
