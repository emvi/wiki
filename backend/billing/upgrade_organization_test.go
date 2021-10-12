package billing

import (
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestUpgradeOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	orga.Expert = false
	orga.MaxStorageGB = constants.DefaultMaxStorageGb

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	testutil.CreateUser(t, orga, 321, "nonadmin@test.com")
	readonly := testutil.CreateUser(t, orga, 322, "readonly@test.com")
	readonly.OrganizationMember.ReadOnly = true

	if err := model.SaveOrganizationMember(nil, readonly.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	inactive := testutil.CreateUser(t, orga, 323, "inactive@test.com")
	inactive.OrganizationMember.Active = false

	if err := model.SaveOrganizationMember(nil, inactive.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := UpgradeOrganization(orga); err != nil {
		t.Fatalf("Organization must have been upgraded, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if !orga.Expert ||
		orga.MaxStorageGB != constants.StorageGBPerUser*2 ||
		!orga.SubscriptionCycle.Valid ||
		orga.SubscriptionCycle.Time.Equal(yesterday()) {
		t.Fatalf("Organization not as expected: %v %v %v", orga.Expert, orga.MaxStorageGB, orga.SubscriptionCycle.Time)
	}
}
