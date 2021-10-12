package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestRenameTag(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateTag(t, orga, "existing")
	tag := testutil.CreateTag(t, orga, "tag")
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleTag(t, article, tag)

	input := []struct {
		id     hide.ID
		userId hide.ID
		name   string
	}{
		{0, user.ID, "valid"},
		{tag.ID, user2.ID, "tagname"},
		{tag.ID, user.ID, ""},
		{tag.ID, user.ID, "0123456789012345678901234567890123456789012345678901234567891"},
		{tag.ID, user.ID, "invalid$"},
		{tag.ID, user.ID, "existing"},
		{tag.ID, user.ID, "tag"},
		{tag.ID, user.ID, "tagname"},
	}
	expected := []error{
		errs.PermissionDenied,
		errs.PermissionDenied,
		errs.TagEmpty,
		errs.TagLen,
		errs.TagInvalid,
		errs.TagNameExistsAlready,
		nil,
		nil,
	}

	for i, in := range input {
		if err := RenameTag(orga, in.userId, in.id, in.name); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}
