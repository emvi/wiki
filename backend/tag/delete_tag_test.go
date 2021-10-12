package tag

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestDeleteTag(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	tag := testutil.CreateTag(t, orga, "tag")
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleTag(t, article, tag)

	input := []struct {
		Orga   *model.Organization
		userId hide.ID
		TagId  hide.ID
	}{
		{orga, user.ID, 0},
		{orga, user.ID, tag.ID + 100},
		{orga, user2.ID, tag.ID},
		{orga, user.ID, tag.ID},
	}
	expected := []error{
		errs.PermissionDenied,
		errs.PermissionDenied,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		if err := DeleteTag(in.Orga, in.userId, in.TagId); err != expected[i] {
			t.Fatalf("Expected %v but was: %v", expected[i], err)
		}
	}

	feed := testutil.AssertFeedCreated(t, orga, "delete_tag")
	refs := model.FindFeedRefByFeedId(feed[0].ID)

	if len(refs) != 1 || refs[0].Key.String != "name" || refs[0].Value.String != "tag" {
		t.Fatalf("Feed must have deleted object name, but was: %v", refs)
	}
}
