package organization

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
)

const (
	invitationCodeLength = 20
)

func GenerateInvitationCode(orga *model.Organization, userId hide.ID, readOnly bool) (string, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return "", err
	}

	code := util.GenRandomString(invitationCodeLength)

	for model.GetOrganizationByInvitationCode(code) != nil {
		code = util.GenRandomString(invitationCodeLength)
	}

	orga.InvitationCode = null.NewString(code, true)
	orga.InvitationReadOnly = readOnly

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while generating invitation code", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId})
		return "", errs.Saving
	}

	return code, nil
}

func GetInvitationCode(orga *model.Organization, userId hide.ID) (string, bool, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return "", false, err
	}

	return orga.InvitationCode.String, orga.InvitationReadOnly, nil
}

func GetInvitationCodeOrganization(userId hide.ID, code string) (*model.Organization, error) {
	orga := model.GetOrganizationByInvitationCode(code)

	if orga == nil {
		return nil, errs.OrganizationNotFound
	}

	if model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, userId) != nil {
		return nil, errs.IsMemberAlready
	}

	return orga, nil
}
