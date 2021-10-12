package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestDeleteInvitation(t *testing.T) {
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
		Id     hide.ID
	}{
		{0, 0},
		{user.ID, 0},
		{0, invitation.ID},
		{user.ID, invitation.ID},
	}
	expected := []error{
		errs.UserNotFound,
		errs.InvitationNotFound,
		errs.UserNotFound,
		nil,
	}

	for i, in := range input {
		if err := DeleteInvitation(in.UserId, in.Id); err != expected[i] {
			t.Fatalf("Expected '%v' but was: %v", expected[i], err)
		}
	}
}
