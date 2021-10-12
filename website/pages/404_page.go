package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var notfoundPageI18n = i18n.Translation{
	"en": {
		"page_title": "Emvi — Page not found",
		"title":      "Not found",
		"text":       "The page does not exist.",
	},
	"de": {
		"page_title": "Emvi — Seite nicht gefunden",
		"title":      "Nicht gefunden",
		"text":       "Die Seite existiert nicht.",
	},
}

func NotFoundPageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	langCode := rest.GetSupportedLangCode(r)
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
	}{
		langCode,
		false,
		i18n.GetVars(langCode, notfoundPageI18n),
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
	}

	if err := tpl.ExecuteTemplate(w, notfoundPageTemplate, &data); err != nil {
		logbuch.Error("Error rendering 404 page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
