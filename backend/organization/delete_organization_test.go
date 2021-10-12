package organization

import (
	"emviwiki/backend/bookmark"
	"emviwiki/backend/errs"
	"emviwiki/backend/member"
	"emviwiki/backend/observe"
	"emviwiki/backend/pinned"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestDeleteOrganization(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@tester.com")
	user3 := testutil.CreateUserWithoutOrganization(t, 333, "user3@tester.com")
	user4 := testutil.CreateUserWithoutOrganization(t, 444, "user4@tester.com")

	// langs
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)

	// groups
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user, true)
	testutil.CreateUserGroupMember(t, group, user2, false)
	testObserveObject(t, orga, user, nil, nil, group)
	testObserveObject(t, orga, user2, nil, nil, group)

	// tags
	tag := testutil.CreateTag(t, orga, "testtag")

	// articles
	articlePublic := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleAccess(t, articlePublic, user, nil, true)
	testutil.CreateArticleAccess(t, articlePublic, user2, nil, false)
	testutil.CreateArticleAccess(t, articlePublic, nil, group, true)
	testutil.CreateArticleTag(t, articlePublic, tag)
	testutil.CreateArticleVisit(t, articlePublic, user)
	testutil.CreateArticleRecommendation(t, articlePublic, user, user2)
	testPinObject(t, orga, user, articlePublic, nil)
	testObserveObject(t, orga, user, articlePublic, nil, nil)
	testObserveObject(t, orga, user2, articlePublic, nil, nil)
	testBookmarkObject(t, orga, user, articlePublic, nil)
	testBookmarkObject(t, orga, user2, articlePublic, nil)

	articlePrivate := testutil.CreateArticle(t, orga, user2, langDe, false, false)
	testutil.CreateArticleAccess(t, articlePrivate, user, nil, false)
	testutil.CreateArticleAccess(t, articlePrivate, user2, nil, true)
	testutil.CreateArticleAccess(t, articlePrivate, nil, group, false)
	testutil.CreateArticleTag(t, articlePublic, tag)
	testObserveObject(t, orga, user, articlePrivate, nil, nil)
	testObserveObject(t, orga, user2, articlePrivate, nil, nil)
	testBookmarkObject(t, orga, user, articlePrivate, nil)
	testBookmarkObject(t, orga, user2, articlePrivate, nil)

	// article lists
	listPublic, _ := testutil.CreateArticleList(t, orga, user, langEn, true)
	testutil.CreateArticleListMember(t, listPublic, user2.ID, 0, false)
	testutil.CreateArticleListMember(t, listPublic, 0, group.ID, true)
	testutil.CreateArticleListEntry(t, listPublic, articlePublic, 1)
	testutil.CreateArticleListEntry(t, listPublic, articlePrivate, 2)
	testutil.CreateArticleListName(t, listPublic, langDe, "name", "info")
	testPinObject(t, orga, user, nil, listPublic)
	testObserveObject(t, orga, user, nil, listPublic, nil)
	testObserveObject(t, orga, user2, nil, listPublic, nil)
	testBookmarkObject(t, orga, user, nil, listPublic)
	testBookmarkObject(t, orga, user2, nil, listPublic)

	listPrivate, _ := testutil.CreateArticleList(t, orga, user2, langDe, false)
	testutil.CreateArticleListMember(t, listPublic, user.ID, 0, true)
	testutil.CreateArticleListMember(t, listPublic, 0, group.ID, false)
	testutil.CreateArticleListEntry(t, listPrivate, articlePublic, 1)
	testutil.CreateArticleListEntry(t, listPrivate, articlePrivate, 2)
	testutil.CreateArticleListName(t, listPrivate, langEn, "name", "info")
	testObserveObject(t, orga, user, nil, listPrivate, nil)
	testObserveObject(t, orga, user2, nil, listPrivate, nil)
	testBookmarkObject(t, orga, user, nil, listPrivate)
	testBookmarkObject(t, orga, user2, nil, listPrivate)

	// invitations
	if err := member.InviteMember(orga, user.ID, member.InviteMemberData{Emails: []string{user3.Email}, ReadOnly: false}, testutil.TestMailSender); err != nil {
		t.Fatal(err)
	}

	if err := member.InviteMember(orga, user.ID, member.InviteMemberData{Emails: []string{user4.Email}, ReadOnly: true}, testutil.TestMailSender); err != nil {
		t.Fatal(err)
	}

	// feed
	testutil.CreateFeedForObject(t, orga, user, articlePublic, nil, nil)
	testutil.CreateFeedForObject(t, orga, user, articlePrivate, nil, nil)
	testutil.CreateFeedForObject(t, orga, user, nil, listPublic, nil)
	testutil.CreateFeedForObject(t, orga, user, nil, listPrivate, nil)
	testutil.CreateFeedForObject(t, orga, user, nil, nil, group)

	// support tickets
	testutil.CreateSupportTicket(t, orga, user)

	// clients
	client := testutil.CreateClient(t, orga, "client", "id", "secret")
	testutil.CreateClientScope(t, client, "scope", true, false)
	testutil.CreateClientScope(t, client, "scope", true, true)

	// test execution
	authProvider := auth.NewMockAuthClient()
	input := []struct {
		OrgaId hide.ID
		UserId hide.ID
		Name   string
	}{
		{0, 0, ""},
		{orga.ID, 0, ""},
		{0, user.ID, ""},
		{orga.ID, user.ID, ""},
		{orga.ID, user.ID, orga.Name},
	}
	expected := []struct {
		Error error
	}{
		{errs.PermissionDenied},
		{errs.PermissionDenied},
		{errs.PermissionDenied},
		{errs.NameDoesNotMatch},
		{nil},
	}

	for i := range input {
		if err := DeleteOrganization(input[i].OrgaId, input[i].UserId, input[i].Name, authProvider); err != expected[i].Error {
			t.Fatalf("Expected error '%v' but was: %v", expected[i].Error, err)
		}
	}

	if authProvider.DeleteClientCalls != 1 {
		t.Fatalf("One client must have been deleted, but was: %v", authProvider.DeleteClientCalls)
	}
}

func testPinObject(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, list *model.ArticleList) {
	var articleId, listId hide.ID

	if article != nil {
		articleId = article.ID
	} else {
		listId = list.ID
	}

	if err := pinned.PinObject(orga, user.ID, articleId, listId); err != nil {
		t.Fatal(err)
	}
}

func testObserveObject(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, list *model.ArticleList, group *model.UserGroup) {
	var articleId, listId, groupId hide.ID

	if article != nil {
		articleId = article.ID
	} else if list != nil {
		listId = list.ID
	} else {
		groupId = group.ID
	}

	if err := observe.ObserveObject(orga, user.ID, articleId, listId, groupId); err != nil {
		t.Fatal(err)
	}
}

func testBookmarkObject(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, list *model.ArticleList) {
	var articleId, listId hide.ID

	if article != nil {
		articleId = article.ID
	} else {
		listId = list.ID
	}

	if err := bookmark.BookmarkObject(orga, user.ID, articleId, listId); err != nil {
		t.Fatal(err)
	}
}
