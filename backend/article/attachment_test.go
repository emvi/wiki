package article

import (
	"emviwiki/backend/content"
	content2 "emviwiki/shared/content"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
	"time"
)

func TestCleanupAttachments(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	con := testutil.CreateArticleContent(t, user, article, lang, 0)
	con.Content = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hier uploade ich einige Testdateien."}]},{"type":"image","attrs":{"src":"http://localhost:4003/api/content/7Irow97ZM6.jpg"}},{"type":"paragraph","content":[{"type":"file","attrs":{"file":"http://localhost:4003/api/content/dVrnCW3ZV9","name":"7Irow97ZM6.jpg","size":"49.28 MB"}}]},{"type":"pdf","attrs":{"src":"https://api.emvi/api/v1/content/keep.pdf"}}]}`

	if err := model.SaveArticleContent(nil, con); err != nil {
		t.Fatal(err)
	}

	testCreateAttachment(t, orga, user, article, lang.ID, "7Irow97ZM6.jpg", time.Now())
	testCreateAttachment(t, orga, user, article, lang.ID, "dVrnCW3ZV9", time.Now().Add(-time.Second*10))
	testCreateAttachment(t, orga, user, article, lang.ID, "keep.pdf", time.Now())
	testCreateAttachment(t, orga, user, article, lang.ID, "keepMeToo.pdf", time.Now().Add(-time.Second*10))
	testCreateAttachment(t, orga, user, article, lang.ID, "deleteMe", time.Now())
	testCreateAttachment(t, orga, user, article, lang.ID, "deleteMe.pdf", time.Now())
	testExistingFiles(t, orga.ID, 6)

	if err := cleanupAttachments(orga.ID, user.ID, article.ID, time.Now().Add(-time.Second*5), con); err != nil {
		t.Fatalf("Expected attachments to be cleaned up, but was: %v", err)
	}

	testExistingFiles(t, orga.ID, 4)
	con.Content = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hier uploade ich einige Testdateien."}]}]}`

	if err := model.SaveArticleContent(nil, con); err != nil {
		t.Fatal(err)
	}

	if err := cleanupAttachments(orga.ID, user.ID, article.ID, time.Now().Add(-time.Second*5), con); err != nil {
		t.Fatalf("Expected attachments to be cleaned up, but was: %v", err)
	}

	testExistingFiles(t, orga.ID, 2)
}

func TestCleanupAttachmentsDifferentLanguage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	con := testutil.CreateArticleContent(t, user, article, langEn, 0)
	con.Content = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hier uploade ich einige Testdateien."}]},{"type":"image","attrs":{"src":"http://localhost:4003/api/content/7Irow97ZM6.jpg"}},{"type":"paragraph","content":[{"type":"file","attrs":{"file":"http://localhost:4003/api/content/dVrnCW3ZV9","name":"7Irow97ZM6.jpg","size":"49.28 MB"}}]}]}`

	if err := model.SaveArticleContent(nil, con); err != nil {
		t.Fatal(err)
	}

	testCreateAttachment(t, orga, user, article, langEn.ID, "7Irow97ZM6.jpg", time.Now())
	testExistingFiles(t, orga.ID, 1)

	if err := cleanupAttachments(orga.ID, user.ID, article.ID, time.Now().Add(-time.Second*5), con); err != nil {
		t.Fatalf("Expected attachments to be cleaned up, but was: %v", err)
	}

	testExistingFiles(t, orga.ID, 1)

	// change language so that attachment must not be deleted, because it is used in the english content
	con.LanguageId = langDe.ID
	con.Content = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hier uploade ich einige Testdateien."}]}]}`

	if err := model.SaveArticleContent(nil, con); err != nil {
		t.Fatal(err)
	}

	if err := cleanupAttachments(orga.ID, user.ID, article.ID, time.Now().Add(-time.Second*5), con); err != nil {
		t.Fatalf("Expected attachments to be cleaned up, but was: %v", err)
	}

	testExistingFiles(t, orga.ID, 1)
}

func TestCleanupAttachmentsDifferentContentVersion(t *testing.T) {
	testutil.CleanBackendDb(t)
	store, _ := content.GetStore().(*content2.DummyStore)
	store.Reset()
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	con := testutil.CreateArticleContent(t, user, article, langEn, 0)
	con.Content = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Hier uploade ich einige Testdateien."}]}]}`

	if err := model.SaveArticleContent(nil, con); err != nil {
		t.Fatal(err)
	}

	testCreateAttachment(t, orga, user, article, langEn.ID, "filename.test", time.Now().Add(-time.Second*6)) // keep
	testCreateAttachment(t, orga, user, article, langEn.ID, "filename.test", time.Now())                     // delete
	testExistingFiles(t, orga.ID, 2)

	if err := cleanupAttachments(orga.ID, user.ID, article.ID, time.Now().Add(-time.Second*3), con); err != nil {
		t.Fatalf("Expected attachments to be cleaned up, but was: %v", err)
	}

	testExistingFiles(t, orga.ID, 1)

	if len(store.Deletes) != 0 {
		t.Fatalf("No files must have been deleted, but was: %v", len(store.Deletes))
	}
}

func testCreateAttachment(t *testing.T, orga *model.Organization, user *model.User, article *model.Article, langId hide.ID, uniqueName string, defTime time.Time) {
	file := testutil.CreateFile(t, orga, user, article, "")
	file.UniqueName = uniqueName
	file.LanguageId = langId
	file.Path = "test"

	if err := model.SaveFile(nil, file); err != nil {
		t.Fatal(err)
	}

	// modify def time
	if _, err := model.GetConnection().Exec(nil, `UPDATE "file" SET def_time = $2 WHERE id = $1`, file.ID, defTime); err != nil {
		t.Fatal(err)
	}
}

func testExistingFiles(t *testing.T, orgaId hide.ID, n int) {
	files := model.FindFileByOrganizationId(orgaId)

	if len(files) != n {
		t.Fatalf("%v files must exist, but was: %v", n, len(files))
	}
}
