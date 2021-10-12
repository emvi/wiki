package pages

import (
	"emviwiki/auth/model"
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"html/template"
	"net/http"
)

var passwordResetPageI18n = i18n.Translation{
	"en": {
		"headline":           "Reset password",
		"oldpwd_label":       "Temporary password",
		"newpwd1_label":      "New password",
		"newpwd2_label":      "Repeat new password",
		"submit_button":      "Confirm password",
		"forgot_password":    "Forgot password?",
		"input_err":          "Please fill out the form.",
		"match_err":          "The new password does not match.",
		"password_rules_err": "The password does not match the password rules.",
		"password_err":       "The temporary password is incorrect.",
	},
	"de": {
		"headline":           "Passwort zurücksetzen",
		"oldpwd_label":       "Temporäres Passwort",
		"newpwd1_label":      "Neues Passwort",
		"newpwd2_label":      "Neues Passwort wiederholen",
		"submit_button":      "Passwort bestätigen",
		"forgot_password":    "Passwort vergessen?",
		"input_err":          "Bitte fülle das Formular aus.",
		"match_err":          "Das neue Passwort stimmt nicht überein.",
		"password_rules_err": "Das neue Passwort entspricht nicht den Passwort Regeln.",
		"password_err":       "Das temporäre Passwort ist falsch.",
	},
}

func PasswordResetPageHandler(w http.ResponseWriter, r *http.Request) {
	user := loggedIn(r)

	if user == nil {
		redirectToLogin(w, r)
		return
	}

	if r.Method == http.MethodGet {
		renderPasswordResetPage(w, r, "")
	} else if r.Method == http.MethodPost {
		handleResetPassword(w, r, user)
	}
}

func handleResetPassword(w http.ResponseWriter, r *http.Request, user *model.User) {
	if err := r.ParseForm(); err != nil {
		logbuch.Warn("Error parsing reset password form", logbuch.Fields{"err": err})
	}

	oldpwd := r.Form.Get("oldpwd")
	newpwd1 := r.Form.Get("newpwd1")
	newpwd2 := r.Form.Get("newpwd2")

	if oldpwd == "" || newpwd1 == "" || newpwd2 == "" {
		renderPasswordResetPage(w, r, "input_err")
		return
	}

	if newpwd1 != newpwd2 {
		renderPasswordResetPage(w, r, "match_err")
		return
	}

	pwd := util.Sha256Base64(oldpwd + user.PasswordSalt.String)

	if pwd != user.Password.String {
		renderPasswordResetPage(w, r, "password_err")
		return
	}

	user.Password = null.NewString(util.Sha256Base64(newpwd1+user.PasswordSalt.String), true)
	user.ResetPassword = false

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving new user password and salt on password reset", logbuch.Fields{"err": err, "user_id": user.ID})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// this should redirect to authorization page
	loginRedirect(w, r, user)
}

func renderPasswordResetPage(w http.ResponseWriter, r *http.Request, resetErr string) {
	tpl := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars    map[string]template.HTML
		EndVars     map[string]template.HTML
		Vars        map[string]template.HTML
		Error       string
		Redirect    string
		WebsiteHost string
	}{
		headI18n[langCode],
		endI18n[langCode],
		passwordResetPageI18n[langCode],
		resetErr,
		getRedirect(r),
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, passwordResetPageTemplate, data); err != nil {
		logbuch.Error("Error executing password reset template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
