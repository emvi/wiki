package pages

import (
	"emviwiki/auth/user"
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
	"time"
)

var logoutPageI18n = i18n.Translation{
	"en": {
		"headline": "Logout",
		"text":     "You have successfully logged out.",
	},
	"de": {
		"headline": "Abmelden",
		"text":     "Du wurdest erfolgreich abgemeldet.",
	},
}

func LogoutPageHandler(w http.ResponseWriter, r *http.Request) {
	token := rest.GetParam(r, "token")
	redirect := rest.GetParam(r, "redirect_uri")
	handleLogout(w, token)

	if redirect == "" {
		renderLogoutPage(w, r)
	} else {
		http.Redirect(w, r, redirect, http.StatusFound)
	}
}

func handleLogout(w http.ResponseWriter, token string) {
	// delete cookie by setting a negative time to live
	user.SetSessionCookie(w, "", time.Now().Add(-time.Second))
}

func renderLogoutPage(w http.ResponseWriter, r *http.Request) {
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
		logoutPageI18n[langCode],
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, logoutPageTemplate, data); err != nil {
		logbuch.Error("Error executing logout template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
