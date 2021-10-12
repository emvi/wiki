package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadUserProfileId(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if user, err := ReadUserProfile(orga, user.ID+1, ""); err != errs.UserNotFound || user != nil {
		t.Fatal("User must not have been found")
	}

	if user, err := ReadUserProfile(orga, user.ID, ""); err != nil || user == nil {
		t.Fatal("User must have been found")
	}
}

func TestReadUserProfileUsername(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if user, err := ReadUserProfile(orga, 0, "unknown"); err != errs.UserNotFound || user != nil {
		t.Fatal("User must not have been found")
	}

	if user, err := ReadUserProfile(orga, 0, user.OrganizationMember.Username); err != nil || user == nil {
		t.Fatal("User must have been found")
	}
}
