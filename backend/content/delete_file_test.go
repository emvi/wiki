package content

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/content"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	file := testutil.CreateFile(t, orga, user2, nil, "")
	input := []struct {
		orga       *model.Organization
		userId     hide.ID
		uniqueName string
	}{
		{nil, 0, ""},
		{orga, 0, ""},
		{orga, user.ID, ""},
		{orga, user.ID, "doesnotexist"},
		{orga, user.ID, file.UniqueName},
	}
	expected := []error{
		errs.PermissionDenied,
		errs.PermissionDenied,
		errs.PermissionDenied,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		if err := DeleteFile(in.orga, in.userId, in.uniqueName); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}

func TestDeleteFileNotAllowed(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	file := testutil.CreateFile(t, nil, user2, nil, "")

	// do not allow to delete files if the context is not set (user picture for example)
	if err := DeleteFile(nil, user.ID, file.UniqueName); err != errs.PermissionDenied {
		t.Fatal("User must not be allowed to delete file")
	}

	if err := DeleteFile(orga, user.ID, file.UniqueName); err != nil {
		t.Fatal("User must be allowed to delete file")
	}
}

func TestDeleteFileForArticleNilChannel(t *testing.T) {
	testutil.CleanBackendDb(t)
	testResetDummyStore(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	DeleteFileForArticle(nil, orga, user.ID, article.ID, nil)
	testFileDeleteCalls(t, 0)
}

func TestDeleteFileForArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	testResetDummyStore(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateFileWithMd5(t, orga, user, article, "", "a")
	testutil.CreateFileWithMd5(t, orga, user, article, "", "b")

	done := make(chan bool)
	DeleteFileForArticle(nil, orga, user.ID, article.ID, done)
	<-done
	close(done)
	testFileDeleteCalls(t, 2)
}

func TestDeleteFileForArticleIgnoreOtherArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	testResetDummyStore(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateFileWithMd5(t, orga, user, article, "", "a")
	testutil.CreateFileWithMd5(t, orga, user, article, "", "b")
	testutil.CreateFileWithMd5(t, orga, user, article2, "", "c")
	testutil.CreateFileWithMd5(t, orga, user, nil, "roomid", "d")

	done := make(chan bool)
	DeleteFileForArticle(nil, orga, user.ID, article.ID, done)
	<-done
	close(done)
	testFileDeleteCalls(t, 2)
}

func TestDeleteFileForArticleIgnoreFileInUseElsewhere(t *testing.T) {
	testutil.CleanBackendDb(t)
	testResetDummyStore(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateFileWithMd5(t, orga, user, article, "", "samemd5")
	testutil.CreateFileWithMd5(t, orga, user, article, "", "b")
	testutil.CreateFileWithMd5(t, orga, user, article2, "", "samemd5")
	testutil.CreateFileWithMd5(t, orga, user, nil, "roomid", "d")

	done := make(chan bool)
	DeleteFileForArticle(nil, orga, user.ID, article.ID, done)
	<-done
	close(done)
	testFileDeleteCalls(t, 1)
}

func testResetDummyStore(t *testing.T) {
	s, ok := store.(*content.DummyStore)

	if !ok {
		t.Fatal("Error getting dummy store")
	}

	s.Reset()
}

func testFileDeleteCalls(t *testing.T, n int) {
	s, ok := store.(*content.DummyStore)

	if !ok {
		t.Fatal("Error getting dummy store")
	}

	if len(s.Deletes) != n {
		t.Log(s.Deletes)
		t.Fatalf("Expected %v delete calls, but was: %v", n, len(s.Deletes))
	}
}
