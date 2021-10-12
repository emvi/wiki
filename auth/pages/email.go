package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var updateUserEmailPageI18n = i18n.Translation{
	"en": {
		"headline":        "Email address confirmed!",
		"text":            "Your email address has been updated successfully. In order for the change to take effect, you must log in again.",
		"back_to_website": "Back to Website",
		"User not found":  "The email address you are trying to update could not be found.",
	},
	"de": {
		"headline":        "E-Mail-Adresse bestätigt!",
		"text":            "Deine E-Mail-Adresse wurde erfolgreich geändert. Damit die Änderung wirksam wird, musst du dich neu anmelden.",
		"back_to_website": "Zurück zur Website",
		"User not found":  "Die E-Mail-Adresse, die du zu ändern versuchst, wurde nicht gefunden.",
	},
}

func UpdateUserEmailPageHandler(w http.ResponseWriter, r *http.Request) {
	renderUpdateEmailPage(w, r)
}

func renderUpdateEmailPage(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars    map[string]template.HTML
		EndVars     map[string]template.HTML
		Vars        map[string]template.HTML
		Error       string
		WebsiteHost string
	}{
		headI18n[langCode],
		endI18n[langCode],
		updateUserEmailPageI18n[langCode],
		rest.GetParam(r, "error"),
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, updateEmailPageTemplate, data); err != nil {
		logbuch.Error("Error executing update email template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
