package organization

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadOrganizations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	orgs := ReadOrganizations(user.ID)

	if len(orgs) != 1 || orgs[0].ID != orga.ID || !orgs[0].IsAdmin {
		t.Fatal("User must be admin of one organization")
	}
}

func TestReadOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateUser(t, orga, 321, "member@user.com")

	if organization, err := ReadOrganization(context.NewEmviUserContext(orga, user.ID)); err != nil || organization == nil {
		t.Fatalf("Organization must be found, but was: %v", err)
	}
}

func TestReadOrganizationClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	testutil.CreateUser(t, orga, 321, "member@user.com")
	ctx := context.NewEmviContext(orga, 0, nil, false)

	if organization, err := ReadOrganization(ctx); err != nil || organization == nil {
		t.Fatalf("Organization must be found, but was: %v", err)
	}
}

func TestGetOrganizationStatistics(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "non@admin.com")
	user2.OrganizationMember.ReadOnly = true

	if err := model.SaveOrganizationMember(nil, user2.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateTag(t, orga, "tag")

	_, err := GetOrganizationStatistics(orga, user2.ID)

	if err != errs.PermissionDenied {
		t.Fatalf("Expected permission denied, but was: %v", err)
	}

	statistics, err := GetOrganizationStatistics(orga, user.ID)

	if err != nil {
		t.Fatalf("Statistics must be returned, but was: %v", err)
	}

	if statistics.ArticleCount != 1 ||
		statistics.ListCount != 1 ||
		statistics.GroupCount != 5 ||
		statistics.MemberCount != 2 ||
		statistics.BillableMemberCount != 1 ||
		statistics.TagCount != 5 ||
		int(statistics.MaxStorage) != constants.DefaultMaxStorageGb ||
		int(statistics.StorageUsage) != 0 {
		t.Fatalf("Statistics not as expected: %v", statistics)
	}
}
