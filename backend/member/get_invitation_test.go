package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestGetInvitation(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	invitation := &model.Invitation{OrganizationId: orga.ID,
		Email: user.Email,
		Code:  "code"}

	if err := model.SaveInvitation(nil, invitation); err != nil {
		t.Fatal(err)
	}

	input := []struct {
		UserId hide.ID
		Code   string
	}{
		{0, ""},
		{user.ID, ""},
		{0, "code"},
		{user.ID, "asdf"},
	}

	for _, in := range input {
		if _, err := GetInvitation(in.UserId, in.Code); err != errs.InvitationNotFound {
			t.Fatalf("Expected error '%v', but was: %v", errs.InvitationNotFound, err)
		}
	}

	org, err := GetInvitation(user.ID, "code")

	if err != nil {
		t.Fatalf("Expected no error, but was: %v", err)
	}

	if org.ID != orga.ID {
		t.Fatal("Organization must be returned")
	}
}

func TestReadInvitations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	invitation := &model.Invitation{OrganizationId: orga.ID,
		Email: user.Email,
		Code:  "code1"}
	invitation2 := &model.Invitation{OrganizationId: orga.ID,
		Email: user.Email,
		Code:  "code2"}

	if err := model.SaveInvitation(nil, invitation); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveInvitation(nil, invitation2); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, "UPDATE invitation SET def_time = NOW() - INTERVAL '31 days' WHERE id = $1", invitation2.ID); err != nil {
		t.Fatal(err)
	}

	invitations, err := ReadInvitations(user.ID)

	if err != nil || len(invitations) != 1 {
		t.Fatalf("Must return one invitation, but was: %v %v", len(invitations), err)
	}

	invitation = model.GetInvitationByEmailAndCode(user.Email, "code1")

	if invitation == nil {
		t.Fatalf("Invitation must have user id set, but was: %v", invitation)
	}
}

func TestReadOrganizationInvitations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	member := testutil.CreateUser(t, orga, 321, "member@test.com")
	testutil.CreateInvitation(t, orga, "user1@test.com", "code1", false)
	testutil.CreateInvitation(t, orga, "user2@test.com", "code2", true)

	if _, err := ReadOrganizationInvitations(orga, member.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	invitations, err := ReadOrganizationInvitations(orga, admin.ID)

	if err != nil {
		t.Fatalf("Invitations must be returned, but was: %v", err)
	}

	if len(invitations) != 2 || invitations[0].Code != "code1" || invitations[1].Code != "code2" {
		t.Fatal("Two invitations must be returned")
	}
}
