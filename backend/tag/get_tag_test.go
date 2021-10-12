package tag

import (
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestGetTagByIdOrName(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	tag := testutil.CreateTag(t, orga, "testtag")
	testutil.CreateArticleTag(t, article, tag)

	input := []struct {
		Id   hide.ID
		Name string
	}{
		{0, ""},
		{tag.ID + 1, ""},
		{tag.ID, ""},
		{0, "unknown"},
		{0, "testtag"},
	}
	expected := []string{
		"",
		"",
		"testtag",
		"",
		"testtag",
	}

	for i, in := range input {
		tag = GetTagByIdOrName(orga, user.ID, in.Id, in.Name)

		if expected[i] == "" && tag != nil {
			t.Fatalf("Expected no tag to be returned for: %v %v", in.Id, in.Name)
		} else if expected[i] != "" && (tag == nil || tag.Name != expected[i]) {
			t.Fatalf("Expected %v but was: %v", expected[i], tag)
		}
	}
}
