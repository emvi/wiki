package search

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSearchAll(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticleList(t, orga, user, lang, true)
	testutil.CreateUserGroup(t, orga, "testgroup")
	tag := testutil.CreateTag(t, orga, "tagname")
	testutil.CreateArticleTag(t, article, tag)
	testutil.CreateFeed(t, orga, user, lang, false)
	ctx := context.NewEmviUserContext(orga, user.ID)

	inOut := []struct {
		Query    string
		Articles int
		Lists    int
		User     int
		Groups   int
		Tags     int
	}{
		{"", 0, 0, 0, 0, 0},
		{"article", 2, 2, 0, 0, 2},
		{"article list", 2, 2, 0, 0, 2},
		{"testuser", 2, 0, 1, 1, 2},
		{"testgroup", 2, 0, 1, 2, 2},
		{"tagname", 1, 0, 1, 0, 1},
	}

	for _, io := range inOut {
		results := SearchAll(ctx, io.Query, nil)

		if len(results.Articles) != io.Articles ||
			len(results.Lists) != io.Lists ||
			len(results.User) != io.User ||
			len(results.Groups) != io.Groups ||
			len(results.Tags) != io.Tags {
			t.Fatalf("Expected search results for '%v' to contain %v articles, %v lists, %v user, %v groups, %v tags, but was %v, %v, %v, %v, %v", io.Query, io.Articles, io.Lists, io.User, io.Groups, io.Tags, len(results.Articles), len(results.Lists), len(results.User), len(results.Groups), len(results.Tags))
		}
	}
}

func TestSearchAllClientCall(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	article.ClientAccess = true
	list, _ := testutil.CreateArticleList(t, orga, user, lang, true)
	list.ClientAccess = true
	testutil.CreateUserGroup(t, orga, "testgroup")
	tag := testutil.CreateTag(t, orga, "tagname")
	testutil.CreateArticleTag(t, article, tag)
	testutil.CreateFeed(t, orga, user, lang, false)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveArticleList(nil, list); err != nil {
		t.Fatal(err)
	}

	inOut := []struct {
		Query    string
		Scopes   []string
		Articles int
		Lists    int
		User     int
		Groups   int
		Tags     int
	}{
		{"article", []string{}, 0, 0, 0, 0, 0},
		{"tagname", []string{}, 0, 0, 0, 0, 0},
		{"article list", []string{}, 0, 0, 0, 0, 0},
		{"article", []string{client.Scopes["articles"].String(), client.Scopes["search_articles"].String()}, 1, 0, 0, 0, 0},
		{"tagname", []string{client.Scopes["tags"].String(), client.Scopes["search_tags"].String()}, 0, 0, 0, 0, 1},
		{"article list", []string{client.Scopes["lists"].String(), client.Scopes["search_lists"].String()}, 0, 1, 0, 0, 0},
	}

	for _, io := range inOut {
		results := SearchAll(context.NewEmviContext(orga, 0, io.Scopes, false), io.Query, nil)

		if len(results.Articles) != io.Articles ||
			len(results.Lists) != io.Lists ||
			len(results.User) != io.User ||
			len(results.Groups) != io.Groups ||
			len(results.Tags) != io.Tags {
			t.Fatalf("Expected search results for '%v' to contain %v articles, %v lists, %v user, %v groups, %v tags, but was %v, %v, %v, %v, %v", io.Query, io.Articles, io.Lists, io.User, io.Groups, io.Tags, len(results.Articles), len(results.Lists), len(results.User), len(results.Groups), len(results.Tags))
		}
	}
}
