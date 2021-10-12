package article

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

func TestArchiveArticleMessage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := ArchiveArticle(orga, user.ID, article.ID, "", false); err != errs.MessageTooShort {
		t.Fatalf("Message must be set, but was: %v", err)
	}

	if err := ArchiveArticle(orga, user.ID, article.ID, "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891", false); err != errs.MessageTooLong {
		t.Fatalf("Message must be set, but was: %v", err)
	}
}

func TestArchiveArticleSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := ArchiveArticle(orga, user.ID, article.ID, "test", false); err != nil {
		t.Fatalf("Expected article to be archived, but was: %v", err)
	}

	article = model.GetArticleByOrganizationIdAndIdIgnoreArchived(orga.ID, article.ID)

	if article == nil || !article.Archived.Valid || article.Archived.String != "test" {
		t.Fatal("Article must be archived")
	}

	testutil.AssertFeedCreatedN(t, orga, "archived_article", 1)

	if err := ArchiveArticle(orga, user.ID, article.ID, "test", false); err != nil {
		t.Fatalf("Expected article to be de-archived, but was: %v", err)
	}

	article = model.GetArticleByOrganizationIdAndId(orga.ID, article.ID)

	if article == nil || article.Archived.Valid {
		t.Fatal("Article must not be archived")
	}

	testutil.AssertFeedCreatedN(t, orga, "restored_article", 1)
}

func TestArchiveArticleAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	user2 := testutil.CreateUser(t, orga, 456, "access@archived.com")

	input := []struct {
		UserId    hide.ID
		ArticleId hide.ID
	}{
		{user.ID, 0},
		{user2.ID, 0},
		{user2.ID, article.ID},
		{user.ID, article.ID},
	}
	expected := []error{
		errs.ArticleNotFound,
		errs.ArticleNotFound,
		errs.PermissionDenied,
		nil,
	}

	for i, in := range input {
		if err := ArchiveArticle(orga, in.UserId, in.ArticleId, "test", false); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}

func TestArchiveArticleAccessModAdmin(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	user2 := testutil.CreateUser(t, orga, 456, "access@archived.com")

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user2.ID)
	testSetModeratorOrAdmin(t, member, true, false)

	if err := ArchiveArticle(orga, user2.ID, article.ID, "test", false); err != nil {
		t.Fatalf("Moderator must have access to archive article, but was: %v", err)
	}

	testSetModeratorOrAdmin(t, member, false, true)

	if err := ArchiveArticle(orga, user2.ID, article.ID, "test", false); err != nil {
		t.Fatalf("Admin must have access to archive article, but was: %v", err)
	}

	testSetModeratorOrAdmin(t, member, false, false)
	access := testutil.CreateArticleAccess(t, article, user2, nil, false)

	if err := ArchiveArticle(orga, user2.ID, article.ID, "test", false); err != errs.PermissionDenied {
		t.Fatalf("User without write access must not have permission to archive article, but was: %v", err)
	}

	access.Write = true

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	if err := ArchiveArticle(orga, user2.ID, article.ID, "test", false); err != nil {
		t.Fatalf("User with write access must have permission to archive article, but was: %v", err)
	}
}

func TestArchiveArticleAccessReadEveryone(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	testSetModeratorOrAdmin(t, user.OrganizationMember, false, false)
	user2 := testutil.CreateUser(t, orga, 321, "test321@user.com")
	lang := testutil.CreateLang(t, orga, "ru", "Russian", true)
	article := testutil.CreateArticle(t, orga, user2, lang, true, false)

	if err := ArchiveArticle(orga, user.ID, article.ID, "test", false); err != errs.PermissionDenied {
		t.Fatalf("User must not have permission to archive article, but was: %v", err)
	}
}

func TestArchiveArticleDelete(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)

	if err := ArchiveArticle(orga, user.ID, article.ID, "", true); err != nil {
		t.Fatalf("Expected article to be archived and deleted, but was: %v", err)
	}

	if model.GetArticleByOrganizationIdAndIdIgnoreArchived(orga.ID, article.ID) != nil {
		t.Fatal("Article must be deleted")
	}

	testutil.AssertFeedCreated(t, orga, "delete_article")
}

func testSetModeratorOrAdmin(t *testing.T, member *model.OrganizationMember, mod, admin bool) {
	member.IsModerator = mod
	member.IsAdmin = admin

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		t.Fatal(err)
	}
}

func TestArchiveArticleRecommended(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleRecommendation(t, article, user, user2)

	if err := ArchiveArticle(orga, user.ID, article.ID, "", true); err != nil {
		t.Fatalf("Recommended article must have been deleted, but was: %v", err)
	}
}
