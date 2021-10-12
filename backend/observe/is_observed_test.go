package observe

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestIsObserved(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	group := testutil.CreateUserGroup(t, orga, "group")

	if IsObserved(user.ID, 0, 0, 0) {
		t.Fatal("No observed object must be found")
	}

	createTestObserved(t, user.ID, article.ID, 0, 0)
	createTestObserved(t, user.ID, 0, list.ID, 0)
	createTestObserved(t, user.ID, 0, 0, group.ID)

	if !IsObserved(user.ID, article.ID, 0, 0) {
		t.Fatal("Observed object must be found")
	}

	if !IsObserved(user.ID, 0, list.ID, 0) {
		t.Fatal("Observed object must be found")
	}

	if !IsObserved(user.ID, 0, 0, group.ID) {
		t.Fatal("Observed object must be found")
	}
}

func createTestObserved(t *testing.T, userId, articleId, listId, groupId hide.ID) *model.ObservedObject {
	observe := &model.ObservedObject{
		UserId:        userId,
		ArticleId:     articleId,
		ArticleListId: listId,
		UserGroupId:   groupId,
	}

	if err := model.SaveObservedObject(nil, observe); err != nil {
		t.Fatal(err)
	}

	return observe
}
