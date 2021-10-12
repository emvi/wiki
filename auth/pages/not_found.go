package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var notFoundPageI18n = i18n.Translation{
	"en": {
		"headline": "Page not found",
		"text":     "The requested page does not exist.",
	},
	"de": {
		"headline": "Seite nicht gefunden",
		"text":     "Die angeforderte Seite existiert nicht.",
	},
}

func NotFoundPageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	renderNotFoundPage(w, r)
}

func renderNotFoundPage(w http.ResponseWriter, r *http.Request) {
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
		notFoundPageI18n[langCode],
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, notFoundPageTemplate, data); err != nil {
		logbuch.Error("Error executing not found template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
