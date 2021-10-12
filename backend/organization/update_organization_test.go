package organization

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestUpdateOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga2, _ := testutil.CreateOrga(t, user, "diff")

	input := []struct {
		UserId hide.ID
		Name   string
		Domain string
	}{
		{user.ID + 1, "", ""},
		{user.ID, "", "domain"},
		{user.ID, "name", ""},
		{user.ID, "newname", orga2.NameNormalized},
		{user.ID, "newname", "newdomain"},
	}
	expected := []error{
		errs.PermissionDenied,
		errs.NameTooShort,
		errs.DomainTooShort,
		errs.DomainInUse,
		nil,
	}

	for i, in := range input {
		err := UpdateOrganization(orga, in.UserId, in.Name, in.Domain)

		if expected[i] == nil && len(err) != 0 {
			t.Fatalf("Expected organization to be updated, but was: %v", err)
		} else if expected[i] != nil && (len(err) != 1 || err[0] != expected[i]) {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.Name != "newname" || orga.NameNormalized != "newdomain" {
		t.Fatalf("Expected organization to be updated, but was: %v", orga)
	}
}

func TestUpdateOrganizationPermissions(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if err := UpdateOrganizationPermissions(orga, user.ID+1, true, true); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	if err := UpdateOrganizationPermissions(orga, user.ID, true, false); err != nil {
		t.Fatalf("Permission must have been updated, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if !orga.CreateGroupAdmin || orga.CreateGroupMod {
		t.Fatalf("Wrong permissions: %v %v", orga.CreateGroupAdmin, orga.CreateGroupMod)
	}

	if err := UpdateOrganizationPermissions(orga, user.ID, false, true); err != nil {
		t.Fatalf("Permission must have been updated, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if !orga.CreateGroupAdmin || !orga.CreateGroupMod {
		t.Fatalf("Wrong permissions: %v %v", orga.CreateGroupAdmin, orga.CreateGroupMod)
	}

	if err := UpdateOrganizationPermissions(orga, user.ID, false, false); err != nil {
		t.Fatalf("Permission must have been updated, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.CreateGroupAdmin || orga.CreateGroupMod {
		t.Fatalf("Wrong permissions: %v %v", orga.CreateGroupAdmin, orga.CreateGroupMod)
	}
}
