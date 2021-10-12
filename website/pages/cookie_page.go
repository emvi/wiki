package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var cookiePageI18n = i18n.Translation{
	"en": {
		"page_title": "Emvi — Cookie Policy",
	},
	"de": {
		"page_title": "Emvi — Cookie-Richtlinien",
	},
}

func CookiePageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		LangCode      string
		IsBlog        bool
		Vars          map[string]template.HTML
		NavbarVars    map[string]template.HTML
		FooterVars    map[string]template.HTML
		BackendHost   string
		AuthHost      string
		WebsiteHost   string
		AuthClientID  string
		Version       string
		CookiesNote   template.HTML
		IsIntegration bool
		Content       template.HTML
	}{
		langCode,
		false,
		i18n.GetVars(langCode, cookiePageI18n),
		i18n.GetVars(langCode, navbarComponentI18n),
		i18n.GetVars(langCode, footerComponentI18n),
		backendHost,
		authHost,
		websiteHost,
		clientId,
		version,
		template.HTML(legal.GetCookieNote(langCode)),
		isIntegration,
		template.HTML(legal.GetCookiePolicy(langCode)),
	}

	if err := tpl.ExecuteTemplate(w, cookiePageTemplate, &data); err != nil {
		logbuch.Error("Error rendering cookie page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
