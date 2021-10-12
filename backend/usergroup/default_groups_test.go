package usergroup

import (
	"emviwiki/shared/constants"
	"emviwiki/shared/testutil"
	"testing"
)

func TestGetAllGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	group := GetAllGroup(nil, orga.ID)

	if group == nil || group.Name != constants.GroupAllName {
		t.Fatalf("Group must be returned, but was: %v", group)
	}
}

func TestGetAdminGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	group := GetAdminGroup(nil, orga.ID)

	if group == nil || group.Name != constants.GroupAdminName {
		t.Fatalf("Group must be returned, but was: %v", group)
	}
}

func TestGetModGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	group := GetModGroup(nil, orga.ID)

	if group == nil || group.Name != constants.GroupModName {
		t.Fatalf("Group must be returned, but was: %v", group)
	}
}

func TestGetReadOnlyGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	group := GetReadOnlyGroup(nil, orga.ID)

	if group == nil || group.Name != constants.GroupReadOnlyName {
		t.Fatalf("Group must be returned, but was: %v", group)
	}
}
