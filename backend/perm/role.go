package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

// CheckUserIsAdmin checks if given user is an administrator of given organization.
// If that's the case, the organization member is returned. Else an error with permission denied is returned.
func CheckUserIsAdmin(orgaId, userId hide.ID) (*model.OrganizationMember, error) {
	admin := model.GetOrganizationMemberByOrganizationIdAndUserIdAndIsAdmin(orgaId, userId)

	if admin == nil {
		return nil, errs.PermissionDenied
	}

	return admin, nil
}

// CheckUserIsMod checks if given user is a moderator of given organization.
// If that's the case, the organization member is returned. Else an error with permission denied is returned.
func CheckUserIsMod(orgaId, userId hide.ID) (*model.OrganizationMember, error) {
	mod := model.GetOrganizationMemberByOrganizationIdAndUserIdAndIsMod(orgaId, userId)

	if mod == nil {
		return nil, errs.PermissionDenied
	}

	return mod, nil
}

func CheckUserIsAdminOrMod(orgaId, userId hide.ID) (*model.OrganizationMember, error) {
	admin, _ := CheckUserIsAdmin(orgaId, userId)

	if admin != nil {
		return admin, nil
	}

	mod, _ := CheckUserIsMod(orgaId, userId)

	if mod != nil {
		return mod, nil
	}

	return nil, errs.PermissionDenied
}
