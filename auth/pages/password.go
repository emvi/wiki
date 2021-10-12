package pages

import (
	"bytes"
	"emviwiki/auth/model"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/rest"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/lib/pq"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var passwordPageI18n = i18n.Translation{
	"en": {
		"headline":      "Reset password",
		"email_label":   "Email",
		"submit_button": "Send",
		"back_to_login": "Back",
		"input_err":     "Please enter your email address.",
		"send_mail_err": "Error sending mail.",
	},
	"de": {
		"headline":      "Passwort zurücksetzen",
		"email_label":   "E-Mail-Adresse",
		"submit_button": "Senden",
		"back_to_login": "Zurück",
		"input_err":     "Bitte gib deine E-Mail-Adresse ein.",
		"send_mail_err": "Fehler beim versenden der Mail.",
	},
}

var passwordSuccessPageI18n = i18n.Translation{
	"en": {
		"headline":      "Reset password",
		"text_start":    "Your password has been reset successfully and a new password was sent to",
		"text_end":      ".",
		"back_to_login": "Back to Login",
	},
	"de": {
		"headline":      "Passwort zurücksetzen",
		"text_start":    "Dein Passwort wurde erfolgreich zurückgesetzt und ein neues Passwort wurde an",
		"text_end":      "gesendet.",
		"back_to_login": "Zurück zur Anmeldung",
	},
}

var passwordResetMailI18n = i18n.Translation{
	"en": {
		"title":    "Your password at Emvi has been reset",
		"text":     "Your password has been reset. You'll be asked to change your password on your next login. Your temporary password is:",
		"greeting": "Your password has been reset.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Passwort bei Emvi wurde zurückgesetzt",
		"text":     "Dein Passwort wurde zurückgesetzt. Du wirst bei deinem nächsten Login aufgefordert ein neues Passwort zu vergeben. Dein temporäres Passwort lautet:",
		"greeting": "Dein Passwort wurde zurückgesetzt.",
		"goodbye":  "Dein Emvi Team",
	},
}

const (
	passwordMailCfgName = "password_mail"
	passwordLength      = 12
	saltLength          = 20
)

func PasswordPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderPasswordPage(w, r, "", "")
	} else if r.Method == http.MethodPost {
		handlePasswordReset(w, r, mailProvider)
	}
}

func handlePasswordReset(w http.ResponseWriter, r *http.Request, mail mail.Sender) {
	if err := r.ParseForm(); err != nil {
		logbuch.Warn("Error parsing reset password email form", logbuch.Fields{"err": err})
	}

	email := strings.TrimSpace(r.Form.Get("email"))

	if email == "" {
		renderPasswordPage(w, r, "input_err", email)
		return
	}

	user := model.GetUserByEmail(email)

	if user == nil {
		// ignore errors to prevent giving attackers information about the existance of a mail address
		renderPasswordSuccessPage(w, r, email)
		return
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction on password reset", logbuch.Fields{"err": err, "user_id": user.ID})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	password, hash, salt := newPassword()
	user.Password = null.NewString(hash, true)
	user.PasswordSalt = null.NewString(salt, true)
	user.ResetPassword = true
	user.LastPasswordReset = pq.NullTime{Time: time.Now(), Valid: true}

	if err := model.SaveUser(tx, user); err != nil {
		logbuch.Error("Error saving new user password and salt", logbuch.Fields{"err": err, "user_id": user.ID})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !sendPasswordResetMail(r, user, password, mail) {
		db.Rollback(tx)
		renderPasswordPage(w, r, "send_mail_err", email)
		return
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction on password reset", logbuch.Fields{"err": err, "user_id": user.ID})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	renderPasswordSuccessPage(w, r, email)
}

func newPassword() (string, string, string) {
	password := util.GenRandomString(passwordLength)
	salt := util.GenRandomString(saltLength)
	return password, util.Sha256Base64(password + salt), salt
}

func sendPasswordResetMail(r *http.Request, user *model.User, password string, mail mail.Sender) bool {
	lang := rest.GetSupportedLangCode(r)
	t := mailTplCache.Get()
	data := struct {
		EndVars  map[string]template.HTML
		Vars     map[string]template.HTML
		User     *model.User
		Password string
	}{
		i18n.GetMailEndI18n(lang),
		i18n.GetVars(lang, passwordResetMailI18n),
		user,
		password,
	}
	var buffer bytes.Buffer

	if err := t.ExecuteTemplate(&buffer, passwordMailTemplate, data); err != nil {
		logbuch.Error("Error executing password reset mail", logbuch.Fields{"err": err})
		return false
	}

	subject := i18n.GetMailTitle(lang)[passwordMailCfgName]

	if err := mail(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending password reset mail", logbuch.Fields{"err": err, "email": user.Email})
		return false
	}

	return true
}

func renderPasswordPage(w http.ResponseWriter, r *http.Request, resetErr string, email string) {
	t := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars    map[string]template.HTML
		EndVars     map[string]template.HTML
		Vars        map[string]template.HTML
		Error       string
		Email       string
		Redirect    string
		WebsiteHost string
	}{
		headI18n[langCode],
		endI18n[langCode],
		passwordPageI18n[langCode],
		resetErr,
		email,
		getRedirect(r),
		websiteHost,
	}

	if err := t.ExecuteTemplate(w, passwordPageTemplate, data); err != nil {
		logbuch.Error("Error executing password template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func renderPasswordSuccessPage(w http.ResponseWriter, r *http.Request, email string) {
	t := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars    map[string]template.HTML
		EndVars     map[string]template.HTML
		Vars        map[string]template.HTML
		Email       string
		Redirect    string
		WebsiteHost string
	}{
		headI18n[langCode],
		endI18n[langCode],
		passwordSuccessPageI18n[langCode],
		email,
		getRedirect(r),
		websiteHost,
	}

	if err := t.ExecuteTemplate(w, passwordSuccessPageTemplate, data); err != nil {
		logbuch.Error("Error executing password success template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
