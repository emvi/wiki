package perm

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"testing"
)

func TestCheckUserTagAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user2, lang, false, false)
	tag := testutil.CreateTag(t, orga, "userpermission")
	testutil.CreateArticleTag(t, article, tag)

	if err := CheckUserTagAccess(orga.ID, user.ID, tag.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	testutil.CreateArticleAccess(t, article, user, nil, false)

	if err := CheckUserTagAccess(orga.ID, user.ID, tag.ID); err != nil {
		t.Fatalf("Permission must be granted, but was: %v", err)
	}
}

func TestCheckUserTagAccessGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user2, lang, false, false)
	tag := testutil.CreateTag(t, orga, "userpermission")
	testutil.CreateArticleTag(t, article, tag)

	if err := CheckUserTagAccess(orga.ID, user.ID, tag.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user, false)
	testutil.CreateArticleAccess(t, article, nil, group, true)

	if err := CheckUserTagAccess(orga.ID, user.ID, tag.ID); err != nil {
		t.Fatalf("Permission must be granted, but was: %v", err)
	}
}

func TestCheckUserTagAccessPublic(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user2, lang, true, true)
	tag := testutil.CreateTag(t, orga, "userpermission")
	testutil.CreateArticleTag(t, article, tag)

	if err := CheckUserTagAccess(orga.ID, user.ID, tag.ID); err != nil {
		t.Fatalf("Permission must be granted, but was: %v", err)
	}
}

func TestCheckUserTagAccessSimple(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	tag := testutil.CreateTag(t, orga, "testtag")
	testutil.CreateArticleTag(t, article, tag)

	if err := CheckUserTagAccess(orga.ID, user.ID, tag.ID); err != nil {
		t.Fatalf("Permission must be granted, but was: %v", err)
	}
}
