package search

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"testing"
	"time"
)

func TestSearchArticleFiltered(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)

	inout := []struct {
		Query  string
		Filter *model.SearchArticleFilter
		N      int
	}{
		// all filters enabled
		{"", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},
		{"notfound", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},
		{"article", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},
		{"First commit", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},
		{"first commit", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},
		{"testuser", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},

		// no title
		{"title", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "non existent", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},

		// no content
		{"content", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "non existent", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},

		// no tags
		{"article", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "non existent", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},

		// no authors
		{"testuser", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{user.ID + 1}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},

		// no commits
		{"commit", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "non existent", time.Time{}, time.Time{}, "", "", ""}, 0},

		// date filtered by created start
		{"article", &model.SearchArticleFilter{model.BaseSearch{CreatedStart: time.Now().AddDate(0, 0, 1)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},
		{"article", &model.SearchArticleFilter{model.BaseSearch{CreatedStart: time.Now().AddDate(0, 0, -1)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},

		// date filtered by created end
		{"article", &model.SearchArticleFilter{model.BaseSearch{CreatedEnd: time.Now().AddDate(0, 0, 1)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},
		{"article", &model.SearchArticleFilter{model.BaseSearch{CreatedEnd: time.Now().AddDate(0, 0, -2)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},

		// date filtered by updated start
		{"article", &model.SearchArticleFilter{model.BaseSearch{UpdatedStart: time.Now().AddDate(0, 0, 1)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},
		{"article", &model.SearchArticleFilter{model.BaseSearch{UpdatedStart: time.Now().AddDate(0, 0, -2)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},

		// date filtered by updated end
		{"article", &model.SearchArticleFilter{model.BaseSearch{UpdatedEnd: time.Now().AddDate(0, 0, 1)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 1},
		{"article", &model.SearchArticleFilter{model.BaseSearch{UpdatedEnd: time.Now().AddDate(0, 0, -2)}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}, 0},

		// date filtered by published start
		{"article", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Now().AddDate(0, 0, 1), time.Time{}, "", "", ""}, 0},
		{"article", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Now().AddDate(0, 0, -1), time.Time{}, "", "", ""}, 1},

		// date filtered by published end
		{"article", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Now().AddDate(0, 0, 1), "", "", ""}, 1},
		{"article", &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Now().AddDate(0, 0, -1), "", "", ""}, 0},
	}

	for i, io := range inout {
		if articles, _ := SearchArticle(ctx, io.Query, io.Filter); len(articles) != io.N {
			t.Fatalf("Expected %v articles to be found, but was: %v (entry %v, query %v, filter %v)", io.N, len(articles), i+1, io.Query, io.Filter)
		}
	}
}

func TestSearchArticleNoAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	userNoAccess := testutil.CreateUser(t, orga, 321, "noaccess@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, false, false)
	ctx := context.NewEmviUserContext(orga, userNoAccess.ID)
	filter := &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}

	if articles, _ := SearchArticle(ctx, "article", filter); len(articles) != 0 {
		t.Fatal("Article must not be found")
	}
}

func TestSearchArticleUserAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	userWithAccess := testutil.CreateUser(t, orga, 321, "access@user.com")
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	testutil.CreateArticleAccess(t, article, userWithAccess, nil, false)
	ctx := context.NewEmviUserContext(orga, userWithAccess.ID)
	filter := &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}

	if articles, _ := SearchArticle(ctx, "article", filter); len(articles) != 1 {
		t.Fatal("Article must be found")
	}
}

func TestSearchArticleGroupAccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	group := testutil.CreateUserGroup(t, orga, "group")
	userWithGroupAccess := testutil.CreateUser(t, orga, 321, "groupaccess@user.com")
	testutil.CreateUserGroupMember(t, group, userWithGroupAccess, false)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, false, false)
	testutil.CreateArticleAccess(t, article, nil, group, false)
	ctx := context.NewEmviUserContext(orga, userWithGroupAccess.ID)
	filter := &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}
	articles, _ := SearchArticle(ctx, "article", filter)

	if len(articles) != 1 {
		t.Fatal("Article must be found")
	}

	if articles[0].LatestArticleContent == nil {
		t.Fatal("Article must have content")
	}
}

func TestSearchArticleOrderByInSelect(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchArticleFilter{model.BaseSearch{}, 0, false, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "asc", ""}

	if articles, count := SearchArticle(ctx, "article", filter); len(articles) != 1 || count != 1 {
		t.Fatal("One article must be found")
	}
}

func TestSearchArticleArchived(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article1 := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	article1.Archived = null.NewString("archived", true)
	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := model.SaveArticle(nil, article1); err != nil {
		t.Fatal(err)
	}

	filter := &model.SearchArticleFilter{model.BaseSearch{}, 0, true, false, false, false, false, false, "", "", "", []hide.ID{}, []hide.ID{}, nil, "", time.Time{}, time.Time{}, "", "", ""}

	if articles, count := SearchArticle(ctx, "article", filter); len(articles) != 1 || count != 1 {
		t.Fatal("One archived article must be found")
	}
}

func TestSearchArticleMultipleTranslations(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleContent(t, user, article, langDe, 0)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchArticleFilter{}

	if articles, count := SearchArticle(ctx, "article", filter); len(articles) != 1 || count != 1 {
		t.Fatalf("One article must have been found, but was: %v %v", len(articles), count)
	}
}

func TestSearchArticleFilterByAuthor(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	langDe := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleContent(t, user, article, langDe, 0)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchArticleFilter{AuthorUserIds: []hide.ID{user.ID}}

	if articles, count := SearchArticle(ctx, "article", filter); len(articles) != 1 || count != 1 {
		t.Fatalf("One article must have been found, but was: %v %v", len(articles), count)
	}
}

func TestSearchArticleContent(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)

	input := []string{
		"",
		"none",
		"content",
		"! content",
		"test | content",
		"test & content",
	}
	expected := []int{
		1,
		0,
		1,
		0,
		1,
		0,
	}

	for i, in := range input {
		filter := &model.SearchArticleFilter{Content: in}

		if articles, count := SearchArticle(ctx, "article", filter); len(articles) != expected[i] || count != expected[i] {
			t.Fatalf("Expected %v articles to be found for '%v', but was: %v %v", expected[i], in, len(articles), count)
		}
	}
}

func TestSearchArticleWIP(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article.WIP = 5
	ctx := context.NewEmviUserContext(orga, user.ID)

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	filter := &model.SearchArticleFilter{WIP: true}

	if articles, count := SearchArticle(ctx, "article", filter); len(articles) != 1 || count != 1 {
		t.Fatalf("WIP article must have been found, but was: %v %v", len(articles), count)
	}
}

func TestSearchArticleClient(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.SetArticleClientAccess(t, article)
	ctx := context.NewEmviContext(orga, 0, []string{client.Scopes["articles"].String(), client.Scopes["search_articles"].String()}, false)
	articles, count := SearchArticle(ctx, "article", nil)

	if len(articles) != 1 || count != 1 {
		t.Fatalf("Article must have been found, but was: %v %v", len(articles), count)
	}
}

func TestSearchArticleTagsById(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	tag1 := testutil.CreateTag(t, orga, "tag1")
	tag2 := testutil.CreateTag(t, orga, "tag2")
	article1 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	article2 := testutil.CreateArticle(t, orga, user, langEn, true, true)
	testutil.CreateArticleTag(t, article1, tag1)
	testutil.CreateArticleTag(t, article1, tag2)
	testutil.CreateArticleTag(t, article2, tag2)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchArticleFilter{TagIds: []hide.ID{tag1.ID}}
	articles, count := SearchArticle(ctx, "article", filter)

	if len(articles) != 1 || count != 1 {
		t.Fatalf("One article must have been found, but was: %v %v", len(articles), count)
	}

	filter.TagIds = []hide.ID{tag1.ID, tag2.ID}
	articles, count = SearchArticle(ctx, "article", filter)

	if len(articles) != 2 || count != 2 {
		t.Fatalf("Two article must have been found, but was: %v %v", len(articles), count)
	}
}

func TestSearchArticlePreview(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchArticleFilter{Preview: true, PreviewImage: true}
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = `{"type":"doc","content":[{"type":"image","attrs":{"src":"http://localhost:4003/api/v1/content/62RTK394jO.jpg"}},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Headline 1"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]},{"type":"headline","attrs":{"level":3},"content":[{"type":"text","text":"Headline 2"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]},{"type":"headline","attrs":{"level":3},"content":[{"type":"text","text":"Headline 3"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]}]}`

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	articles, count := SearchArticle(ctx, "article", filter)

	if len(articles) != 1 || count != 1 {
		t.Fatalf("Expected one article to be found, but was: %v %v", len(articles), count)
	}

	if articles[0].PreviewImage != "http://localhost:4003/api/v1/content/62RTK394jO.jpg" {
		t.Fatalf("Expected preview image to be returned, but was: %v", articles[0].PreviewImage)
	}

	if articles[0].LatestArticleContent.Content != `<img src="http://localhost:4003/api/v1/content/62RTK394jO.jpg" alt="http://localhost:4003/api/v1/content/62RTK394jO.jpg" /><h2>Headline 1</h2><p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.</p><h3>Headline 2</h3><p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.</p><h3>Headline 3</h3><p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.</p>` {
		t.Fatalf("Preview content not as expected: %v", articles[0].LatestArticleContent.Content)
	}
}

func TestSearchArticlePreviewParagraph(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	langEn := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, langEn, true, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	filter := &model.SearchArticleFilter{Preview: true, PreviewParagraph: true, PreviewImage: true}
	content := model.GetArticleContentLatestByArticleIdAndLanguageId(article.ID, langEn.ID, false)
	content.Content = `{"type":"doc","content":[{"type":"image","attrs":{"src":"http://localhost:4003/api/v1/content/62RTK394jO.jpg"}},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Headline 1"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]},{"type":"headline","attrs":{"level":3},"content":[{"type":"text","text":"Headline 2"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]},{"type":"headline","attrs":{"level":3},"content":[{"type":"text","text":"Headline 3"}]},{"type":"paragraph","content":[{"type":"text","text":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."}]}]}`

	if err := model.SaveArticleContent(nil, content); err != nil {
		t.Fatal(err)
	}

	articles, count := SearchArticle(ctx, "article", filter)

	if len(articles) != 1 || count != 1 {
		t.Fatalf("Expected one article to be found, but was: %v %v", len(articles), count)
	}

	if articles[0].PreviewImage != "http://localhost:4003/api/v1/content/62RTK394jO.jpg" {
		t.Fatalf("Expected preview image to be returned, but was: %v", articles[0].PreviewImage)
	}

	if articles[0].LatestArticleContent.Content != `<p>Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.</p>` {
		t.Fatalf("Preview content not as expected: %v", articles[0].LatestArticleContent.Content)
	}
}

func TestSearchArticleTranslationOrder(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang1 := testutil.CreateLang(t, orga, "en", "English", true)
	lang2 := testutil.CreateLang(t, orga, "de", "Deutsch", false)

	// set user to prefer non standard language
	user.Language.SetValid("de")

	if err := model.SaveUser(nil, user, false); err != nil {
		t.Fatal(err)
	}

	article1 := createTestArticle(t, orga, user, lang1, lang2, "Aaa", "Bbb", "", "")
	article2 := createTestArticle(t, orga, user, lang1, lang2, "Bbb", "Aaa", "", "")
	ctx := context.NewEmviUserContext(orga, user.ID)
	articles, count := SearchArticle(ctx, "", &model.SearchArticleFilter{SortTitle: "asc"})

	if count != 2 {
		t.Fatalf("Two results must be returned, but was: %v", count)
	}

	if articles[0].ID != article2.ID || articles[1].ID != article1.ID {
		t.Fatalf("Articles must be sorted by translated title, but was: %v %v", articles[0].LatestArticleContent.Title, articles[1].LatestArticleContent.Title)
	}
}

func TestSearchArticleLanguage(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang1 := testutil.CreateLang(t, orga, "en", "English", true)
	lang2 := testutil.CreateLang(t, orga, "de", "Deutsch", false)
	createTestArticle(t, orga, user, lang1, lang2, "English 1", "Deutsch 1", "", "")
	createTestArticle(t, orga, user, lang1, lang2, "English 2", "Deutsch 2", "", "")
	ctx := context.NewEmviUserContext(orga, user.ID)
	articles, count := SearchArticle(ctx, "", &model.SearchArticleFilter{LanguageId: lang1.ID})

	if len(articles) != 2 ||
		count != 2 ||
		articles[0].LatestArticleContent.LanguageId != lang1.ID ||
		articles[1].LatestArticleContent.LanguageId != lang1.ID {
		t.Fatalf("Expected articles to be returned in english, but was: %v %v", len(articles), count)
	}

	articles, count = SearchArticle(ctx, "", &model.SearchArticleFilter{LanguageId: lang2.ID})

	if len(articles) != 2 ||
		count != 2 ||
		articles[0].LatestArticleContent.LanguageId != lang2.ID ||
		articles[1].LatestArticleContent.LanguageId != lang2.ID {
		t.Fatalf("Expected articles to be returned in german, but was: %v %v", len(articles), count)
	}
}

func TestSearchArticleRelevance(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article1 := createTestArticle(t, orga, user, lang, lang, "This has relevance", "This has relevance", "Probably less relevant than the second article", "Probably less relevant than the second article")
	article2 := createTestArticle(t, orga, user, lang, lang, "Still relevant", "Still relevant", "Relevance is relevant", "Relevance is relevant")
	ctx := context.NewEmviUserContext(orga, user.ID)
	articles, count := SearchArticle(ctx, "relevance", nil)

	if len(articles) != 2 || count != 2 {
		t.Fatalf("Two results must be found, but was: %v %v", len(articles), count)
	}

	if articles[0].ID != article2.ID || articles[1].ID != article1.ID {
		t.Fatal("Article 2 must be more relevant than article 1")
	}
}

func TestSearchArticleUserGroup(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article := testutil.CreateArticle(t, orga, user, lang, true, true)
	testutil.CreateArticle(t, orga, user, lang, true, true)
	group := testutil.CreateUserGroup(t, orga, "group")
	testutil.CreateArticleAccess(t, article, nil, group, true)
	ctx := context.NewEmviUserContext(orga, user.ID)
	articles, count := SearchArticle(ctx, "", &model.SearchArticleFilter{UserGroupIds: []hide.ID{group.ID}})

	if len(articles) != 1 || count != 1 {
		t.Fatalf("One article must have been found, but was: %v %v", len(articles), count)
	}

	if articles[0].ID != article.ID {
		t.Fatal("Article with group access must have been returned")
	}
}

func TestSearchArticleSortByPublished(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	lang := testutil.CreateLang(t, orga, "en", "English", true)
	article1 := testutil.CreateArticle(t, orga, user, lang, true, true)
	article2 := testutil.CreateArticle(t, orga, user, lang, true, true)
	article3 := testutil.CreateArticle(t, orga, user, lang, true, true)
	article1.Published = null.NewTime(time.Date(2020, 7, 8, 0, 0, 0, 0, time.UTC), true)
	article2.Published = null.NewTime(time.Date(2020, 6, 8, 0, 0, 0, 0, time.UTC), true)
	article3.Published = null.NewTime(time.Date(2020, 7, 5, 0, 0, 0, 0, time.UTC), true)

	if err := model.SaveArticle(nil, article1); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveArticle(nil, article2); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveArticle(nil, article3); err != nil {
		t.Fatal(err)
	}

	tag := testutil.CreateTag(t, orga, "blog")
	testutil.CreateArticleTag(t, article1, tag)
	testutil.CreateArticleTag(t, article2, tag)
	testutil.CreateArticleTag(t, article3, tag)
	ctx := context.NewEmviUserContext(orga, user.ID)
	articles, count := SearchArticle(ctx, "", &model.SearchArticleFilter{Tags: "blog", SortPublished: "asc"})

	if len(articles) != 3 || count != 3 {
		t.Fatalf("Three articles must have been found, but was: %v %v", len(articles), count)
	}

	if articles[0].ID != article2.ID ||
		articles[1].ID != article3.ID ||
		articles[2].ID != article1.ID {
		t.Fatal("Articles in wrong order")
	}

	articles, count = SearchArticle(ctx, "", &model.SearchArticleFilter{Tags: "blog", SortPublished: "desc"})

	if len(articles) != 3 || count != 3 {
		t.Fatalf("Three articles must have been found, but was: %v %v", len(articles), count)
	}

	if articles[0].ID != article1.ID ||
		articles[1].ID != article3.ID ||
		articles[2].ID != article2.ID {
		t.Fatal("Articles in wrong order")
	}
}

func createTestArticle(t *testing.T, orga *model.Organization, user *model.User, lang1, lang2 *model.Language, title1, title2, contentText1, contentText2 string) *model.Article {
	if contentText1 == "" {
		contentText1 = "content 1"
	}

	if contentText2 == "" {
		contentText2 = "content 2"
	}

	article := &model.Article{OrganizationId: orga.ID,
		Views:         54,
		ReadEveryone:  true,
		WriteEveryone: true,
		WIP:           -1,
		Published:     null.NewTime(time.Now(), true)}

	if err := model.SaveArticle(nil, article); err != nil {
		t.Fatal(err)
	}

	content1 := &model.ArticleContent{Title: title1,
		Content:         contentText1,
		TitleTsvector:   title1,
		ContentTsvector: contentText1,
		Version:         0,
		Commit:          null.NewString("Lang 1", true),
		ArticleId:       article.ID,
		LanguageId:      lang1.ID,
		UserId:          user.ID,
		User:            user}

	if err := model.SaveArticleContent(nil, content1); err != nil {
		t.Fatal(err)
	}

	author1 := &model.ArticleContentAuthor{ArticleContentId: content1.ID,
		UserId: user.ID}

	if err := model.SaveArticleContentAuthor(nil, author1); err != nil {
		t.Fatal(err)
	}

	content2 := &model.ArticleContent{Title: title2,
		Content:         contentText2,
		TitleTsvector:   title2,
		ContentTsvector: contentText2,
		Version:         0,
		Commit:          null.NewString("Lang 2", true),
		ArticleId:       article.ID,
		LanguageId:      lang2.ID,
		UserId:          user.ID,
		User:            user}

	if err := model.SaveArticleContent(nil, content2); err != nil {
		t.Fatal(err)
	}

	author2 := &model.ArticleContentAuthor{ArticleContentId: content2.ID,
		UserId: user.ID}

	if err := model.SaveArticleContentAuthor(nil, author2); err != nil {
		t.Fatal(err)
	}

	access := &model.ArticleAccess{Write: true,
		UserId:    user.ID,
		ArticleId: article.ID}

	if err := model.SaveArticleAccess(nil, access); err != nil {
		t.Fatal(err)
	}

	return article
}
