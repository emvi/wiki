package schema

import (
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

const (
	testContentVersion1 = `{"type":"doc","content":[{"type":"image","attrs":{"src":"test.png"}}]}`
	testContentVersion2 = `{"type":"doc","content":[{"type":"image","attrs":{"src":"test.png"},"content":[{"type":"paragraph"}]}]}`
)

func TestMigrateNoMigration(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := testutil.CreateArticleContent(t, user, article, lang, 0)
	modTime := content.ModTime

	if err := Migrate(nil); err != nil {
		t.Fatalf("Migrate on nil must not return an error, but was: %v", err)
	}

	if err := Migrate(content); err != nil {
		t.Fatalf("Latest schema version must not be migrated, but was: %v", err)
	}

	if content.ModTime.UnixNano() != modTime.UnixNano() {
		t.Fatal("Article content must not have been updated")
	}
}

func TestMigrateEmptyContent(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := testutil.CreateArticleContent(t, user, article, lang, 0)
	content.SchemaVersion = 1
	content.Content = ""

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	if err := Migrate(content); err != nil {
		t.Fatalf("Empty article must have been migrated, but was: %v", err)
	}

	content = model.GetArticleContentById(content.ID)

	if content.Content != "" {
		t.Fatalf("Content must not have been changed, but was: %v", content.Content)
	}

	if content.SchemaVersion > 1 {
		t.Fatalf("Schema version must have been set, but was: %v", content.SchemaVersion)
	}
}

func TestMigrateVersion2(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := testutil.CreateArticleContent(t, user, article, lang, 0)
	content.SchemaVersion = 1
	content.Content = testContentVersion1

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	if err := Migrate(content); err != nil {
		t.Fatalf("Schema version 1 must have been migrated, but was: %v", err)
	}

	content = model.GetArticleContentById(content.ID)

	if content.Content != testContentVersion2 {
		t.Fatalf("Content must have been changed, but was: %v", content.Content)
	}

	if content.SchemaVersion != 2 {
		t.Fatalf("Schema version must have been set, but was: %v", content.SchemaVersion)
	}
}
