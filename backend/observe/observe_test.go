package observe

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestObserveObjectError(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if err := ObserveObject(orga, user.ID, 0, 0, 0); err != errs.NoObjectToObserve {
		t.Fatal("There must have been no object to observe")
	}

	if err := ObserveObject(orga, user.ID, 123, 0, 0); err != errs.ArticleNotFound {
		t.Fatal("Article must not be found")
	}

	if err := ObserveObject(orga, user.ID, 0, 123, 0); err != errs.ArticleListNotFound {
		t.Fatal("Article list must not be found")
	}

	if err := ObserveObject(orga, user.ID, 0, 0, 123); err != errs.GroupNotFound {
		t.Fatal("User group must not be found")
	}
}

func TestObserveObjectSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := ObserveObject(orga, user.ID, article.ID, 0, 0); err != nil {
		t.Fatal("Observed object for article must have been saved")
	}
}

func TestObserveObjectToggle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := ObserveObject(orga, user.ID, article.ID, 0, 0); err != nil {
		t.Fatal("Observed object for article must have been saved")
	}

	if err := ObserveObject(orga, user.ID, article.ID, 0, 0); err != nil {
		t.Fatal("Observed object for article must have been removed")
	}

	if model.GetObservedObjectByUserIdAndArticleId(user.ID, article.ID) != nil {
		t.Fatal("Observed object must have been deleted")
	}
}
