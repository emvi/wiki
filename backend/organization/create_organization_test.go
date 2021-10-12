package organization

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestCreateOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	user := testutil.CreateUserWithoutOrganization(t, 2, "test@user.com")
	user2 := testutil.CreateUserWithoutOrganization(t, 3, "test2@user.com")
	org, _ := testutil.CreateOrga(t, user2, "name")

	if err := CreateOrganization(user.ID, CreateOrganizationData{"name", "name", "username", "en"}); err[0] != errs.DomainInUse {
		t.Fatalf("Domain must be in use already, but was: %v", err)
	}

	if err := CreateOrganization(user.ID, CreateOrganizationData{"org Name 123", "org-Name-123", "username", "en"}); err != nil {
		t.Fatalf("Organization must be created: %v", err)
	}

	// check organization was created and can be found by name
	org = model.GetOrganizationByName("org name 123")

	if org == nil {
		t.Fatal("New organization must exist")
	}

	if org.NameNormalized != "org-name-123" {
		t.Fatal("Normalized name wrong")
	}

	// check organization has one member
	member := model.FindOrganizationMemberByOrganizationId(org.ID)

	if len(member) != 1 {
		t.Fatalf("New organization must have one member, but was: %v", len(member))
	}

	if member[0].SendNotificationsInterval != defaultNotificationInterval ||
		!member[0].RecommendationMail ||
		!member[0].ShowActionButtons ||
		!member[0].ShowNavigation ||
		!member[0].ShowCreateButton {
		t.Fatalf("Member not as expected: %v", member[0])
	}

	// check languages are created
	langs := model.FindLanguagesByOrganizationId(org.ID)

	if len(langs) != 1 {
		t.Fatal("Organization must have one language")
	}

	if langs[0].Name != "English" || langs[0].Code != "en" || !langs[0].Default {
		t.Fatal("First language must be en and default")
	}
}

func TestCreateOrganizationAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	testutil.CreateUserWithoutOrganization(t, 2, "test@user.com")

	if err := CreateOrganization(2, CreateOrganizationData{"name", "name", "username", "en"}); err != nil {
		t.Fatal("Organization must be created")
	}

	orgs := model.FindOrganizationsByUserId(2)

	if len(orgs) != 1 {
		t.Fatal("New organization must exist")
	}
}

func TestCreateOrganizationNameNormalized(t *testing.T) {
	testutil.CleanBackendDb(t)
	testutil.CreateUserWithoutOrganization(t, 2, "test@user.com")
	user2 := testutil.CreateUserWithoutOrganization(t, 3, "test@user.com")
	org := &model.Organization{Name: "name", NameNormalized: "name-123", OwnerUserId: 2}

	if err := model.SaveOrganization(nil, org); err != nil {
		t.Fatal(err)
	}

	if err := CreateOrganization(user2.ID, CreateOrganizationData{"name 123", "Name-123", "username", "en"}); err[0] != errs.DomainInUse {
		t.Fatalf("Name normalized must be in use already, but was: %v", err)
	}
}

func TestCreateOrganizationMaxFreeOrganizations(t *testing.T) {
	testutil.CleanBackendDb(t)
	testutil.CreateUserWithoutOrganization(t, 2, "test@user.com")
	org := &model.Organization{Name: "name", NameNormalized: "name-123", OwnerUserId: 2}

	if err := model.SaveOrganization(nil, org); err != nil {
		t.Fatal(err)
	}

	if err := CreateOrganization(2, CreateOrganizationData{"name 123", "name-123", "username", "en"}); err[0] != errs.MaxFreeOrganizationsReached {
		t.Fatalf("Maximum free organization must have been reached, but was: %v", err)
	}
}

func TestCreateOrganizationExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	testutil.CreateUserWithoutOrganization(t, 2, "test@user.com")
	org := &model.Organization{Name: "name", NameNormalized: "name-123", OwnerUserId: 2, Expert: true}

	if err := model.SaveOrganization(nil, org); err != nil {
		t.Fatal(err)
	}

	orga := model.GetOrganizationByName("name")

	if orga == nil || !orga.Expert {
		t.Fatal("Organization must be expert")
	}
}

func TestCreateOrganizationStandardGroups(t *testing.T) {
	testutil.CleanBackendDb(t)
	user := testutil.CreateUserWithoutOrganization(t, 2, "test@user.com")
	data := CreateOrganizationData{"name", "name", "username", "en"}

	if err := CreateOrganization(user.ID, data); len(err) != 0 {
		t.Fatalf("Organization must have been created, but was: %v", err)
	}

	orga := model.GetOrganizationByName("name")

	if orga == nil {
		t.Fatal("Organization must exist")
	}

	all := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupAllName)
	admin := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupAdminName)
	mod := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupModName)
	readonly := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupReadOnlyName)

	if all == nil || all.Info.String != constants.GroupAllInfo {
		t.Fatalf("Default group all must have been created, but was: %v", all)
	}

	if admin == nil || admin.Info.String != constants.GroupAdminInfo {
		t.Fatalf("Default group admin must have been created, but was: %v", all)
	}

	if mod == nil || mod.Info.String != constants.GroupModInfo {
		t.Fatalf("Default group moderator must have been created, but was: %v", all)
	}

	if readonly == nil || readonly.Info.String != constants.GroupReadOnlyInfo {
		t.Fatalf("Default group readonly must have been created, but was: %v", all)
	}
}

func TestCheckDomainValid(t *testing.T) {
	testutil.CleanBackendDb(t)
	validNames := []string{"name", "name123", "Name123", "name-123", "Name-1-2-3test"}
	invalidNames := []string{"123name", "-name", "name-", "näme", "with space"}

	for _, name := range validNames {
		data := CreateOrganizationData{Domain: name}

		if data.CheckDomainValid(name) != nil {
			t.Fatalf("Name '%v' must be valid", name)
		}
	}

	for _, name := range invalidNames {
		data := CreateOrganizationData{Domain: name}

		if data.CheckDomainValid(name) == nil {
			t.Fatalf("Name '%v' must be invalid", name)
		}
	}
}

func TestCheckDomainNotAllowed(t *testing.T) {
	testutil.CleanBackendDb(t)

	if err := model.SaveDomainBlacklist(nil, &model.DomainBlacklist{Name: "test"}); err != nil {
		t.Fatal(err)
	}

	data := CreateOrganizationData{Domain: "test"}

	if data.CheckDomainValid("test") != errs.DomainNotAllowed {
		t.Fatal("Domain must not be allowed")
	}

	data.Domain = "okay"

	if data.CheckDomainValid("okay") != nil {
		t.Fatal("Domain must be allowed")
	}
}

func TestCheckUsernameValid(t *testing.T) {
	validNames := []string{"name", "name123", "Name123", "name-123", "Name_123test", "01Müller-123"}
	invalidNames := []string{"-name", "name-", "name!"}

	for _, name := range validNames {
		data := CreateOrganizationData{Username: name}

		if data.CheckUsernameValid(name) != nil {
			t.Fatalf("Username '%v' must be valid", name)
		}
	}

	for _, name := range invalidNames {
		data := CreateOrganizationData{Username: name}

		if data.CheckUsernameValid(name) == nil {
			t.Fatalf("Username '%v' must be invalid", name)
		}
	}
}

func TestCheckNameValid(t *testing.T) {
	data := CreateOrganizationData{Name: ""}

	if data.CheckNameValid(data.Name) != errs.NameTooShort {
		t.Fatal("Name must not be valid")
	}

	data.Name = "0123456789012345678901234567890123456789012345678901234567891"

	if data.CheckNameValid(data.Name) != errs.NameTooLong {
		t.Fatal("Name must not be valid")
	}

	data.Name = "012345678901234567890123456789012345678901234567890123456789"

	if data.CheckNameValid(data.Name) != nil {
		t.Fatal("Name must be valid")
	}
}

func TestDomainValid(t *testing.T) {
	data := CreateOrganizationData{Domain: ""}

	if data.CheckDomainValid(data.Domain) != errs.DomainTooShort {
		t.Fatal("Domain must not be valid")
	}

	data.Domain = "d12345678901234567891"

	if data.CheckDomainValid(data.Domain) != errs.DomainTooLong {
		t.Fatal("Domain must not be valid")
	}

	data.Domain = "d1234567890123456789"

	if data.CheckDomainValid(data.Domain) != nil {
		t.Fatal("Domain must be valid")
	}
}

func TestCreateIntroduction(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)

	if err := createIntroduction(nil, orga, user.ID, lang); err != nil {
		t.Fatal(err)
	}
}
