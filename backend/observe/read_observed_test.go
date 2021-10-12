package observe

import (
	"emviwiki/shared/testutil"
	"testing"
)

func TestReadObserved(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateObservedObject(t, user, article, nil, nil)
	testutil.CreateObservedObject(t, user, nil, list, nil)
	testutil.CreateObservedObject(t, user, nil, nil, group)

	a, l, g := ReadObserved(orga, user.ID, true, true, true, 0, 0, 0)

	if len(a) != 1 || len(l) != 1 || len(g) != 1 {
		t.Fatalf("Expected 1 article, 1 list and 1 group, but was: %v %v %v", len(a), len(l), len(g))
	}
}
