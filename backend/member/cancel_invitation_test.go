package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestCancelInvitation(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, admin := testutil.CreateOrgaAndUser(t)
	member := testutil.CreateUser(t, orga, 222, "member@test.com")
	user1 := testutil.CreateUserWithoutOrganization(t, 333, "new1@test.com")
	user2 := testutil.CreateUserWithoutOrganization(t, 444, "new2@test.com")
	invitation1 := testutil.CreateInvitation(t, orga, "new1@test.com", "code", false)
	invitation2 := testutil.CreateInvitation(t, orga, "new2@test.com", "code", false)

	input := []struct {
		UserId       hide.ID
		InvitationId hide.ID
	}{
		{member.ID, invitation1.ID},
		{admin.ID, 0},
		{admin.ID, invitation1.ID},
	}
	expected := []error{
		errs.PermissionDenied,
		errs.InvitationNotFound,
		nil,
	}

	for i, in := range input {
		if err := CancelInvitation(orga, in.UserId, in.InvitationId); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	if model.GetInvitationByEmailAndId(user1.Email, invitation1.ID) != nil {
		t.Fatal("First invitation must not exist anymore")
	}

	if model.GetInvitationByEmailAndId(user2.Email, invitation2.ID) == nil {
		t.Fatal("Second invitation must still exist")
	}
}
