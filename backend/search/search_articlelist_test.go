package search

import (
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestSearchArticleList(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	lang2 := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	createTestArticleList(t, orga, user, lang, lang2, "name", "name1", "info", false)
	createTestArticleList(t, orga, user, lang, lang2, "second list", "second list1", "info", false)
	createTestArticleList(t, orga, nil, lang, lang2, "another list", "another list1", "no user access", false)
	createTestArticleList(t, orga, user, lang, lang2, "secondus listus", "secondus listus1", "should be found", false)
	createTestArticleList(t, orga, nil, lang, lang2, "second", "second1", "out of ideas", true)
	ctx := context.NewEmviUserContext(orga, user.ID)

	inout := []struct {
		Query string
		N     int
	}{
		{"", 3},
		{"xzuyxcf", 0},
		{"name", 1},
		{"second", 2},
		{"another", 0},
		{"user access", 0},
		{"should", 1},
		{"should be found", 1},
	}

	for _, io := range inout {
		if lists, _ := SearchArticleList(ctx, io.Query, nil); len(lists) != io.N {
			t.Fatalf("Expected %v lists to be found for '%v', but was: %v", io.N, io.Query, len(lists))
		}
	}
}

func TestSearchArticleListClientAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	lang2 := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	list, _ := createTestArticleList(t, orga, user, lang, lang2, "listname", "listname1", "info", true)
	testutil.SetListClientAccess(t, list)
	ctx := context.NewEmviUserContext(orga, user.ID)
	lists, count := SearchArticleList(ctx, "listname", nil)

	if len(lists) != 1 || count != 1 {
		t.Fatalf("One list must be found, but was: %v %v", len(lists), count)
	}

	if !lists[0].ClientAccess {
		t.Fatalf("List must have client access set to true")
	}
}

func TestSearchArticleListTranslationOrder(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	lang2 := testutil.CreateLang(t, orga, "de", "Deutsch", false)

	// set user to prefer non standard language
	user.Language.SetValid("de")

	if err := model.SaveUser(nil, user, false); err != nil {
		t.Fatal(err)
	}

	list1, _ := createTestArticleList(t, orga, user, lang, lang2, "Aaa", "Bbb", "Second result", true)
	list2, _ := createTestArticleList(t, orga, user, lang, lang2, "Bbb", "Aaa", "First result", true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	lists, count := SearchArticleList(ctx, "", &model.SearchArticleListFilter{SortName: "ASC"})

	if count != 2 {
		t.Fatalf("Two results must be returned, but was: %v", count)
	}

	if lists[0].ID != list2.ID || lists[1].ID != list1.ID {
		t.Fatalf("Lists must be sorted by translated name, but was: %v %v", lists[0].Name.Name, lists[1].Name.Name)
	}
}

func TestSearchArticleListArchivedArticles(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.Archived = null.NewString("archived", true)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateArticleListEntry(t, list, article, 1)
	ctx := context.NewEmviUserContext(orga, user.ID)
	lists, count := SearchArticleList(ctx, "", nil)

	if len(lists) != 1 || count != 1 {
		t.Fatalf("Expected one list to be returned, but was: %v %v", len(lists), count)
	}

	if lists[0].ArticleCount != 1 {
		t.Fatalf("List must show one article, but was: %v", lists[0].ArticleCount)
	}
}

func TestSearchArticleListUserGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateArticleList(t, orga, user, lang, true)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateArticleListMember(t, list, 0, group.ID, false)
	ctx := context.NewEmviUserContext(orga, user.ID)
	lists, count := SearchArticleList(ctx, "", &model.SearchArticleListFilter{UserGroupIds: []hide.ID{group.ID}})

	if len(lists) != 1 || count != 1 {
		t.Fatalf("Expected one list to be returned, but was: %v %v", len(lists), count)
	}

	if lists[0].ID != list.ID {
		t.Fatal("List with group access must have been returned")
	}
}

func TestSearchArticleListUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateArticleList(t, orga, user2, lang, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	lists, count := SearchArticleList(ctx, "", &model.SearchArticleListFilter{UserIds: []hide.ID{user.ID}})

	if len(lists) != 1 || count != 1 {
		t.Fatalf("Expected one list to be returned, but was: %v %v", len(lists), count)
	}

	if lists[0].ID != list.ID {
		t.Fatal("List with group access must have been returned")
	}
}

func createTestArticleList(t *testing.T, orga *model.Organization, user *model.User, lang, lang2 *model.Language, name, name2, info string, public bool) (*model.ArticleList, *model.ArticleListMember) {
	list := &model.ArticleList{OrganizationId: orga.ID, Public: public}

	if err := model.SaveArticleList(nil, list); err != nil {
		t.Fatal(err)
	}

	listname := &model.ArticleListName{ArticleListId: list.ID,
		LanguageId: lang.ID,
		Name:       name,
		Info:       null.NewString(info, info != "")}

	if err := model.SaveArticleListName(nil, listname); err != nil {
		t.Fatal(err)
	}

	// add additional user to test join on list names
	listname = &model.ArticleListName{ArticleListId: list.ID,
		LanguageId: lang2.ID,
		Name:       name2,
		Info:       null.NewString(info, info != "")}

	if err := model.SaveArticleListName(nil, listname); err != nil {
		t.Fatal(err)
	}

	var member *model.ArticleListMember

	if user != nil {
		member = &model.ArticleListMember{ArticleListId: list.ID,
			UserId:      user.ID,
			IsModerator: true}

		if err := model.SaveArticleListMember(nil, member); err != nil {
			t.Fatal(err)
		}

		// add additional user to test join on list member
		member = &model.ArticleListMember{ArticleListId: list.ID,
			UserId:      user.ID,
			IsModerator: false}

		if err := model.SaveArticleListMember(nil, member); err != nil {
			t.Fatal(err)
		}
	}

	return list, member
}
