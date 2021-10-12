package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var clientUnknownPageI18n = i18n.Translation{
	"en": {
		"headline": "Client unknown",
		"text":     "The client you're trying to authorize was not found in our system.",
	},
	"de": {
		"headline": "Anwendung unbekannt",
		"text":     "Die Anwendung die du versuchst zu autorisieren ist unbekannt.",
	},
}

func renderClientUnknownPage(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars    map[string]template.HTML
		EndVars     map[string]template.HTML
		Vars        map[string]template.HTML
		WebsiteHost string
	}{
		headI18n[langCode],
		endI18n[langCode],
		clientUnknownPageI18n[langCode],
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, clientUnknownPageTemplate, data); err != nil {
		logbuch.Error("Error executing client unknown template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
