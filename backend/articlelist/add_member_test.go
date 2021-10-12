package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestAddArticleListMemberUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	user2 := testutil.CreateUser(t, orga, 9999, "test2@user.com")
	user3 := testutil.CreateUser(t, orga, 9998, "test3@user.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user3, false)
	var emptyIds []hide.ID

	input := []struct {
		UserId    hide.ID
		ListId    hide.ID
		MemberIds []hide.ID
		GroupIds  []hide.ID
	}{
		{user.ID, 0, emptyIds, emptyIds},
		{0, list.ID, emptyIds, emptyIds},
		{user.ID, list.ID, emptyIds, emptyIds},
		{user.ID, list.ID, []hide.ID{user.ID + 1}, emptyIds},
		{user.ID, list.ID, emptyIds, []hide.ID{123}},
		{user.ID, list.ID, []hide.ID{user2.ID}, emptyIds},
		{user.ID, list.ID, emptyIds, []hide.ID{group.ID}},
	}
	expected := []error{
		errs.ArticleListNotFound,
		errs.PermissionDenied,
		nil,
		errs.UserNotFound,
		errs.GroupNotFound,
		nil,
		nil,
	}

	for i, in := range input {
		if _, err := AddArticleListMember(orga, in.UserId, in.ListId, in.MemberIds, in.GroupIds); err != expected[i] {
			t.Fatalf("Expected %v, but was: %v", expected[i], err)
		}
	}

	testutil.AssertFeedCreatedN(t, orga, "add_article_list_member", 3)
}

func TestAddArticleListMemberGroups(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	user2 := testutil.CreateUser(t, orga, 9999, "test2@user.com")
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user2, false)
	userIds := []hide.ID{user.ID, user2.ID}
	groupIds := []hide.ID{group.ID}

	member := model.FindArticleListMemberByArticleListId(list.ID)

	if len(member) != 1 {
		t.Fatalf("Article list must have one member, but was: %v", len(member))
	}

	if _, err := AddArticleListMember(orga, user.ID, list.ID, userIds, groupIds); err != nil {
		t.Fatalf("Users and groups must be added to article list, but was: %v", err)
	}

	testutil.AssertFeedCreatedN(t, orga, "add_article_list_member", 1)
	member = model.FindArticleListMemberByArticleListId(list.ID)

	if len(member) != 3 {
		t.Fatalf("Article list must have three members, but was: %v", len(member))
	}

	if hide.ID(member[0].UserId) != user.ID && hide.ID(member[1].UserId) != user.ID && hide.ID(member[2].UserId) != user.ID {
		t.Fatal("User must be member of article list")
	}

	if hide.ID(member[0].UserId) != user2.ID && hide.ID(member[1].UserId) != user2.ID && hide.ID(member[2].UserId) != user2.ID {
		t.Fatal("User 2 must be member of article list")
	}

	if hide.ID(member[0].UserGroupId) != group.ID && hide.ID(member[1].UserGroupId) != group.ID && hide.ID(member[2].UserGroupId) != group.ID {
		t.Fatal("Group must be member of article list")
	}
}

func TestAddArticleListMemberNotification(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	user2 := testutil.CreateUser(t, orga, 9999, "test2@user.com")

	if _, err := AddArticleListMember(orga, user.ID, list.ID, []hide.ID{user2.ID}, nil); err != nil {
		t.Fatalf("Member must have been added, but was: %v", err)
	}

	feed := testutil.AssertFeedCreated(t, orga, "add_article_list_member")

	if model.GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotification(orga.ID, user2.ID, feed[0].ID, true) == nil {
		t.Fatal("User must have received a notification")
	}
}
