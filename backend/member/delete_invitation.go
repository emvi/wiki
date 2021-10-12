package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func DeleteInvitation(userId, invitationId hide.ID) error {
	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	invitation := model.GetInvitationByEmailAndId(user.Email, invitationId)

	if invitation == nil {
		return errs.InvitationNotFound
	}

	if err := model.DeleteInvitationById(nil, invitation.ID); err != nil {
		logbuch.Error("Error deleting invitation", logbuch.Fields{"err": err, "user_id": userId, "invitation_id": invitationId})
		return errs.Saving
	}

	return nil
}
