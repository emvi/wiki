package member

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func CancelInvitation(orga *model.Organization, userId, invitationId hide.ID) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	if model.GetInvitationByOrganizationIdAndId(orga.ID, invitationId) == nil {
		return errs.InvitationNotFound
	}

	if err := model.DeleteInvitationById(nil, invitationId); err != nil {
		logbuch.Error("Error deleting invitation while cancelling invitation", logbuch.Fields{"err": err})
		return errs.Saving
	}

	return nil
}
