package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"fmt"
	emvi "github.com/emvi/api-go"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
)

const (
	publishedFormat = "2006-01-02"
)

var articlePageI18n = i18n.Translation{
	"en": {
		"page_title": "Emvi — Blog — %s",
		"title":      "Blog",
		"back":       "Go to all Articles",
	},
	"de": {
		"page_title": "Emvi — Blog — %s",
		"title":      "Blog",
		"back":       "Zu allen Artikeln",
	},
}

func ArticlePageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	params := mux.Vars(r)
	langCode := rest.GetSupportedLangCode(r)
	articleId := getIdFromSlug(params["slug"])
	langId := GetBlogLangId(langCode)
	article, content, _, err := blogClient.GetArticle(articleId, langId, 0)

	if err != nil {
		logbuch.Debug("Error finding blog article", logbuch.Fields{"err": err, "id": articleId, "lang": langId})
		http.Redirect(w, r, "/notfound", http.StatusFound)
		return
	}

	vars := copyMap(i18n.GetVars(langCode, articlePageI18n))
	vars["page_title"] = template.HTML(fmt.Sprintf(string(vars["page_title"]), content.Title))
	data := struct {
		LangCode       string
		IsBlog         bool
		Vars           map[string]template.HTML
		NavbarVars     map[string]template.HTML
		FooterVars     map[string]template.HTML
		NewsletterVars map[string]template.HTML
		BackendHost    string
		AuthHost       string
		WebsiteHost    string
		AuthClientID   string
		Version        string
		CookiesNote    template.HTML
		IsIntegration  bool
		Article        *emvi.Article
		Title          string
		Content        template.HTML
		Published      string
	}{
		langCode,
		true,
		vars,
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
		article,
		content.Title,
		template.HTML(content.Content),
		article.Published.Format(publishedFormat),
	}

	if err := tpl.ExecuteTemplate(w, articlePageTemplate, &data); err != nil {
		logbuch.Error("Error rendering article page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getIdFromSlug(slug string) string {
	parts := strings.Split(slug, "-")

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return ""
}

func copyMap(m map[string]template.HTML) map[string]template.HTML {
	n := make(map[string]template.HTML, len(m))

	for k, v := range m {
		n[k] = v
	}

	return n
}
