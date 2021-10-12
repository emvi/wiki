package articlelist

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestToggleArticleListModerator(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	list, userMember := testutil.CreateArticleList(t, orga, user, lang, true)
	user2 := testutil.CreateUser(t, orga, 9999, "test2@user.com")
	member := testutil.CreateArticleListMember(t, list, user2.ID, 0, false)

	in := []struct {
		userId   hide.ID
		listId   hide.ID
		memberId hide.ID
	}{
		{0, 0, 0},
		{user.ID + 1, list.ID, member.ID},
		{user.ID, list.ID, userMember.ID},
		{user.ID, list.ID, member.ID + userMember.ID},
		{user.ID, list.ID, member.ID},
	}
	out := []struct {
		err error
	}{
		{errs.PermissionDenied},
		{errs.PermissionDenied},
		{errs.ModeratorYourself},
		{errs.MemberNotFound},
		{nil},
	}

	for i := range in {
		if err := ToggleArticleListModerator(orga, in[i].userId, in[i].listId, in[i].memberId); err != out[i].err {
			t.Fatalf("Error '%v' does not match expected error '%v'", err, out[i].err)
		}
	}

	testutil.AssertFeedCreated(t, orga, "set_article_list_moderator")

	if err := ToggleArticleListModerator(orga, user.ID, list.ID, member.ID); err != nil {
		t.Fatalf("Moderator must have been toggled, but was: %v", err)
	}

	testutil.AssertFeedCreated(t, orga, "remove_article_list_moderator")
}
