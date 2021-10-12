package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/observe"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSaveArticleListNoPermission(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	if _, err := SaveArticleList(orga, user.ID+1, SaveArticleListData{list.ID, []SaveArticleListNameData{}, true, true}); err[0] != errs.PermissionDenied {
		t.Fatal("Permission must be denied")
	}
}

func TestSaveArticleListNameInvalid(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	names := []SaveArticleListNameData{{lang.ID, "", ""}}

	if _, err := SaveArticleList(orga, user.ID, SaveArticleListData{list.ID, names, true, true}); err[0].Error() != errs.NameTooShort.Error() {
		t.Fatalf("Name must be invalid, but was %v", err[0])
	}

	names[0].Name = "01234567890123456789012345678901234567891"

	if _, err := SaveArticleList(orga, user.ID, SaveArticleListData{list.ID, names, true, true}); err[0].Error() != errs.NameTooLong.Error() {
		t.Fatalf("Name must be invalid, but was %v", err[0])
	}

	names[0].Name = "name"
	names[0].Info = "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891"

	if _, err := SaveArticleList(orga, user.ID, SaveArticleListData{list.ID, names, true, true}); err[0].Error() != errs.InfoTooLong.Error() {
		t.Fatalf("Info must be invalid, but was %v", err[0])
	}

	names[0].Info = "info"
	names[0].LanguageId = 0

	if _, err := SaveArticleList(orga, user.ID, SaveArticleListData{list.ID, names, true, true}); err[0].Error() != errs.LanguageNotFound.Error() {
		t.Fatalf("Language must not be found, but was %v", err[0])
	}
}

func TestSaveArticleListCreateNew(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	names := []SaveArticleListNameData{{lang.ID, "name", "info"}, {langDe.ID, "name de", "info de"}}

	if id, err := SaveArticleList(orga, user.ID, SaveArticleListData{0, names, true, true}); err != nil || id == 0 {
		t.Fatal("New article list must have been created")
	}

	testutil.AssertFeedCreated(t, orga, "create_article_list")
	lists := model.FindArticleListsByOrganizationId(orga.ID)

	if len(lists) != 1 {
		t.Fatalf("One list must have been created, but was: %v", len(lists))
	}

	if !observe.IsObserved(user.ID, 0, lists[0].ID, 0) {
		t.Fatal("New list must be observed")
	}

	if !lists[0].Public || !lists[0].ClientAccess {
		t.Fatalf("List not as expected: %v", lists[0])
	}
}

func TestSaveArticleListUpdateExisting(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	names := []SaveArticleListNameData{{lang.ID, "name", "info"}, {langDe.ID, "name de", "info de"}}

	if id, err := SaveArticleList(orga, user.ID, SaveArticleListData{list.ID, names, true, true}); err != nil || id == 0 {
		t.Fatal("Article list must have been updated")
	}

	testutil.AssertFeedCreated(t, orga, "update_article_list")
	lists := model.FindArticleListsByOrganizationId(orga.ID)

	if len(lists) != 1 {
		t.Fatalf("One list must have been created, but was: %v", len(lists))
	}

	if observe.IsObserved(user.ID, 0, lists[0].ID, 0) {
		t.Fatal("Existing list must not be observed")
	}
}

func TestSaveArticleListOptionalNames(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	names := []SaveArticleListNameData{{lang.ID, "name", "info"}}

	if id, err := SaveArticleList(orga, user.ID, SaveArticleListData{0, names, true, true}); err != nil || id == 0 {
		t.Fatal("Article list must have been created")
	}

	testutil.AssertFeedCreated(t, orga, "create_article_list")
	lists := model.FindArticleListsByOrganizationId(orga.ID)

	if len(lists) != 1 {
		t.Fatalf("One list must be found, but was: %v", len(lists))
	}

	listNames := model.FindArticleListNamesByArticleListId(lists[0].ID)

	if len(listNames) != 1 {
		t.Fatalf("Article list must have one name, but was: %v", len(listNames))
	}

	names = []SaveArticleListNameData{{lang.ID, "name", "info"}, {langDe.ID, "name de", "info de"}}

	if id, err := SaveArticleList(orga, user.ID, SaveArticleListData{lists[0].ID, names, true, true}); err != nil || id == 0 {
		t.Fatal("Article list must have been updated")
	}

	testutil.AssertFeedCreated(t, orga, "update_article_list")
	listNames = model.FindArticleListNamesByArticleListId(lists[0].ID)

	if len(listNames) != 2 {
		t.Fatalf("Article list must have two names, but was: %v", len(listNames))
	}

	names = []SaveArticleListNameData{{lang.ID, "name", "info"}}

	if id, err := SaveArticleList(orga, user.ID, SaveArticleListData{lists[0].ID, names, true, true}); err != nil || id == 0 {
		t.Fatal("Article list must have been updated")
	}

	testutil.AssertFeedCreatedN(t, orga, "update_article_list", 2)
	listNames = model.FindArticleListNamesByArticleListId(lists[0].ID)

	if len(listNames) != 1 {
		t.Fatalf("Article list name must have been removed, but was: %v", len(listNames))
	}
}

func TestSaveArticleListMaxLists(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	lang := testutil.CreateLang(t, orga, "en", "English", true)
	names := []SaveArticleListNameData{{lang.ID, "name", "info"}}

	for i := 0; i < maxLists; i++ {
		testutil.CreateArticleList(t, orga, user, lang, true)
	}

	if _, err := SaveArticleList(orga, user.ID, SaveArticleListData{0, names, true, true}); err[0] != errs.MaxListsReached {
		t.Fatalf("Maximum number of lists must be reached, but was: %v", err)
	}
}
