package organization

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
)

func UpdateOrganization(orga *model.Organization, userId hide.ID, name, domain string) []error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return []error{err}
	}

	name = strings.TrimSpace(name)
	domain = strings.ToLower(strings.TrimSpace(domain))
	err := make([]error, 0)
	data := CreateOrganizationData{}

	if e := data.CheckNameValid(name); e != nil {
		err = append(err, e)
	}

	if e := data.CheckDomainValid(domain); e != nil {
		err = append(err, e)
	}

	if o := model.GetOrganizationByNameNormalized(domain); o != nil && o.ID != orga.ID {
		err = append(err, errs.DomainInUse)
	}

	if len(err) != 0 {
		return err
	}

	orga.Name = name
	orga.NameNormalized = domain

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization when changing name and/or domain", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "name": name, "domain": domain})
		return []error{errs.Saving}
	}

	return nil
}

func UpdateOrganizationPermissions(orga *model.Organization, userId hide.ID, createGroupAdmin, createGroupMod bool) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	orga.CreateGroupAdmin = createGroupAdmin || createGroupMod
	orga.CreateGroupMod = createGroupMod

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization when changing permissions", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId})
		return errs.Saving
	}

	return nil
}
