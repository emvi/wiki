package member

import (
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestUpdateObjectPermissionsForInactiveUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	user2 := testutil.CreateUser(t, orga, 321, "user2@testutil.com")
	groupIndrectAccess := testutil.CreateUserGroup(t, orga, "indirect")
	testutil.CreateUserGroupMember(t, groupIndrectAccess, user, true)

	// create articles
	testutil.CreateArticle(t, orga, user2, lang, true, true)               // don't update
	article2 := testutil.CreateArticle(t, orga, user2, lang, false, false) // don't update
	testutil.CreateArticleAccess(t, article2, user, nil, true)
	article3 := testutil.CreateArticle(t, orga, user2, lang, false, false) // don't update
	testSetArticlePrivate(t, article3)
	testutil.CreateArticleAccess(t, article3, user, nil, false)
	article4 := testutil.CreateArticle(t, orga, user2, lang, true, false) // update
	testutil.CreateArticleAccess(t, article4, user, nil, false)
	testutil.CreateArticle(t, orga, user2, lang, false, false) // don't update
	testutil.CreateArticleAccess(t, article4, nil, groupIndrectAccess, true)
	testutil.CreateArticle(t, orga, user2, lang, false, false) // update
	testutil.CreateArticleAccess(t, article4, nil, groupIndrectAccess, false)

	// create lists
	testutil.CreateArticleList(t, orga, user2, lang, true)              // update
	testutil.CreateArticleList(t, orga, user2, lang, false)             // update
	list2, _ := testutil.CreateArticleList(t, orga, user2, lang, false) // update
	testutil.CreateArticleListMember(t, list2, user.ID, 0, false)
	list3, _ := testutil.CreateArticleList(t, orga, user2, lang, false) // don't update
	testutil.CreateArticleListMember(t, list3, user.ID, 0, true)
	list4, _ := testutil.CreateArticleList(t, orga, user2, lang, false) // don't update
	testutil.CreateArticleListMember(t, list4, 0, groupIndrectAccess.ID, false)

	// create groups
	group1 := testutil.CreateUserGroup(t, orga, "group1") // update
	testutil.CreateUserGroupMember(t, group1, user2, true)
	group2 := testutil.CreateUserGroup(t, orga, "group2") // don't update
	testutil.CreateUserGroupMember(t, group2, user, true)

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		t.Fatal(err)
	}

	if err := updateObjectPermissionsForInactiveUser(tx, orga, user2.ID, false); err != nil {
		t.Fatalf("Objects must have been updated, but was: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	feeds := testutil.AssertFeedCreated(t, orga, "transfer_ownership")
	refs := model.FindFeedRefByOrganizationIdAndLanguageIdAndFeedId(orga.ID, lang.ID, feeds[0].ID)
	articleRefs := make([]model.FeedRef, 0)
	listRefs := make([]model.FeedRef, 0)
	groupRefs := make([]model.FeedRef, 0)

	for _, ref := range refs {
		if ref.ArticleID != 0 {
			articleRefs = append(articleRefs, ref)
		} else if ref.ArticleListID != 0 {
			listRefs = append(listRefs, ref)
		} else if ref.UserGroupID != 0 {
			groupRefs = append(groupRefs, ref)
		} else {
			t.Fatalf("Unknown feed reference: %v", ref)
		}
	}

	if len(articleRefs) != 2 {
		t.Fatalf("2 articles must have been transfered, but was: %v", len(articleRefs))
	}

	if len(listRefs) != 4 {
		t.Fatalf("4 lists must have been transfered, but was: %v", len(listRefs))
	}

	if len(groupRefs) != 1 {
		t.Fatalf("1 groups must have been transfered, but was: %v", len(groupRefs))
	}
}

func TestUpdateObjectPermissionsForInactiveUserRemovePermissions(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	user2 := testutil.CreateUser(t, orga, 321, "user2@testutil.com")

	// create articles
	article1 := testutil.CreateArticle(t, orga, user2, lang, false, false) // don't update
	testSetArticlePrivate(t, article1)
	article2 := testutil.CreateArticle(t, orga, user2, lang, false, false) // update

	// create lists
	list1, _ := testutil.CreateArticleList(t, orga, user2, lang, true) // update

	// create groups
	group1 := testutil.CreateUserGroup(t, orga, "group1") // update
	testutil.CreateUserGroupMember(t, group1, user2, true)

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		t.Fatal(err)
	}

	if err := updateObjectPermissionsForInactiveUser(tx, orga, user2.ID, true); err != nil {
		t.Fatalf("Objects must have been updated, but was: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	testutil.AssertFeedCreated(t, orga, "transfer_ownership")

	// check new counts for group
	group1 = model.GetUserGroupByOrganizationIdAndId(orga.ID, group1.ID)

	if group1.MemberCount != 1 {
		t.Fatalf("Users from admin/mod group must have been granted access to list, but was: %v", group1.MemberCount)
	}

	// check new permissions
	access1 := model.FindArticleAccessByOrganizationIdAndArticleId(orga.ID, article1.ID)

	if len(access1) != 1 && access1[0].ID != user2.ID {
		t.Fatalf("User must still have access to private article")
	}

	access2 := model.FindArticleAccessByOrganizationIdAndArticleId(orga.ID, article2.ID)

	if len(access2) != 2 {
		t.Fatalf("Admin and mod group must have access to article")
	}

	listMember1 := model.FindArticleListMemberByArticleListId(list1.ID)

	if len(listMember1) != 2 {
		t.Fatalf("Admin and mod group must have access to list")
	}

	groupMember1 := model.FindUserGroupMemberOnlyByUserGroupIdTx(nil, group1.ID)

	if len(groupMember1) != 1 {
		t.Fatalf("Admin user must have access to group")
	}
}

func TestUpdateArticlePermissionsForInactiveUser(t *testing.T) {
	orga, user2 := testSetupUserToUpdatePermissions(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user2, lang, false, false)
	tx, _ := model.GetConnection().Beginx()
	err := updateObjectPermissionsForInactiveUser(tx, orga, user2.ID, true)

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatalf("Article must have been transfered, but was: %v", err)
	}

	access := len(model.FindArticleAccessByOrganizationIdAndArticleId(orga.ID, article.ID))

	if access != 2 {
		t.Fatalf("Article must have only two members, but was: %v", access)
	}
}

func TestUpdateListPermissionsForInactiveUser(t *testing.T) {
	orga, user2 := testSetupUserToUpdatePermissions(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user2, lang, false)
	tx, _ := model.GetConnection().Beginx()
	err := updateObjectPermissionsForInactiveUser(tx, orga, user2.ID, true)

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatalf("List must have been transfered, but was: %v", err)
	}

	member := len(model.FindArticleListMemberByArticleListId(list.ID))

	if member != 2 {
		t.Fatalf("List must have only two members, but was: %v", member)
	}
}

func TestUpdateGroupPermissionsForInactiveUser(t *testing.T) {
	orga, user2 := testSetupUserToUpdatePermissions(t)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user2, true)
	tx, _ := model.GetConnection().Beginx()
	err := updateObjectPermissionsForInactiveUser(tx, orga, user2.ID, true)

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatalf("Group must have been transfered, but was: %v", err)
	}

	member := len(model.FindUserGroupMemberOnlyByUserGroupIdTx(nil, group.ID))

	if member != 1 {
		t.Fatalf("Group must have only one member, but was: %v", member)
	}
}

func testSetupUserToUpdatePermissions(t *testing.T) (*model.Organization, *model.User) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "remove@me.com")
	user2.OrganizationMember.IsAdmin = true
	user2.OrganizationMember.IsModerator = true

	if err := model.SaveOrganizationMember(nil, user2.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	adminGroup := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupAdminName)
	testutil.CreateUserGroupMember(t, adminGroup, user2, false)
	modGroup := model.GetUserGroupByOrganizationIdAndName(orga.ID, constants.GroupModName)
	testutil.CreateUserGroupMember(t, modGroup, user2, false)
	return orga, user2
}

func testSetArticlePrivate(t *testing.T, article *model.Article) {
	article.Private = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}
}
