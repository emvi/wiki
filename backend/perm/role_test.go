package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestCheckUserIsAdmin(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	user := testutil.CreateUser(t, orga, 321, "tester@testutil.com")

	if member, err := CheckUserIsAdmin(orga.ID, user.ID); err != errs.PermissionDenied || member != nil {
		t.Fatalf("Expected user not to be an admin, but was: %v %v", err, member)
	}

	if member, err := CheckUserIsAdmin(orga.ID, admin.ID); err != nil || member == nil {
		t.Fatalf("Expected user to be an admin, but was: %v %v", err, member)
	}
}

func TestCheckUserIsMod(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, mod := testutil.CreateOrgaAndUser(t)
	user := testutil.CreateUser(t, orga, 321, "tester@testutil.com")

	if member, err := CheckUserIsMod(orga.ID, user.ID); err != errs.PermissionDenied || member != nil {
		t.Fatalf("Expected user not to be a moderator, but was: %v %v", err, member)
	}

	if member, err := CheckUserIsAdmin(orga.ID, mod.ID); err != nil || member == nil {
		t.Fatalf("Expected user to be a moderator, but was: %v %v", err, member)
	}
}

func TestCheckUserIsAdminOrMod(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	mod := testutil.CreateUser(t, orga, 322, "mod@testutil.com")
	user := testutil.CreateUser(t, orga, 321, "tester@testutil.com")

	mod.OrganizationMember.IsModerator = true

	if err := model.SaveOrganizationMember(nil, mod.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if member, err := CheckUserIsAdminOrMod(orga.ID, user.ID); err != errs.PermissionDenied || member != nil {
		t.Fatalf("Expected user not to be an admin or moderator, but was: %v %v", err, member)
	}

	if member, err := CheckUserIsAdminOrMod(orga.ID, admin.ID); err != nil || member == nil {
		t.Fatalf("Expected user to be an admin or moderator, but was: %v %v", err, member)
	}

	if member, err := CheckUserIsAdminOrMod(orga.ID, mod.ID); err != nil || member == nil {
		t.Fatalf("Expected user to be an admin or moderator, but was: %v %v", err, member)
	}
}
