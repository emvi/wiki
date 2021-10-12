package article

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/null"
	"testing"
	"time"
)

const (
	emptySampleDoc     = `{"type":"doc","content":[{"type":"paragraph"}]}`
	nonEmptySampleDoc  = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Text"}]}]}`
	sampleDoc          = `{"type":"doc","content":[{"type":"image","attrs":{"src":"http://localhost:4003/api/v1/content/62RTK394jO.jpg"}},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Headline 1"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]},{"type":"headline","attrs":{"level":3},"content":[{"type":"text","text":"Headline 2"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]},{"type":"headline","attrs":{"level":3},"content":[{"type":"text","text":"Headline 3"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]}]}`
	sampleDocImageOnly = `{"type":"doc","content":[{"type":"image","attrs":{"src":"http://localhost:4003/api/v1/content/KmCMFEvoFi.png"}},{"type":"paragraph"}]}`
)

func TestGetArticlePreviewError(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user2, langEn, false, false)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = sampleDoc

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, false)

	if err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	if contentStr != "" {
		t.Fatalf("Content must be empty, but was: %v", contentStr)
	}
}

func TestGetArticlePreviewArticleNotFound(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = sampleDoc

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID+1, langEn.ID, false)

	if err != errs.ArticleNotFound {
		t.Fatalf("Article must not be found, but was: %v", err)
	}

	if contentStr != "" {
		t.Fatalf("Content must be empty, but was: %v", contentStr)
	}
}

func TestGetArticlePreviewParagraph(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = sampleDoc

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, true)

	if err != nil {
		t.Fatalf("Expected preview to be returned, but was: %v", err)
	}

	if contentStr != `<p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.</p>` {
		t.Fatalf("Preview content not as expected: %v", contentStr)
	}
}

func TestGetArticlePreviewEmptyUnpublishedArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := createEmptyUnpublishedArticle(t, orga, user, langEn)
	ctx := context.NewEmviUserContext(orga, user.ID)
	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, false)

	if err != nil {
		t.Fatalf("Expected preview to be returned, but was: %v", err)
	}

	if contentStr != `` {
		t.Fatalf("Preview content not as expected: %v", contentStr)
	}
}

func TestGetArticlePreviewEmptyParagraphFirst(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = `{"type":"doc","content":[{"type":"paragraph"},{"type":"paragraph","content":[{"type":"text","text":"Next paragraph!"}]}]}`

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, false)

	if err != nil {
		t.Fatalf("Expected preview to be returned, but was: %v", err)
	}

	if contentStr != `<p><br /></p><p>Next paragraph!</p>` {
		t.Fatalf("Preview content not as expected: %v", contentStr)
	}
}

func TestGetArticlePreviewImageOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = sampleDocImageOnly

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, false)

	if err != nil {
		t.Fatalf("Expected preview to be returned, but was: %v", err)
	}

	if contentStr != `<img src="http://localhost:4003/api/v1/content/KmCMFEvoFi.png" alt="http://localhost:4003/api/v1/content/KmCMFEvoFi.png" /><p><br /></p>` {
		t.Fatalf("Preview content not as expected: %v", contentStr)
	}
}

func TestGetArticlePreviewEmpty(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = emptySampleDoc

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, false)

	if err != nil {
		t.Fatalf("Expected preview to be returned, but was: %v", err)
	}

	if contentStr != "" {
		t.Fatalf("Preview content not as expected: %v", contentStr)
	}
}

func TestGetArticlePreviewParagraphOnly(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = nonEmptySampleDoc

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	contentStr, err := GetArticlePreview(ctx, article.ID, langEn.ID, false)

	if err != nil {
		t.Fatalf("Expected preview to be returned, but was: %v", err)
	}

	if contentStr != "<p>Text</p>" {
		t.Fatalf("Preview content not as expected: %v", contentStr)
	}
}

func createEmptyUnpublishedArticle(t *testing.T, orga *model.Organization, user *model.User, lang *model.Language) *model.Article {
	article := &model.Article{OrganizationId: orga.ID,
		Views:         54,
		ReadEveryone:  false,
		WriteEveryone: false,
		WIP:           1,
		Published:     null.NewTime(time.Now(), true)}

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	content := &model.ArticleContent{
		Title:           "Empty article",
		Content:         `{"type":"doc","content":[{"type":"paragraph"}]}`,
		TitleTsvector:   "empty article",
		ContentTsvector: "type doc content paragraph",
		Version:         1,
		WIP:             true,
		Commit:          null.NewString("WIP", true),
		ArticleId:       article.ID,
		LanguageId:      lang.ID,
		UserId:          user.ID,
		User:            user,
	}

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	access := &model.ArticleAccess{Write: true,
		UserId:    user.ID,
		ArticleId: article.ID}

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	article.Access = []model.ArticleAccess{*access}
	article.LatestArticleContent = content
	return article
}
