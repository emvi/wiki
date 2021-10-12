package articlelist

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
)

func TestReadPrivateArticleLists(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test2@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateArticleList(t, orga, user, lang, false)
	withMember, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	testutil.CreateArticleListMember(t, withMember, user2.ID, 0, false)
	lists := ReadPrivateArticleLists(orga, user.ID, 0)

	if len(lists) != 1 {
		t.Fatalf("Expected one list to be found, but was: %v", len(lists))
	}
}

func TestReadArticleListNotFound(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	result, _, _, _, err := ReadArticleList(context.NewEmviUserContext(orga, user.ID), 0, 9999)

	if result != nil || err != errs.ArticleListNotFound {
		t.Fatal("Article list must not be found")
	}
}

func TestReadArticleList(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	testutil.CreateUser(t, orga, 321, "member@user.com")
	list, _ := testutil.CreateArticleList(t, orga, user, langEn, false)
	article0 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article1 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article2 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleListEntry(t, list, article0, 2)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 3)
	testutil.CreateArticleListName(t, list, langDe, "de", "")
	result, isMod, _, _, err := ReadArticleList(context.NewEmviUserContext(orga, user.ID), langEn.ID, list.ID)

	if result == nil || err != nil {
		t.Fatal("Article list must be found")
	}

	// default name and info...
	if result.Name == nil || result.Name.Name != "article list name" || result.Name.Info.String != "article list info" {
		t.Fatalf("Article list name not as expected, was: '%v' and '%v'", result.Name.Name, result.Name.Info.String)
	}

	if !isMod {
		t.Fatal("User must be moderator of list")
	}

	testHasName(t, result.Names, "article list name", "article list info")
	testHasName(t, result.Names, "de", "")
}

func TestReadArticleListPermissions(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	_, _, _, _, err := ReadArticleList(context.NewEmviUserContext(orga, user2.ID), lang.ID, list.ID)

	if err != errs.PermissionDenied {
		t.Fatalf("Permission to article list must be denied, but was: %v", err)
	}

	_, _, _, _, err = ReadArticleList(context.NewEmviUserContext(orga, user.ID), lang.ID, list.ID)

	if err != nil {
		t.Fatalf("Permission to article list must not be denied, but was: %v", err)
	}
}

func TestReadArticleListFallbackName(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	list, _ := testutil.CreateArticleList(t, orga, user, langDe, false)
	list, _, _, _, err := ReadArticleList(context.NewEmviUserContext(orga, user.ID), langDe.ID, list.ID)

	if err != nil {
		t.Fatalf("Article list must be found, but was: %v", err)
	}

	if len(list.Names) != 1 {
		t.Fatalf("Article list must have one name, but was: %v", len(list.Names))
	}

	if list.Name == nil || list.Name.LanguageId != langDe.ID {
		t.Fatal("Article list name must fallback to available language")
	}
}

func TestReadArticleListClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	ctx := context.NewEmviContext(orga, 0, nil, false)
	_, _, _, _, err := ReadArticleList(ctx, lang.ID, list.ID)

	// not found because reading the list checks for client access
	if err != errs.ArticleListNotFound {
		t.Fatalf("List must not be found, but was: %v", err)
	}

	testutil.SetListClientAccess(t, list)
	_, _, _, _, err = ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatalf("List must be returned, but was: %v", err)
	}
}

func TestReadArticleListEntries(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	list, _ := testutil.CreateArticleList(t, orga, user, langEn, false)
	article0 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article1 := testutil.CreateArticle(t, orga, user, langDe, true, true)
	article2 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article3 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article3.Archived = null.NewString("archived", true)

	if err := model.SaveArticle(nil, article3); err != nil {
		t.Fatal(err)
	}

	testutil.CreateArticleListEntry(t, list, article0, 2)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 3)
	testutil.CreateArticleListEntry(t, list, article3, 4)
	testutil.CreateArticleContent(t, user, article0, langDe, 0)
	_, _, _, err := ReadArticleListEntries(context.NewEmviUserContext(orga, user.ID), langEn.ID, list.ID+1, nil)

	if err != errs.ArticleListNotFound {
		t.Fatal("List must not be found")
	}

	entries, count, pos, err := ReadArticleListEntries(context.NewEmviUserContext(orga, user.ID), langDe.ID, list.ID, nil)

	if err != nil {
		t.Fatalf("List entries must be returned, but was: %v", err)
	}

	if len(entries) != 4 || count != 4 {
		t.Fatalf("Four list entries must be returned, but was: %v %v", len(entries), count)
	}

	if pos != 1 {
		t.Fatalf("Pos must be 0, but was: %v", pos)
	}

	// archived only
	entries, count, pos, err = ReadArticleListEntries(context.NewEmviUserContext(orga, user.ID), langDe.ID, list.ID, &model.SearchArticleListEntryFilter{Archived: true})

	if err != nil {
		t.Fatalf("List entries must be returned, but was: %v", err)
	}

	if len(entries) != 1 || count != 1 {
		t.Fatalf("One list entries must be returned, but was: %v %v", len(entries), count)
	}

	if pos != 1 {
		t.Fatalf("Pos must be 0, but was: %v", pos)
	}
}

func TestReadArticleListEntriesClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	article1 := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 1)
	testutil.SetArticleClientAccess(t, article2)
	ctx := context.NewEmviContext(orga, 0, nil, false)
	_, _, _, err := ReadArticleListEntries(ctx, lang.ID, list.ID, nil)

	if err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	testutil.SetListClientAccess(t, list)
	entries, count, pos, err := ReadArticleListEntries(ctx, lang.ID, list.ID, nil)

	if err != nil || len(entries) != 1 || count != 1 {
		t.Fatalf("List entries must be returned, but was: %v %v %v", err, len(entries), count)
	}

	if pos != 1 {
		t.Fatalf("Pos must be 0, but was: %v", pos)
	}
}

func TestReadArticleListEntriesPermissions(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	list, _ := testutil.CreateArticleList(t, orga, user, langEn, false)
	article0 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article1 := testutil.CreateArticle(t, orga, user, langDe, true, true)
	article2 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleListEntry(t, list, article0, 2)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 3)
	testutil.CreateArticleContent(t, user, article0, langDe, 0)
	_, _, _, err := ReadArticleListEntries(context.NewEmviUserContext(orga, user2.ID), langEn.ID, list.ID, nil)

	if err != errs.PermissionDenied {
		t.Fatalf("Permission to article list entries must be denied, but was: %v", err)
	}

	entries, count, pos, err := ReadArticleListEntries(context.NewEmviUserContext(orga, user.ID), langEn.ID, list.ID, nil)

	if err != nil {
		t.Fatalf("Permission to article list entries must not be denied, but was: %v", err)
	}

	if len(entries) != 3 || count != 3 {
		t.Fatalf("Three list entries must be returned, but was: %v %v", len(entries), count)
	}

	if pos != 1 {
		t.Fatalf("Pos must be 0, but was: %v", pos)
	}
}

func TestReadArticleListEntriesCenterArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, langEn, false)
	article0 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article1 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article2 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article3 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article4 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article5 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleListEntry(t, list, article0, 1)
	testutil.CreateArticleListEntry(t, list, article1, 2) // -1
	testutil.CreateArticleListEntry(t, list, article2, 3) // center
	testutil.CreateArticleListEntry(t, list, article3, 4) // +1
	testutil.CreateArticleListEntry(t, list, article4, 5) // +2
	testutil.CreateArticleListEntry(t, list, article5, 6)
	filter := &model.SearchArticleListEntryFilter{CenterArticleId: article2.ID, CenterBefore: 1, BaseSearch: model.BaseSearch{Limit: 4}}
	results, count, pos, err := ReadArticleListEntries(context.NewEmviUserContext(orga, user.ID), langEn.ID, list.ID, filter)

	if err != nil {
		t.Fatalf("Results must be returned, but was: %v", err)
	}

	if count != 6 || len(results) != 4 {
		t.Fatalf("Must return four results of six, but was: %v %v", count, len(results))
	}

	if results[1].ID != article2.ID {
		t.Fatal("Third article must be centered")
	}

	if pos != 2 {
		t.Fatalf("Pos must be 2, but was: %v", pos)
	}
}

func TestReadArticleListMemberUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)

	if member, count, err := ReadArticleListMember(orga, user.ID, list.ID, nil); err != nil || len(member) != 1 || count != 1 {
		t.Fatalf("Article list must have 1 members, but was %v and %v %v member", err, len(member), count)
	}
}

func TestReadArticleListMemberGroups(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	testutil.CreateArticleListMember(t, list, 0, group.ID, false)
	member, count, err := ReadArticleListMember(orga, user.ID, list.ID, nil)

	if err != nil || len(member) != 2 || count != 2 {
		t.Fatalf("Article list must have 2 members, but was %v and %v %v member", err, len(member), count)
	}

	if member[1].UserGroupId != group.ID {
		t.Fatal("Second user must be a group")
	}

	if member[1].UserGroup.MemberCount != 1 {
		t.Fatalf("Returned group must have one member, but was: %v", member[1].UserGroup.MemberCount)
	}
}

func TestReadArticleListMemberPermissions(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)

	if _, _, err := ReadArticleListMember(orga, user2.ID, list.ID, nil); err != errs.PermissionDenied {
		t.Fatalf("Permission to article list member must be denied, but was: %v", err)
	}

	if member, count, err := ReadArticleListMember(orga, user.ID, list.ID, nil); err != nil || len(member) != 1 || count != 1 {
		t.Fatalf("Permission to article list member must not be denied, but was: %v %v %v", err, len(member), count)
	}
}

func TestReadArticleListArticleCountAccessUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	article1 := testutil.CreateArticle(t, orga, user, lang, false, false)
	article2 := testutil.CreateArticle(t, orga, user2, lang, false, false)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 2)
	ctx := context.NewEmviUserContext(orga, user.ID)
	list, _, _, _, err := ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatal(err)
	}

	if list.ArticleCount != 1 {
		t.Fatalf("Article count must be 1, but was: %v", list.ArticleCount)
	}

	testutil.CreateArticleAccess(t, article2, user, nil, false)
	list, _, _, _, err = ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatal(err)
	}

	if list.ArticleCount != 2 {
		t.Fatalf("Article count must be 2, but was: %v", list.ArticleCount)
	}
}

func TestReadArticleListArticleCountAccessUserThroughGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	article1 := testutil.CreateArticle(t, orga, user, lang, false, false)
	article2 := testutil.CreateArticle(t, orga, user2, lang, false, false)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 2)
	group := testutil.CreateUserGroup(t, orga, "name")
	testutil.CreateUserGroupMember(t, group, user, false)
	ctx := context.NewEmviUserContext(orga, user.ID)
	list, _, _, _, err := ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatal(err)
	}

	if list.ArticleCount != 1 {
		t.Fatalf("Article count must be 1, but was: %v", list.ArticleCount)
	}

	testutil.CreateArticleAccess(t, article2, nil, group, false)
	list, _, _, _, err = ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatal(err)
	}

	if list.ArticleCount != 2 {
		t.Fatalf("Article count must be 2, but was: %v", list.ArticleCount)
	}
}

func TestReadArticleListArticleCountAccessClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "member@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, false)
	testutil.SetListClientAccess(t, list)
	article1 := testutil.CreateArticle(t, orga, user, lang, false, false)
	article2 := testutil.CreateArticle(t, orga, user2, lang, false, false)
	testutil.CreateArticleListEntry(t, list, article1, 1)
	testutil.CreateArticleListEntry(t, list, article2, 2)
	testutil.SetArticleClientAccess(t, article1)
	ctx := context.NewEmviContext(orga, 0, nil, false)
	list, _, _, _, err := ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatal(err)
	}

	if list.ArticleCount != 1 {
		t.Fatalf("Article count must be 1, but was: %v", list.ArticleCount)
	}

	testutil.SetArticleClientAccess(t, article2)
	list, _, _, _, err = ReadArticleList(ctx, lang.ID, list.ID)

	if err != nil {
		t.Fatal(err)
	}

	if list.ArticleCount != 2 {
		t.Fatalf("Article count must be 2, but was: %v", list.ArticleCount)
	}
}

func testHasName(t *testing.T, names []model.ArticleListName, name, info string) {
	for _, n := range names {
		if n.Name == name && n.Info.String == info {
			return
		}
	}

	t.Fatalf("Expected name '%v' and info '%v' not found in article list names", name, info)
}
