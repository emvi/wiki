package history

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
)

func TestReadArticleHistory(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user, article, lang := setupArticleHistory(t, -1)

	input := []struct {
		ArticleId hide.ID
		LangId    hide.ID
		Offset    int
	}{
		{0, 0, 0},
		{article.ID, 0, 0},
		{article.ID, lang.ID, 0},
		{article.ID, lang.ID, 3},
	}
	expected := []struct {
		Error error
		Len   int
		Count int
	}{
		{errs.ArticleNotFound, 0, 0},
		{nil, 4, 4},
		{nil, 4, 4},
		{nil, 1, 4},
	}

	for i, in := range input {
		history, count, err := ReadArticleHistory(context.NewEmviUserContext(orga, user.ID), in.ArticleId, in.LangId, in.Offset)

		if err != expected[i].Error || len(history) != expected[i].Len || count != expected[i].Count {
			t.Fatalf("Expected %v and %v/%v items, but was: %v %v %v", expected[i].Error, expected[i].Len, expected[i].Count, err, len(history), count)
		}

		if expected[i].Error == nil && history[0].User == nil {
			t.Fatal("User must be set")
		}
	}
}

func TestReadArticleHistoryNonExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user, article, lang := setupArticleHistory(t, -1)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	history, count, err := ReadArticleHistory(context.NewEmviUserContext(orga, user.ID), article.ID, lang.ID, 0)

	if err != nil || len(history) != 3 || count != 4 {
		t.Fatalf("History must be returned and contain 3 items, but was: %v %v %v", err, len(history), count)
	}
}

func TestReadArticleHistoryClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _, article, lang := setupArticleHistory(t, -1)
	ctx := context.NewEmviContext(orga, 0, nil, false)
	history, count, err := ReadArticleHistory(ctx, article.ID, lang.ID, 0)

	if err != errs.PermissionDenied || len(history) != 0 || count != 0 {
		t.Fatalf("Permission must be denied, but was: %v %v %v", err, len(history), count)
	}

	testutil.SetArticleClientAccess(t, article)
	history, count, err = ReadArticleHistory(ctx, article.ID, lang.ID, 0)

	if err != nil || len(history) != 4 || count != 4 {
		t.Fatalf("History must be returned and contain 4 items, but was: %v %v %v", err, len(history), count)
	}

	for _, entry := range history {
		if entry.Authors != nil || entry.User != nil {
			t.Fatal("Authors and user must not be returned")
		}
	}
}

func setupArticleHistory(t *testing.T, wip int) (*model.Organization, *model.User, *model.Article, *model.Language) {
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)

	article := &model.Article{OrganizationId: orga.ID,
		Views:         54,
		ReadEveryone:  true,
		WriteEveryone: true,
		WIP:           wip}

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	content := []model.ArticleContent{{
		Title:      "title 1",
		Content:    "content 1",
		Version:    1,
		Commit:     null.NewString("First commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 2",
		Content:    "content 2",
		Version:    2,
		Commit:     null.NewString("Second commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 3",
		Content:    "content 3",
		Version:    3,
		Commit:     null.NewString("Third commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 4",
		Content:    "content 4",
		Version:    4,
		Commit:     null.NewString("Forth commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 4",
		Content:    "content 4",
		Version:    0,
		Commit:     null.NewString("Forth commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}}

	for _, c := range content {
		if err := model.SaveArticleContent(nil, &c); err != nil {
			t.Fatal(err)
		}

		author := &model.ArticleContentAuthor{ArticleContentId: c.ID,
			UserId: user.ID}

		if err := model.SaveArticleContentAuthor(nil, author); err != nil {
			t.Fatal(err)
		}
	}

	access := &model.ArticleAccess{Write: true,
		UserId:    user.ID,
		ArticleId: article.ID}

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	return orga, user, article, lang
}
