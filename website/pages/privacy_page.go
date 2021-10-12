package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var privacyPageI18n = i18n.Translation{
	"en": {
		"page_title": "Emvi — Privacy Policy",
	},
	"de": {
		"page_title": "Emvi — Datenschutzerklärung",
	},
}

func PrivacyPageHandler(w http.ResponseWriter, r *http.Request) {
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
		i18n.GetVars(langCode, privacyPageI18n),
		i18n.GetVars(langCode, navbarComponentI18n),
		i18n.GetVars(langCode, footerComponentI18n),
		backendHost,
		authHost,
		websiteHost,
		clientId,
		version,
		template.HTML(legal.GetCookieNote(langCode)),
		isIntegration,
		template.HTML(legal.GetPrivacyPolicy(langCode)),
	}

	if err := tpl.ExecuteTemplate(w, privacyPageTemplate, &data); err != nil {
		logbuch.Error("Error rendering privacy page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
