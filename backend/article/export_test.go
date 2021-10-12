package article

import (
	"emviwiki/backend/content"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/config"
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExportArticleFailure(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	ctx := context.NewEmviUserContext(orga, user.ID)

	if _, err := ExportArticle(ctx, article.ID, lang.ID, "invalid", true); err != errs.UnknownArticleFormat {
		t.Fatalf("Article format must be unknown, but was: %v", err)
	}

	if _, err := ExportArticle(ctx, article.ID+1, lang.ID, formatHTML, true); err != errs.ArticleNotFound {
		t.Fatalf("Article must not be found, but was: %v", err)
	}

	ctx = context.NewEmviUserContext(orga, user2.ID)

	if _, err := ExportArticle(ctx, article.ID, lang.ID, formatHTML, true); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}
}

func TestExportArticleUnpublished(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := &model.Article{
		OrganizationId: orga.ID,
		ReadEveryone:   true,
		WriteEveryone:  true,
		WIP:            1,
	}

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	version0 := testutil.CreateArticleContent(t, user, article, lang, 0)
	version0.Content = ""

	if err := model.SaveArticleContent(nil, version0); err != nil {
		t.Fatal(err)
	}

	testutil.CreateArticleContent(t, user, article, lang, 1)
	ctx := context.NewEmviUserContext(orga, user.ID)

	if _, err := ExportArticle(ctx, article.ID, lang.ID, formatHTML, false); err != errs.UnpublishedArticle {
		t.Fatalf("Article must be unpublished, but was: %v", err)
	}
}

func TestCreateExportZip(t *testing.T) {
	os.Setenv("EMVI_WIKI_STORE_TYPE", "file")
	config.Load()
	LoadConfig()
	content.LoadConfig()
	testutil.CleanBackendDb(t)

	if err := os.Mkdir("bucket", 0777); err != nil {
		t.Fatal(err)
	}

	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	c := &model.ArticleContent{ArticleId: article.ID, Title: "Test article", Content: "content"}
	file := testutil.CreateFile(t, orga, user, article, "")
	image := testutil.CreateFile(t, orga, user, article, "")
	pdf := testutil.CreateFile(t, orga, user, article, "")
	files := []model.File{*file, *image, *pdf}

	for _, f := range files {
		f.Path = ""

		if err := model.SaveFile(nil, &f); err != nil {
			t.Fatal(err)
		}

		if err := ioutil.WriteFile(filepath.Join("bucket", f.UniqueName), []byte(f.UniqueName), 0777); err != nil {
			t.Fatal(err)
		}
	}

	reader, writer := io.Pipe()
	go createExportZip(writer, c, zipExt, files)
	out, err := ioutil.ReadAll(reader)

	if err != nil {
		t.Fatalf("Zip must have been created, but was: %v", err)
	}

	if len(out) == 0 {
		t.Fatalf("Zip must have been filled, but was: %v", len(out))
	}
}

func TestFindFilesForFilesAndImages(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	file := testutil.CreateFile(t, orga, user, article, "")
	image := testutil.CreateFile(t, orga, user, article, "")
	pdf := testutil.CreateFile(t, orga, user, article, "")
	file.UniqueName = "j06rqfiflKwSRgmtw5li.txt"
	image.UniqueName = "hoADCzmBrtfmzM3LlzDJ.jpg"
	pdf.UniqueName = "BQyOlhF3RqEw9vHJNEBN.pdf"

	if err := model.SaveFile(nil, file); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveFile(nil, image); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveFile(nil, pdf); err != nil {
		t.Fatal(err)
	}

	doc, _ := prosemirror.ParseDoc(exportSampleDoc)
	files := findAndUpdateFilePathsForFilesAndImages(doc)

	if len(files) != 3 {
		t.Fatalf("Must have found three files, but was: %v", len(files))
	}

	for _, file := range files {
		if file.UniqueName != image.UniqueName &&
			file.UniqueName != file.UniqueName &&
			file.UniqueName != pdf.UniqueName {
			t.Fatalf("File '%v' must be in list of files", file.UniqueName)
		}
	}

	images := prosemirror.FindNodes(doc, 1, "image")

	if len(images) != 1 || images[0].Attrs["src"] != "files/original_name_hoADCzmBrtfmzM3LlzDJ.jpg" {
		t.Fatalf("Image src not as expected: %v", images)
	}
}

func TestRenderArticleTemplate(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	tag := testutil.CreateTag(t, orga, "tag")
	testutil.CreateArticleTag(t, article, tag)
	c := &model.ArticleContent{Title: "title", Content: exportSampleDoc, LanguageId: lang.ID, BaseEntity: db.BaseEntity{DefTime: time.Now()}}
	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := renderArticleTemplate(ctx, article, c); err != nil {
		t.Fatalf("Article content must have been rendered for export, but was: %v", err)
	}

	if !strings.Contains(c.Content, `<html lang="en"`) ||
		!strings.Contains(c.Content, `<title>title</title>`) ||
		!strings.Contains(c.Content, user.Firstname) ||
		!strings.Contains(c.Content, user.Lastname) ||
		!strings.Contains(c.Content, "tag") {
		t.Fatalf("Content not as expected: %v", c.Content)
	}
}
