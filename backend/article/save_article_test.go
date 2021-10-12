package article

import (
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"testing"
)

const (
	simpleSampleDoc  = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"content"}]}]}`
	simpleSampleDoc2 = `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"new content"}]}]}`
)

func TestSaveArticleNotFound(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	data := SaveArticleData{Organization: orga, Id: 123}

	if _, err := SaveArticle(data); err[0] != errs.ArticleNotFound {
		t.Fatal("Article must not be found")
	}
}

func TestSaveArticleInvalidInput(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, _ := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{UserId: 1, LanguageId: lang.ID, Organization: orga}
	_, err := SaveArticle(data)

	if len(err) == 0 {
		t.Fatal("Input must contain errors")
	}

	if err[0] != errs.NoTitle {
		t.Fatal("Title must not be set")
	}

	data.Title = "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"
	data.CommitMsg = "12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901"
	_, err = SaveArticle(data)

	if len(err) == 0 {
		t.Fatal("Input must contain errors")
	}

	if err[0] != errs.TitleLen {
		t.Fatal("Title must be too long")
	}

	if err[1] != errs.CommitMsgLen {
		t.Fatal("Commit message must be too long")
	}
}

func TestSaveArticleNewSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	tag := &model.Tag{OrganizationId: orga.ID,
		Name: "foo"}

	if err := model.SaveTag(nil, tag); err != nil {
		t.Fatal(err)
	}

	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{Organization: orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "Commit message",
		Wip:           false,
		ReadEveryone:  false,
		WriteEveryone: false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new tag", "foo", "bar"},
		RTL:           true}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatal("Article must have been saved")
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article == nil {
		t.Fatal("Article must exist")
	}

	tags := model.FindTagByOrganizationIdAndUserIdAndArticleId(orga.ID, user.ID, article.ID)

	if article.OrganizationId == 0 ||
		article.WIP != -1 ||
		article.ReadEveryone ||
		article.WriteEveryone ||
		article.ClientAccess ||
		!article.Published.Valid ||
		len(tags) != 3 {
		t.Fatal("Article fields not as expected")
	}

	content := model.FindArticleContentByArticleIdAndLanguageId(id, lang.ID)

	if len(content) != 2 {
		t.Fatal("Two article contents must exist")
	}

	if content[0].ArticleId != id ||
		content[0].Title != "Title" ||
		content[0].Content != simpleSampleDoc ||
		content[0].Version != 0 ||
		content[0].Commit.String != "" ||
		content[0].SchemaVersion != constants.LatestSchemaVersion ||
		!content[0].RTL {
		t.Fatal("Article content 0 fields not as expected")
	}

	if content[1].ArticleId != id ||
		content[1].Title != "Title" ||
		content[1].Content != simpleSampleDoc ||
		content[1].Version != 1 ||
		content[1].Commit.String != "Commit message" ||
		content[1].SchemaVersion != constants.LatestSchemaVersion ||
		!content[1].RTL {
		t.Fatal("Article content 1 fields not as expected")
	}

	if model.GetTagByOrganizationIdAndName(orga.ID, "new tag") == nil {
		t.Fatal("Tag 'new tag' must have been created")
	}

	if model.GetTagByOrganizationIdAndName(orga.ID, "foo") == nil {
		t.Fatal("Tag 'foo' must have been created")
	}

	if model.GetTagByOrganizationIdAndName(orga.ID, "bar") == nil {
		t.Fatal("Tag 'bar' must have been created")
	}
}

func TestSaveArticleExistingSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{Organization: orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "Commit message",
		Wip:           false,
		ReadEveryone:  false,
		WriteEveryone: false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"},
		RTL:           true}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatal("Article must have been saved")
	}

	data.Id = id
	data.CommitMsg = "New commit message"
	data.Title = "Changed title"
	data.Content = simpleSampleDoc2
	data.ReadEveryone = true
	data.WriteEveryone = true
	data.ClientAccess = true
	data.Tags = []string{"bar", "foo", "duba"}
	data.RTL = false // test it gets updated
	id, errors = SaveArticle(data)

	if errors != nil {
		t.Fatal("Article must have been updated")
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article == nil {
		t.Fatal("Article must exist")
	}

	if article.OrganizationId == 0 ||
		article.WIP != -1 ||
		!article.ReadEveryone ||
		!article.WriteEveryone ||
		!article.ClientAccess {
		t.Fatal("Article fields not as expected")
	}

	tags := model.FindTagByOrganizationIdAndUserIdAndArticleId(orga.ID, user.ID, article.ID)

	if len(tags) != 4 {
		t.Fatalf("Tags must have been saved, but was: %v", article.Tags)
	}

	content := model.FindArticleContentByArticleIdAndLanguageId(id, lang.ID)

	if len(content) != 3 {
		t.Fatalf("Three article contents must exist, but was %v", len(content))
	}

	if content[0].ArticleId != id ||
		content[0].Title != "Changed title" ||
		content[0].Content != simpleSampleDoc2 ||
		content[0].Version != 0 ||
		content[0].Commit.String != "" ||
		content[0].SchemaVersion != constants.LatestSchemaVersion ||
		content[0].RTL {
		t.Fatalf("Article content 0 fields not as expected: %v", content[0])
	}

	if content[1].ArticleId != id ||
		content[1].Title != "Title" ||
		content[1].Content != simpleSampleDoc ||
		content[1].Version != 1 ||
		content[1].Commit.String != "Commit message" ||
		content[1].SchemaVersion != constants.LatestSchemaVersion ||
		!content[1].RTL {
		t.Fatalf("Article content 1 fields not as expected: %v", content[1])
	}

	if content[2].ArticleId != id ||
		content[2].Title != "Changed title" ||
		content[2].Content != simpleSampleDoc2 ||
		content[2].Version != 2 ||
		content[2].Commit.String != "New commit message" ||
		content[2].SchemaVersion != constants.LatestSchemaVersion ||
		content[2].RTL {
		t.Fatalf("Article content 2 fields not as expected: %v", content[2])
	}
}

func TestSaveArticleWIPSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	author := testutil.CreateUser(t, orga, 321, "author@testutil.com")
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{Organization: orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		Authors:       []hide.ID{author.ID},
		CommitMsg:     "Commit message",
		Wip:           false,
		ReadEveryone:  false,
		WriteEveryone: false,
		Private:       false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatal("Article must have been saved")
	}

	data.Id = id
	data.CommitMsg = "New commit message"
	data.Title = "Changed title"
	data.Content = simpleSampleDoc2
	data.Wip = true
	data.ReadEveryone = true
	data.WriteEveryone = true
	data.ClientAccess = true
	id, errors = SaveArticle(data)

	if errors != nil {
		t.Fatal("Article must have been updated")
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article == nil {
		t.Fatal("Article must exist")
	}

	if article.OrganizationId == 0 ||
		article.WIP != 1 ||
		!article.ReadEveryone ||
		!article.WriteEveryone ||
		!article.ClientAccess ||
		!article.Published.Valid {
		t.Fatal("Article fields not as expected")
	}

	content := model.FindArticleContentByArticleIdAndLanguageId(id, lang.ID)

	if len(content) != 3 {
		t.Fatalf("Three article contents must exist, but was %v", len(content))
	}

	if content[0].ArticleId != id ||
		content[0].Title != "Title" ||
		content[0].Content != simpleSampleDoc ||
		content[0].Version != 0 ||
		content[0].Commit.String != "" {
		t.Fatal("Article content 0 fields not as expected")
	}

	if content[1].ArticleId != id ||
		content[1].Title != "Title" ||
		content[1].Content != simpleSampleDoc ||
		content[1].Version != 1 ||
		content[1].Commit.String != "Commit message" {
		t.Fatal("Article content 1 fields not as expected")
	}

	if content[2].ArticleId != id ||
		content[2].Title != "Changed title" ||
		content[2].Content != simpleSampleDoc2 ||
		content[2].Version != 2 ||
		content[2].Commit.String != "New commit message" {
		t.Fatal("Article content 2 fields not as expected")
	}

	authors := model.FindArticleContentAuthorByArticleContentId(content[0].ID)

	if len(authors) != 0 {
		t.Fatal("No authors must be found for content version 0")
	}

	authors = model.FindArticleContentAuthorByArticleContentId(content[1].ID)

	if len(authors) != 2 {
		t.Fatal("Two authors must be found for content version 1")
	}

	authors = model.FindArticleContentAuthorByArticleContentId(content[2].ID)

	if len(authors) != 2 {
		t.Fatal("Two authors must be found for content version 2")
	}
}

func TestSaveArticleWIPFirstCommit(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	author := testutil.CreateUser(t, orga, 321, "author@testutil.com")
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{Organization: orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		Authors:       []hide.ID{author.ID},
		CommitMsg:     "Commit message",
		Wip:           true,
		ReadEveryone:  false,
		WriteEveryone: false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	result, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), id, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatalf("Article must have been found after creation, but was: %v", err)
	}

	if result.Article == nil {
		t.Fatalf("Article must be found, but was: %v", err)
	}

	if result.Article.WIP != 1 {
		t.Fatalf("Article must be WIP, but was: %v", result.Article.WIP)
	}

	if result.Article.Published.Valid {
		t.Fatalf("Article published date must not be set, but was: %v", result.Article.Published)
	}

	if result.Content == nil {
		t.Fatal("Article must have content version 0")
	}

	if result.Content.Content != "" {
		t.Fatalf("Content must not be filled, but was: %v", result.Content.Content)
	}
}

func TestSaveArticleChangeLanguage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	data := SaveArticleData{Id: article.ID,
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    langDe.ID,
		CommitMsg:     "Changed language",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	_, errors := SaveArticle(data)

	if errors != nil {
		t.Fatal("Article must have been saved")
	}

	content := model.FindArticleContentByArticleIdAndLanguageId(article.ID, langEn.ID)

	if len(content) != 3 {
		t.Fatalf("Expected 3 commits to exist, but was: %v", len(content))
	}

	content = model.FindArticleContentByArticleIdAndLanguageId(article.ID, langDe.ID)

	if len(content) != 2 {
		t.Fatalf("Expected 2 commits to exist, but was: %v", len(content))
	}

	if content[0].Version != 0 || content[1].Version != 1 {
		t.Fatalf("Expected commits to have propper versions, but was: %v %v", content[0].Version, content[1].Version)
	}
}

func TestSaveArticleUpdateAttachments(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	file := testutil.CreateFile(t, orga, user, nil, "roomid")
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "Update attachments",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	file = model.GetFileByOrganizationIdAndUniqueName(orga.ID, "unique_name0")

	if file == nil {
		t.Fatal("File must have been found")
	}

	if file.RoomId.Valid || file.ArticleId != id {
		t.Fatalf("File room ID must have been updated, but was: %v %v", file.RoomId, file.ArticleId)
	}
}

func TestSaveArticleKeepWIP(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "Create article",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	data.Id = id
	data.RoomId = ""
	data.CommitMsg = "first wip"
	data.Wip = true
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article first WIP commit must have been saved, but was: %v", errors)
	}

	data.CommitMsg = "second wip"
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article second WIP commit must have been saved, but was: %v", errors)
	}

	data.CommitMsg = "third wip"
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article third WIP commit must have been saved, but was: %v", errors)
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article.WIP != 1 {
		t.Fatalf("Article WIP version must be 1, but was: %v", article.WIP)
	}
}

func TestSaveArticlePublishAfterWIPAndDeleteWIP(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "Create article",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	data.Id = id
	data.RoomId = ""
	data.CommitMsg = "wip commit"
	data.Wip = true
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article WIP commit must have been saved, but was: %v", errors)
	}

	data.CommitMsg = "publish"
	data.Wip = false
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article.WIP != -1 {
		t.Fatalf("Article WIP version must be -1, but was: %v", article.WIP)
	}

	if !article.Published.Valid {
		t.Fatal("Article published date must have been set")
	}

	content := model.FindArticleContentByArticleId(id)

	if len(content) != 3 {
		t.Fatalf("Expected two contents to exist, but was: %v", len(content))
	}
}

func TestSaveArticlePublishAfterWIP(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "WIP save",
		Wip:           true,
		ReadEveryone:  true,
		WriteEveryone: true,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	data.Id = id
	data.Wip = false
	data.Access = []perm.SaveArticleAccess{{user.ID, 0, true}}
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been published, but was: %v", errors)
	}
}

func TestSaveArticlePrivate(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test2@user.com")
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	access := []perm.SaveArticleAccess{
		{UserId: user.ID, Write: true},
		{UserId: user2.ID, Write: true},
	}
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "save",
		Wip:           false,
		ReadEveryone:  false,
		WriteEveryone: false,
		Private:       false, // save access permissions
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"},
		Access:        access}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	savedAccess := model.FindArticleAccessByOrganizationIdAndArticleId(orga.ID, id)

	if len(savedAccess) != 2 {
		t.Fatalf("Two access permission must have been saved, but was: %v", len(savedAccess))
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article.Private {
		t.Fatal("Aricle must not be private")
	}

	data.Private = true // now save committing user only
	id, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	savedAccess = model.FindArticleAccessByOrganizationIdAndArticleId(orga.ID, id)

	if len(savedAccess) != 1 {
		t.Fatalf("One access permission must have been saved, but was: %v", len(savedAccess))
	}

	if savedAccess[0].UserId != data.UserId {
		t.Fatalf("The user who saved the article must have access to it when private, but was: %v", savedAccess[0])
	}

	article = model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if !article.Private {
		t.Fatal("Aricle must be private")
	}
}

func TestSaveArticlePermissionsChanged(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test2@user.com")
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	access := []perm.SaveArticleAccess{
		{UserId: user.ID, Write: true},
		{UserId: user2.ID, Write: true},
	}
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "save",
		Wip:           false,
		ReadEveryone:  false,
		WriteEveryone: false,
		Private:       false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"},
		Access:        access}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	data.Id = id
	data.RoomId = ""
	data.Private = true
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	data.UserId = user2.ID
	_, errors = SaveArticle(data)

	if len(errors) != 1 && errors[0] != errs.PermissionDenied {
		t.Fatalf("Access must be denied, but was: %v", errors)
	}
}

func TestSaveArticleRemovePinnedPrivate(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "save",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		Private:       false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)
	article.Pinned = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	data.Id = id
	data.RoomId = ""
	data.ReadEveryone = false
	data.WriteEveryone = false
	data.Private = true
	_, errors = SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	if model.GetArticleByOrganizationIdAndIdAndPinned(orga.ID, id) != nil {
		t.Fatal("Article must not be pinned after saving private")
	}
}

func TestSaveArticleClientAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	lang := testutil.CreateLang(t, orga, "de", "Deutsch", true)
	data := SaveArticleData{RoomId: "roomid",
		Organization:  orga,
		UserId:        user.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "save",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		Private:       false,
		ClientAccess:  true,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{"new_tag", "foo", "bar"}}
	id, errors := SaveArticle(data)

	if errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	article := model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if article.ClientAccess {
		t.Fatal("Client must not have access to article")
	}

	orga.Expert = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	data.Id = id

	if _, errors := SaveArticle(data); errors != nil {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}

	article = model.GetArticleByOrganizationIdAndId(orga.ID, id)

	if !article.ClientAccess {
		t.Fatal("Client must have access to article")
	}
}

func TestSaveArticleDeleteWIPFeedRefs(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := testutil.CreateArticleContent(t, user, article, lang, 3)
	content.WIP = true

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	feed := &model.Feed{Public: false,
		Reason:            "reason",
		OrganizationId:    orga.ID,
		TriggeredByUserId: user.ID}

	if err := model.SaveFeed(nil, feed); err != nil {
		t.Fatal(err)
	}

	ref := &model.FeedRef{FeedId: feed.ID,
		ArticleContentID: content.ID}

	if err := model.SaveFeedRef(nil, ref); err != nil {
		t.Fatal(err)
	}

	access := &model.FeedAccess{FeedId: feed.ID,
		UserId: user.ID}

	if err := model.SaveFeedAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	data := SaveArticleData{Organization: orga,
		UserId:        user.ID,
		Id:            article.ID,
		LanguageId:    lang.ID,
		CommitMsg:     "save",
		Wip:           false,
		ReadEveryone:  true,
		WriteEveryone: true,
		Private:       false,
		ClientAccess:  false,
		Title:         "Title",
		Content:       simpleSampleDoc,
		Tags:          []string{}}
	_, errors := SaveArticle(data)

	if len(errors) != 0 {
		t.Fatalf("Article must have been saved, but was: %v", errors)
	}
}
