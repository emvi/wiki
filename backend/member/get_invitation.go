package member

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func GetInvitation(userId hide.ID, code string) (*model.Organization, error) {
	user := model.GetUserById(userId)

	if user == nil {
		logbuch.Warn("Invitation user not found", logbuch.Fields{"user_id": userId, "code": code})
		return nil, errs.InvitationNotFound
	}

	invitation := model.GetInvitationByEmailAndCode(user.Email, code)

	if invitation == nil {
		logbuch.Warn("Invitation by email and code not found", logbuch.Fields{"user_id": userId, "code": code})
		return nil, errs.InvitationNotFound
	}

	orga := model.GetOrganizationById(invitation.OrganizationId)

	if orga == nil {
		logbuch.Warn("Invitation organization not found", logbuch.Fields{"user_id": userId, "code": code})
		return nil, errs.InvitationNotFound
	}

	return orga, nil
}

func ReadInvitations(userId hide.ID) ([]model.Invitation, error) {
	user := model.GetUserById(userId)

	if user == nil {
		return nil, errs.UserNotFound
	}

	return model.FindInvitationsByEmail(user.Email), nil
}

func ReadOrganizationInvitations(orga *model.Organization, userId hide.ID) ([]model.Invitation, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return nil, err
	}

	return model.FindInvitationsByOrganizationId(orga.ID), nil
}
