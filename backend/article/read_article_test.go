package article

import (
	"emviwiki/backend/article/schema"
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"strings"
	"testing"
	"time"
)

func TestReadArticleNotFound(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)

	if _, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), 123, 321, 0, false, formatHTML); err != errs.ArticleNotFound {
		t.Fatal("Article must not be found")
	}
}

func TestReadArticleNoAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)

	if err := model.DeleteArticleAccessByArticleId(nil, article.ID); err != nil {
		t.Fatal(err)
	}

	if _, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, lang.ID, 0, false, formatHTML); err != errs.PermissionDenied {
		t.Fatal("Permission must be denied")
	}
}

func TestReadArticleContentDoesNotExist(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langRu := testutil.CreateLang(t, orga, "ru", "Russian", false)
	article := testutil.CreateArticle(t, orga, user, langRu, true, true)

	if result, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, langEn.ID, 0, false, formatHTML); err != nil || result.Content == nil {
		t.Fatal("Empty content must be returned if it does not exist")
	}

	if result, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, langEn.ID, 0, false, formatHTML); err != nil || result.Content == nil {
		t.Fatal("Empty content must be returned if it does not exist")
	}
}

func TestReadArticleGroupAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	//article := testutil.CreateArticle(t, orga, user, lang, false, false)
	article := testutil.CreateArticleWithoutContent(t, orga, user, lang, false, false)
	user2 := testutil.CreateUser(t, orga, hide.ID(567), "group@access.com")

	if _, err := ReadArticle(context.NewEmviUserContext(orga, user2.ID), article.ID, lang.ID, 0, false, formatHTML); err != errs.PermissionDenied {
		t.Fatal("Permission must be denied for user without access via group")
	}

	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateUserGroupMember(t, group, user2, false)
	testutil.CreateArticleAccess(t, article, nil, group, true)

	if _, err := ReadArticle(context.NewEmviUserContext(orga, user2.ID), article.ID, lang.ID, 0, false, formatHTML); err != nil {
		t.Fatal("User must have access via group")
	}
}

func TestReadArticle(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticleWithoutContent(t, orga, user, lang, false, false)
	doc := `{"type":"doc","content":[]}`
	content := testutil.CreateArticleContent(t, user, article, lang, 0)
	content.Content = doc

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	testutil.CreateArticleContentAuthor(t, user, content)
	testutil.CreateObservedObject(t, user, article, nil, nil)
	testutil.CreateBookmark(t, orga, user, article, nil)
	result, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, 0, 0, false, formatHTML)

	if err != nil || result.Article == nil || result.Content == nil {
		t.Fatal("Article and content must be returned")
	}

	if result.Article.Views != 55 || len(result.Article.Tags) != 4 {
		t.Fatalf("Article not as expected, was: %v", result.Article)
	}

	if result.Content.Title != "content title" ||
		result.Content.Content != doc ||
		result.Content.Version != 0 ||
		result.Content.Commit.String != "content commit" {
		t.Fatalf("Article content not as expected, was: %v", result.Content)
	}

	if !result.IsObserved {
		t.Fatal("Article must be observed")
	}

	if !result.IsBookmarked {
		t.Fatal("Article must not be bookmarked")
	}

	if !result.WriteAccess {
		t.Fatal("Must have write access")
	}

	if len(result.Authors) != 1 {
		t.Fatal("Article must have one author")
	}

	if result.Authors[0].OrganizationMember == nil {
		t.Fatal("Author must have the organization member joined")
	}
}

func TestReadArticleMarkdown(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticleWithoutContent(t, orga, user, lang, false, false)
	content := testutil.CreateArticleContent(t, user, article, lang, 0)
	content.Content = sampleArticleContent

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	ctx := context.NewEmviUserContext(orga, user.ID)
	result, err := ReadArticle(ctx, article.ID, lang.ID, 0, true, formatMarkdown)

	if err != nil {
		t.Fatalf("Article with markdown content must have been returned, but was: %v", err)
	}

	if result.Content.Content == "" {
		t.Fatal("Article contente must not be empty")
	}
}

func TestReadArticleClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateObservedObject(t, user, article, nil, nil)
	testutil.CreateBookmark(t, orga, user, article, nil)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, lang.ID, true)
	content.Content = sampleArticleContent

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	input := []struct {
		scopes       []string
		clientAccess bool
	}{
		{clientAccess: false},
		{[]string{client.Scopes["articles"].String()}, true},
		{[]string{client.Scopes["articles"].String(), client.Scopes["tags"].String()}, true},
		{[]string{client.Scopes["articles"].String(), client.Scopes["article_authors"].String()}, true},
		{[]string{client.Scopes["articles"].String(), client.Scopes["article_authors"].String(), client.Scopes["article_authors_mails"].String()}, true},
	}
	expected := []struct {
		err            error
		hasAuthors     bool
		hasAuthorMails bool
		hasTags        bool
	}{
		{errs.PermissionDenied, false, false, false},
		{nil, false, false, false},
		{nil, false, false, true},
		{nil, true, false, false},
		{nil, true, true, false},
	}

	for i, in := range input {
		if in.clientAccess {
			testutil.SetArticleClientAccess(t, article)
		}

		ctx := context.NewEmviContext(orga, 0, in.scopes, false)
		result, err := ReadArticle(ctx, article.ID, 0, 0, false, formatHTML)

		if err != expected[i].err {
			t.Fatalf("Expected '%v', but was: %v", expected[i].err, err)
		}

		if expected[i].err != nil {
			continue
		}

		if err != nil {
			t.Fatal("Error must be nil")
		}

		if result.Article == nil || result.Content == nil {
			t.Fatalf("Article and content must have been returned: %v %v", result.Article != nil, result.Content != nil)
		}

		if result.WriteAccess || result.IsObserved || result.IsBookmarked || result.Content.User != nil {
			t.Fatal("Write access, observed, bookmarked, content user must not be returned")
		}

		if !expected[i].hasAuthors && (len(result.Authors) != 0 || len(result.Content.Authors) != 0) {
			t.Fatalf("Authors must not be returned, but was: %v %v", len(result.Authors), len(result.Content.Authors))
		} else if expected[i].hasAuthors && (len(result.Authors) == 0 || len(result.Content.Authors) == 0) {
			t.Fatalf("Authors must be returned, but was: %v %v", len(result.Authors), len(result.Content.Authors))
		}

		if expected[i].hasAuthors && (result.Authors[0].OrganizationMember != nil || result.Authors[0].AcceptMarketing) {
			t.Fatalf("Authors must not have member or marketing set, but was: %v %v", result.Authors[0].OrganizationMember, result.Authors[0].AcceptMarketing)
		}

		if expected[i].hasAuthors {
			if !expected[i].hasAuthorMails && result.Authors[0].Email != "" {
				t.Fatal("Author must not have mail")
			} else if expected[i].hasAuthorMails && result.Authors[0].Email == "" {
				t.Fatal("Authors must have mail")
			}
		}

		if !expected[i].hasTags && len(result.Article.Tags) != 0 {
			t.Fatal("Tags must not be returned")
		} else if expected[i].hasTags && len(result.Article.Tags) == 0 {
			t.Fatal("Tags must be returned")
		}
	}
}

func TestReadArticleUserTagUsageCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticleWithoutContent(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticleWithoutContent(t, orga, user2, lang, false, false)
	testutil.SetArticleClientAccess(t, article)
	tag := testutil.CreateTag(t, orga, "usagecount")
	testutil.CreateArticleTag(t, article, tag)
	testutil.CreateArticleTag(t, article2, tag)
	ctx := context.NewEmviUserContext(orga, user.ID)
	result, err := ReadArticle(ctx, article.ID, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	testUsageCountTagFound(t, result, 1)
	testutil.CreateArticleAccess(t, article2, user, nil, false)
	result, err = ReadArticle(ctx, article.ID, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	testUsageCountTagFound(t, result, 2)
}

func TestReadArticleUserTagUsageCountThroughGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "user2@test.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticleWithoutContent(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticleWithoutContent(t, orga, user2, lang, false, false)
	testutil.SetArticleClientAccess(t, article)
	tag := testutil.CreateTag(t, orga, "usagecount")
	testutil.CreateArticleTag(t, article, tag)
	testutil.CreateArticleTag(t, article2, tag)
	group := testutil.CreateUserGroup(t, orga, "name")
	testutil.CreateUserGroupMember(t, group, user, false)
	ctx := context.NewEmviUserContext(orga, user.ID)
	result, err := ReadArticle(ctx, article.ID, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	testUsageCountTagFound(t, result, 1)
	testutil.CreateArticleAccess(t, article2, nil, group, false)
	result, err = ReadArticle(ctx, article.ID, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	testUsageCountTagFound(t, result, 2)
}

func TestReadArticleClientTagUsageCount(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, false, false)
	testutil.SetArticleClientAccess(t, article)
	tag := testutil.CreateTag(t, orga, "usagecount")
	testutil.CreateArticleTag(t, article, tag)
	testutil.CreateArticleTag(t, article2, tag)
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, lang.ID, true)
	content.Content = sampleArticleContent

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	ctx := context.NewEmviContext(orga, 0, []string{client.Scopes["tags"].String()}, false)
	result, err := ReadArticle(ctx, article.ID, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	testUsageCountTagFound(t, result, 1)
	testutil.SetArticleClientAccess(t, article2)
	result, err = ReadArticle(ctx, article.ID, lang.ID, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	testUsageCountTagFound(t, result, 2)
}

func TestReadArticleContentVersion(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user, article, lang := setupArticleHistory(t, -1)

	result, _ := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, lang.ID, 0, false, formatHTML)

	if result.Content == nil || !strings.Contains(result.Content.Content, "content 4") {
		t.Fatalf("Latest article content must be returned, but was: %v", result.Content)
	}

	result, _ = ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, lang.ID, 2, false, formatHTML)

	if result.Content == nil || !strings.Contains(result.Content.Content, "content 2") {
		t.Fatalf("Article content with version 2 must be returned, but was: %v", result.Content)
	}
}

func TestReadArticleContentVersionWIP(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user, article, lang := setupArticleHistory(t, 2)

	result, _ := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, lang.ID, 99999999, false, formatHTML)

	if result.Content == nil || !strings.Contains(result.Content.Content, "content 4") {
		t.Fatalf("WIP article content must be returned, but was: %v", result.Content)
	}
}

func TestReadArticleNonExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user, article, lang := setupArticleHistory(t, -1)
	orga.Expert = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	_, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, lang.ID, 1, false, formatHTML)

	if err != errs.RequiresExpertVersion {
		t.Fatalf("Article version must not be returned, but was: %v", err)
	}
}

func TestReadArticleWithoutLanguage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "ru", "Russian", true)
	article := testutil.CreateArticleWithoutContent(t, orga, user, lang, true, true)
	testutil.CreateArticleContent(t, user, article, lang, 0)
	result, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, 0, 0, false, formatHTML)

	if err != nil || result.Content == nil {
		t.Fatal("Content must be returned")
	}

	if result.Content.LanguageId != lang.ID {
		t.Fatalf("Content language must be %v, but was: %v", lang.ID, result.Content.LanguageId)
	}
}

func TestReadArticleWriteAccessReadEveryone(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test321@user.com")
	lang := testutil.CreateLang(t, orga, "ru", "Russian", true)
	article := testutil.CreateArticleWithoutContent(t, orga, user2, lang, true, false)
	result, err := ReadArticle(context.NewEmviUserContext(orga, user.ID), article.ID, 0, 0, false, formatHTML)

	if err != nil {
		t.Fatal(err)
	}

	if result.WriteAccess {
		t.Fatal("User must not have write access")
	}
}

func TestRenderArticleContent(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	content := &model.ArticleContent{Content: sampleArticleContent}
	ctx := context.NewEmviUserContext(orga, user.ID)
	out, err := renderArticleContent(ctx, orga.ID, user.ID, content, schema.HTMLSchema)

	if err != nil {
		t.Fatalf("Content must be rendered, but was: %v", err)
	}

	t.Log(out)
}

func TestFindAndReplaceMentions(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test2@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	articleNoAccess := testutil.CreateArticle(t, orga, user2, lang, false, false)
	group := testutil.CreateUserGroup(t, orga, "group")
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	listNoAccess, _ := testutil.CreateArticleList(t, orga, user2, lang, false)
	tag := testutil.CreateTag(t, orga, "tag")
	nullId := idToString(hide.ID(0))

	content := sampleArticleContentMentions
	content = strings.Replace(content, "ARTICLE_ID", idToString(article.ID), 1)
	content = strings.Replace(content, "ARTICLE_NOT_FOUND_ID", nullId, 1)
	content = strings.Replace(content, "ARTICLE_NO_ACCESS_ID", idToString(articleNoAccess.ID), 1)
	content = strings.Replace(content, "GROUP_ID", idToString(group.ID), 1)
	content = strings.Replace(content, "GROUP_NOT_FOUND_ID", nullId, 1)
	content = strings.Replace(content, "LIST_ID", idToString(list.ID), 1)
	content = strings.Replace(content, "LIST_NOT_FOUND_ID", nullId, 1)
	content = strings.Replace(content, "LIST_NO_ACCESS_ID", idToString(listNoAccess.ID), 1)
	content = strings.Replace(content, "TAG_NAME", tag.Name, 1)
	content = strings.Replace(content, "TAG_NOT_FOUND_NAME", "doesnotexist", 1)
	content = strings.Replace(content, "USER_NAME", user.OrganizationMember.Username, 1)
	content = strings.Replace(content, "USER_NOT_FOUND_NAME", "doesnotexist", 1)
	doc, err := prosemirror.ParseDoc(content)

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.NewEmviUserContext(orga, user.ID)
	findAndReplaceMentions(ctx, doc, orga.ID, user.ID, lang.ID)
	mentions := prosemirror.FindNodes(doc, -1, "mention")

	if len(mentions) != 12 {
		t.Fatalf("Unexpected number of mentions: %v", len(mentions))
	}

	testMentionTitle(t, mentions[0], article.LatestArticleContent.Title)
	testMentionTitle(t, mentions[1], mentionNotFoundEN)
	testMentionTitle(t, mentions[2], mentionNoAccessEN)
	testMentionTitle(t, mentions[3], group.Name)
	testMentionTitle(t, mentions[4], mentionNotFoundEN)
	testMentionTitle(t, mentions[5], list.Name.Name)
	testMentionTitle(t, mentions[6], mentionNotFoundEN)
	testMentionTitle(t, mentions[7], mentionNoAccessEN)
	testMentionTitle(t, mentions[8], tag.Name)
	testMentionTitle(t, mentions[9], mentionNotFoundEN)
	testMentionTitle(t, mentions[10], fmt.Sprintf("%v %v", user.Firstname, user.Lastname))
	testMentionTitle(t, mentions[11], mentionNotFoundEN)
}

func TestFindAndReplaceMentionsClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 321, "test2@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.SetArticleClientAccess(t, article)
	articleNoAccess := testutil.CreateArticle(t, orga, user2, lang, false, false)
	group := testutil.CreateUserGroup(t, orga, "group")
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.SetListClientAccess(t, list)
	listNoAccess, _ := testutil.CreateArticleList(t, orga, user2, lang, false)
	tag := testutil.CreateTag(t, orga, "tag")
	nullId := idToString(hide.ID(0))

	content := sampleArticleContentMentions
	content = strings.Replace(content, "ARTICLE_ID", idToString(article.ID), 1)
	content = strings.Replace(content, "ARTICLE_NOT_FOUND_ID", nullId, 1)
	content = strings.Replace(content, "ARTICLE_NO_ACCESS_ID", idToString(articleNoAccess.ID), 1)
	content = strings.Replace(content, "GROUP_ID", idToString(group.ID), 1)
	content = strings.Replace(content, "GROUP_NOT_FOUND_ID", nullId, 1)
	content = strings.Replace(content, "LIST_ID", idToString(list.ID), 1)
	content = strings.Replace(content, "LIST_NOT_FOUND_ID", nullId, 1)
	content = strings.Replace(content, "LIST_NO_ACCESS_ID", idToString(listNoAccess.ID), 1)
	content = strings.Replace(content, "TAG_NAME", tag.Name, 1)
	content = strings.Replace(content, "TAG_NOT_FOUND_NAME", "doesnotexist", 1)
	content = strings.Replace(content, "USER_NAME", user.OrganizationMember.Username, 1)
	content = strings.Replace(content, "USER_NOT_FOUND_NAME", "doesnotexist", 1)
	doc, err := prosemirror.ParseDoc(content)

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.NewEmviContext(orga, 0, []string{"articles:r", "lists:r", "tags:r"}, false)
	findAndReplaceMentions(ctx, doc, orga.ID, user.ID, lang.ID)
	mentions := prosemirror.FindNodes(doc, -1, "mention")

	if len(mentions) != 12 {
		t.Fatalf("Unexpected number of mentions: %v", len(mentions))
	}

	testMentionTitle(t, mentions[0], article.LatestArticleContent.Title)
	testMentionTitle(t, mentions[1], mentionNotFoundEN)
	testMentionTitle(t, mentions[2], mentionNoAccessEN)
	testMentionTitle(t, mentions[3], mentionNoAccessEN)
	testMentionTitle(t, mentions[4], mentionNoAccessEN)
	testMentionTitle(t, mentions[5], list.Name.Name)
	testMentionTitle(t, mentions[6], mentionNotFoundEN)
	testMentionTitle(t, mentions[7], mentionNoAccessEN)
	testMentionTitle(t, mentions[8], tag.Name)
	testMentionTitle(t, mentions[9], mentionNotFoundEN)
	testMentionTitle(t, mentions[10], mentionNoAccessEN)
	testMentionTitle(t, mentions[11], mentionNoAccessEN)
}

func TestUpdateArticleViews(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	views := article.Views

	updateArticleViews(article, user.ID)
	updateArticleViews(article, user.ID)
	article = model.GetArticleByOrganizationIdAndId(orga.ID, article.ID)

	if article.Views != views+1 {
		t.Fatalf("Article views must have been updated, but was: %v -> %v", views, article.Views)
	}

	visit := model.GetArticleVisitByArticleIdAndUserIdAndDefTimeAfter(article.ID, user.ID, time.Now().Add(-time.Second))

	if visit == nil {
		t.Fatal("Visit must exist")
	}
}

func TestUpdateArticleViewsClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.ClientAccess = true

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, lang.ID, true)
	content.Content = sampleArticleContent

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	views := article.Views
	ctx := context.NewEmviContext(orga, 0, []string{client.Scopes["articles"].String()}, false)
	_, err := ReadArticle(ctx, article.ID, 0, 0, false, formatHTML)

	if err != nil {
		t.Fatalf("Article must be returned, but was: %v", err)
	}

	article = model.GetArticleByOrganizationIdAndId(orga.ID, article.ID)

	if article.Views != views {
		t.Fatalf("Article views must not have been updated, but was: %v", article.Views)
	}
}

func TestReadArticleRecommendations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user2 := testutil.CreateUser(t, orga, 222, "test2@user.com")
	user3 := testutil.CreateUser(t, orga, 333, "test3@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleRecommendation(t, article, user, user3)
	testutil.CreateArticleRecommendation(t, article, user2, user3)
	recommendations := getRecommendations(article.ID, user3.ID)

	if len(recommendations) != 2 {
		t.Fatalf("Two recommendations must be returned, but was: %v", len(recommendations))
	}

	if recommendations[0].Member == nil ||
		recommendations[1].Member == nil ||
		recommendations[0].Member.User == nil ||
		recommendations[1].Member.User == nil {
		t.Fatalf("Recommendations must have user and member set, but was: %v", recommendations)
	}
}

func TestCheckVersionAllowed(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	expert := &model.Organization{Expert: true}
	nonExpert := &model.Organization{Expert: false}
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	content := &model.ArticleContent{
		ArticleId:  article.ID,
		LanguageId: lang.ID,
	}
	input := []struct {
		orga    *model.Organization
		content *model.ArticleContent
		version int
	}{
		{expert, content, 0},
		{nonExpert, content, 0},
	}
	expected := []error{
		nil,
		nil,
	}

	for i, in := range input {
		if err := checkVersionAllowed(in.orga, in.content, in.version); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
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
		Content:    `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"content 1"}]}]}`,
		Version:    1,
		Commit:     null.NewString("First commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 2",
		Content:    `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"content 2"}]}]}`,
		Version:    2,
		Commit:     null.NewString("Second commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 3",
		Content:    `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"content 3"}]}]}`,
		Version:    3,
		Commit:     null.NewString("Third commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 4",
		Content:    `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"content 4"}]}]}`,
		Version:    4,
		Commit:     null.NewString("Forth commit", true),
		ArticleId:  article.ID,
		LanguageId: lang.ID,
		UserId:     user.ID,
	}, {
		Title:      "title 4",
		Content:    `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"content 4"}]}]}`,
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

func testUsageCountTagFound(t *testing.T, result ArticleResult, usages int) {
	found := false

	for _, tag := range result.Article.Tags {
		if tag.Name == "usagecount" {
			found = true

			if tag.Usages != usages {
				t.Fatalf("Tag usage count must be %v, but was: %v", usages, tag.Usages)
			}
		}
	}

	if !found {
		t.Fatal("Tag must be found")
	}
}

func idToString(id hide.ID) string {
	idStr, _ := hide.ToString(id)
	return idStr
}

func testMentionTitle(t *testing.T, node prosemirror.Node, title string) {
	if node.Attrs[mentionTitleAttr] != title {
		t.Fatalf("Unexpected title '%v' expected '%v'", node.Attrs[mentionTitleAttr], title)
	}
}
