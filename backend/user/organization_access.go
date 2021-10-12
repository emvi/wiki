package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func HasAccessToOrganization(userId, orgaId hide.ID) error {
	if model.GetOrganizationMemberByOrganizationIdAndUserId(orgaId, userId) == nil {
		return errs.PermissionDenied
	}

	return nil
}
