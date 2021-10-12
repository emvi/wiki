package invitation

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
	"time"
)

func TestCleanupInvitations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	invitation1 := testutil.CreateInvitation(t, orga, "test1@user.com", "code1", false)
	invitation2 := testutil.CreateInvitation(t, orga, "test2@user.com", "code2", true)
	invitation3 := testutil.CreateInvitation(t, orga, "test3@user.com", "code3", false)
	invitation4 := testutil.CreateInvitation(t, orga, "test4@user.com", "code4", false)
	setInvitationDefTime(t, invitation1, time.Now().Add(-time.Hour*24*31))
	setInvitationDefTime(t, invitation2, time.Now().Add(-time.Hour*24*42))
	setInvitationDefTime(t, invitation3, time.Now().Add(-time.Hour*24*28))
	setInvitationDefTime(t, invitation4, time.Now().Add(-time.Hour*24*11))
	CleanupInvitations()

	if model.GetInvitationByEmailAndCode("test1@user.com", "code1") != nil {
		t.Fatal("First invitation must not exist")
	}

	if model.GetInvitationByEmailAndCode("test2@user.com", "code2") != nil {
		t.Fatal("Second invitation must not exist")
	}

	if model.GetInvitationByEmailAndCode("test3@user.com", "code3") == nil {
		t.Fatal("Third invitation must exist")
	}

	if model.GetInvitationByEmailAndCode("test4@user.com", "code4") == nil {
		t.Fatal("Forth invitation must exist")
	}
}

func setInvitationDefTime(t *testing.T, invitation *model.Invitation, defTime time.Time) {
	if _, err := model.GetConnection().Exec(nil, `UPDATE "invitation" SET def_time = $2 WHERE id = $1`, invitation.ID, defTime); err != nil {
		t.Fatal(err)
	}
}
