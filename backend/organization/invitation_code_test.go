package organization

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestGenerateInvitationCode(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")

	if _, err := GenerateInvitationCode(orga, user2.ID, false); err != errs.PermissionDenied {
		t.Fatalf("Permission must have been denied, but was: %v", err)
	}

	code, err := GenerateInvitationCode(orga, user.ID, true)

	if err != nil {
		t.Fatalf("Code must have been generated, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.InvitationCode.String != code || len(orga.InvitationCode.String) != invitationCodeLength {
		t.Fatalf("Invitation code must have been set, but was: %v", orga.InvitationCode)
	}

	if !orga.InvitationReadOnly {
		t.Fatal("Invitation link must be read only")
	}
}

func TestGetInvitationCode(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	orga.InvitationCode = null.NewString("invitationcode", true)
	orga.InvitationReadOnly = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if _, _, err := GetInvitationCode(orga, user2.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must have been denied, but was: %v", err)
	}

	if code, ro, err := GetInvitationCode(orga, user.ID); err != nil || code != "invitationcode" || !ro {
		t.Fatalf("Invitation code must have been returned, but was: %v %v %v", code, ro, err)
	}
}

func TestGetInvitationCodeOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUserWithoutOrganization(t, 321, "test2@user.com")

	if _, err := GetInvitationCodeOrganization(user2.ID, "invitationcode"); err != errs.OrganizationNotFound {
		t.Fatalf("Organization must not be found, but was: %v", err)
	}

	orga.InvitationCode = null.NewString("invitationcode", true)

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if _, err := GetInvitationCodeOrganization(user.ID, "invitationcode"); err != errs.IsMemberAlready {
		t.Fatalf("User must be a member already, but was: %v", err)
	}

	if orga, err := GetInvitationCodeOrganization(user2.ID, "invitationcode"); err != nil || orga == nil {
		t.Fatalf("Organization must be found, but was: %v", err)
	}
}
