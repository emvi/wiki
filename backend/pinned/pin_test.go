package pinned

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestPinObject(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)

	input := []struct {
		articleId hide.ID
		listId    hide.ID
	}{
		{0, 0},
		{article.ID, 0},
		{0, list.ID},
	}
	expected := []error{
		errs.NoObjectToPin,
		nil,
		nil,
	}

	for i, in := range input {
		if err := PinObject(orga, user.ID, in.articleId, in.listId); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}
