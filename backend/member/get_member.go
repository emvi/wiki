package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
)

func GetMember(organization *model.Organization, userId hide.ID) (*model.OrganizationMember, error) {
	member := model.GetOrganizationMemberByOrganizationIdAndUserId(organization.ID, userId)

	if member == nil {
		return nil, errs.MemberNotFound
	}

	return member, nil
}
