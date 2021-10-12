package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func CheckUserTagAccess(orgaId, userId, tagId hide.ID) error {
	if model.GetTagByOrganizationIdAndUserIdAndId(orgaId, userId, tagId) == nil {
		return errs.PermissionDenied
	}

	return nil
}
