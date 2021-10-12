package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	emvi "github.com/emvi/api-go"
	"github.com/emvi/logbuch"
	"github.com/gosimple/slug"
	"html/template"
	"net/http"
)

const articlesPerPage = 20

var blogPageI18n = i18n.Translation{
	"en": {
		"page_title":         "Emvi — Blog",
		"title":              "Blog",
		"search_placeholder": "Search...",
		"load_more":          "Load more",
		"no_results":         "No results found.",
	},
	"de": {
		"page_title":         "Emvi — Blog",
		"title":              "Blog",
		"search_placeholder": "Suchen...",
		"load_more":          "Mehr laden",
		"no_results":         "Keine Ergebnisse gefunden.",
	},
}

type Article struct {
	Article emvi.Article
	Content template.HTML
	Slug    string
}

func BlogPageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	langCode := rest.GetSupportedLangCode(r)
	langId := GetBlogLangId(langCode)
	articlesOnly := rest.GetBoolParam(r, "articles_only")
	query := rest.GetParam(r, "query")
	offset, _ := rest.GetIntParam(r, "offset")
	articles, count, err := blogClient.FindArticles(query, &emvi.ArticleFilter{
		BaseSearch: emvi.BaseSearch{
			Offset: offset,
			Limit:  articlesPerPage,
		},
		SortPublished:    emvi.SortDescending,
		PreviewParagraph: true,
		PreviewImage:     true,
		LanguageId:       langId,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var articlesWithContent []Article

	for _, article := range articles {
		articlesWithContent = append(articlesWithContent, Article{
			article,
			template.HTML(article.LatestArticleContent.Content),
			slug.Make(article.LatestArticleContent.Title) + "-" + article.Id,
		})
	}

	data := struct {
		LangCode           string
		IsBlog             bool
		Vars               map[string]template.HTML
		NavbarVars         map[string]template.HTML
		FooterVars         map[string]template.HTML
		NewsletterVars     map[string]template.HTML
		BackendHost        string
		AuthHost           string
		WebsiteHost        string
		AuthClientID       string
		Version            string
		CookiesNote        template.HTML
		IsIntegration      bool
		Articles           []Article
		ArticleCount       int
		ArticlesPerPage    int
		RenderArticlesOnly bool
		Published          string
	}{
		langCode,
		true,
		i18n.GetVars(langCode, blogPageI18n),
		i18n.GetVars(langCode, navbarComponentI18n),
		i18n.GetVars(langCode, footerComponentI18n),
		i18n.GetVars(langCode, newsletterComponentI18n),
		backendHost,
		authHost,
		websiteHost,
		clientId,
		version,
		template.HTML(legal.GetCookieNote(langCode)),
		isIntegration,
		articlesWithContent,
		count,
		articlesPerPage,
		articlesOnly,
		"",
	}

	if err := tpl.ExecuteTemplate(w, blogPageTemplate, &data); err != nil {
		logbuch.Error("Error rendering blog page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
