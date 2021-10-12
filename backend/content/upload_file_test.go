package content

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"strconv"
	"testing"
	"time"
)

func TestGetUploadDir(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)

	got := getUploadDir(orga, "some/path")
	expected := "some/path/test_" + strconv.Itoa(int(orga.ID)) + "/" + time.Now().Format(timeDirFormat)

	if got != expected {
		t.Fatalf("Expected %v as upload dir, but was: %v", expected, got)
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	if len(generateUniqueFilename()) != filenameLength {
		t.Fatalf("Generated string must be %v characters long", filenameLength)
	}
}

func TestCheckUploadLimitReached(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.MaxStorageGB = 1

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	for i := 0; i < 5; i++ {
		file := &model.File{UniqueName: "file1", OrganizationId: orga.ID, UserId: user.ID, Size: 1000000000, ArticleId: article.ID}

		if err := model.SaveFile(nil, file); err != nil {
			t.Fatal(err)
		}
	}

	if err := checkUploadLimitReached(orga); err != nil {
		t.Fatal("Upload limit must not be reached")
	}

	file := &model.File{UniqueName: "last_file", OrganizationId: orga.ID, UserId: user.ID, Size: 1000000001, ArticleId: article.ID}

	if err := model.SaveFile(nil, file); err != nil {
		t.Fatal(err)
	}

	if err := checkUploadLimitReached(orga); err != errs.MaxStorageReached {
		t.Fatal("Upload limit must be reached")
	}
}

func TestFindExistingFile(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateFileWithMd5(t, orga, user, nil, "", "file1")
	testutil.CreateFileWithMd5(t, orga, user, article, "", "file2")
	testutil.CreateFileWithMd5(t, orga, user, nil, "room", "file3")

	if findExistingFile(nil, orga, 0, "room", "file0") != nil {
		t.Fatal("File must not be found")
	}

	if findExistingFile(nil, orga, 0, "room", "file1") != nil {
		t.Fatal("File must not be found")
	}

	if findExistingFile(nil, orga, article.ID, "room", "file2") == nil {
		t.Fatal("File must be found")
	}

	if findExistingFile(nil, orga, 0, "room", "file3") == nil {
		t.Fatal("File must be found")
	}
}
