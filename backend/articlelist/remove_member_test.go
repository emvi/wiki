package articlelist

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"testing"

	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
)

func TestRemoveArticleListMemberForUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, userMember := testutil.CreateArticleList(t, orga, user, lang, true)
	user2 := testutil.CreateUser(t, orga, 9999, "test2@user.com")
	member := testutil.CreateArticleListMember(t, list, user2.ID, 0, false)
	var emptyIds []hide.ID

	in := []struct {
		userId        hide.ID
		listId        hide.ID
		removeUserIds []hide.ID
	}{
		{0, 0, emptyIds},
		{user.ID + 1, list.ID, emptyIds},
		{user.ID, list.ID, []hide.ID{userMember.ID}},
		{user.ID, list.ID, []hide.ID{member.ID}},
	}
	out := []struct {
		err error
	}{
		{errs.PermissionDenied},
		{errs.PermissionDenied},
		{errs.RemoveYourself},
		{nil},
	}

	for i := range in {
		if err := RemoveArticleListMember(orga, in[i].userId, in[i].listId, in[i].removeUserIds); err != out[i].err {
			t.Fatalf("Error '%v' does not match expected error '%v'", err, out[i].err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "remove_article_list_member")
}

func TestRemoveArticleListMemberForGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	group := testutil.CreateUserGroup(t, orga, "group")
	user2 := testutil.CreateUser(t, orga, 9999, "test2@user.com")
	user3 := testutil.CreateUser(t, orga, 9998, "test3@user.com")
	testutil.CreateUserGroupMember(t, group, user, false)
	testutil.CreateUserGroupMember(t, group, user2, true)
	testutil.CreateUserGroupMember(t, group, user3, false)
	member := testutil.CreateArticleListMember(t, list, 0, group.ID, false)

	if err := RemoveArticleListMember(orga, user.ID, list.ID, []hide.ID{member.ID}); err != nil {
		t.Fatalf("Article list member must have been removed, but was: %v", err)
	}

	testutil.AssertFeedCreatedN(t, orga, "remove_article_list_member", 1)
}

func TestRemoveArticleListMemberWhenGroupOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, userMember := testutil.CreateArticleList(t, orga, user, lang, true)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user, true)
	member := testutil.CreateArticleListMember(t, list, 0, group.ID, true)

	if err := model.DeleteArticleListMemberByArticleListIdAndId(nil, list.ID, userMember.ID); err != nil {
		t.Fatal(err)
	}

	if err := RemoveArticleListMember(orga, user.ID, list.ID, []hide.ID{member.ID}); err != errs.RemoveModeratorAccess {
		t.Fatalf("Article list member must not have been removed, but was: %v", err)
	}

	testutil.AssertFeedCreatedN(t, orga, "remove_article_list_member", 0)
}
