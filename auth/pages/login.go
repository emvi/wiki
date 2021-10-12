package pages

import (
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	authuser "emviwiki/auth/user"
	"emviwiki/shared/constants"
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxLoginAttempts    = 5
	loginBlockedMinutes = 5
)

var loginPageI18n = i18n.Translation{
	"en": {
		"headline":               "Login",
		"email_label":            "Email",
		"password_label":         "Password",
		"forgot_password":        "Forgot password?",
		"submit_button":          "Login",
		"input_err":              "Please enter your email address and password.",
		"login_err":              "Incorrect email address or password.",
		"attempts_err":           "Maximum login attempts reached, please wait 5 minutes and try again.",
		"account":                "No account yet?",
		"signup":                 "Sign up",
		"button_login_google":    "Sign in with Google",
		"button_login_slack":     "Sign in with Slack",
		"button_login_github":    "Sign in with GitHub",
		"button_login_microsoft": "Sign in with Microsoft",
		"or":                     "or",
	},
	"de": {
		"headline":               "Anmelden",
		"email_label":            "E-Mail-Adresse",
		"password_label":         "Passwort",
		"forgot_password":        "Passwort vergessen?",
		"submit_button":          "Anmelden",
		"input_err":              "Bitte gib deine E-Mail-Adresse und Passwort ein.",
		"login_err":              "Ungültige E-Mail-Adresse oder Passwort.",
		"attempts_err":           "Maximum Loginversuchen erreicht, bitte warte 5 Minuten und versuche es erneut.",
		"account":                "Noch kein Konto?",
		"signup":                 "Konto anlegen",
		"button_login_google":    "Mit Google anmelden",
		"button_login_slack":     "Mit Slack anmelden",
		"button_login_github":    "Mit GitHub anmelden",
		"button_login_microsoft": "Mit Microsoft anmelden",
		"or":                     "oder",
	},
}

var loginSuccessPageI18n = i18n.Translation{
	"en": {
		"headline": "Login successful",
		"text":     "Your login was successful.",
	},
	"de": {
		"headline": "Anmeldung erfolgreich",
		"text":     "Deine Anmeldung war erfolgreich.",
	},
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	user := loggedIn(r)

	if user != nil {
		loginRedirect(w, r, user)
		return
	}

	if r.Method == http.MethodGet {
		renderLoginPage(w, r, "", "")
	} else if r.Method == http.MethodPost {
		handleLogin(w, r)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		logbuch.Warn("Error parsing login form", logbuch.Fields{"err": err})
	}

	email := strings.TrimSpace(r.PostForm.Get("email"))
	pwd := r.PostForm.Get("password")

	if email == "" || pwd == "" {
		renderLoginPage(w, r, "input_err", email)
		return
	}

	// check email and password
	user := model.GetUserByEmail(email)

	if user == nil {
		renderLoginPage(w, r, "login_err", email)
		return
	}

	// check login attempts
	lastLoginAttempt := user.LastLoginAttempt
	fiveMinAgo := time.Now().Add(-loginBlockedMinutes * time.Minute)

	if user.LoginAttempts >= maxLoginAttempts && lastLoginAttempt.After(fiveMinAgo) {
		renderLoginPage(w, r, "attempts_err", email)
		return
	}

	userAttempts := user
	pwd = util.Sha256Base64(pwd + user.PasswordSalt.String)
	user = model.GetUserByEmailAndPassword(email, pwd)

	if user == nil {
		updateLoginAttempts(userAttempts)
		renderLoginPage(w, r, "login_err", email)
		return
	}

	// create session
	session, expires, err := jwt.NewUserToken(&jwt.UserTokenClaims{UserId: user.ID, Language: user.Language.String})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authuser.SetSessionCookie(w, session, expires)

	// update last login and login attempts
	user.LastLogin = time.Now()
	user.LoginAttempts = 0
	user.LastLoginAttempt = user.LastLogin

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when updating last login time", logbuch.Fields{"err": err, "id": user.ID})
		// continue login
	}

	go saveUserLogin(user.ID)
	loginRedirect(w, r, user)
}

func updateLoginAttempts(user *model.User) {
	user.LoginAttempts++

	if user.LastLoginAttempt.Before(time.Now().Add(-loginBlockedMinutes * time.Minute)) {
		user.LastLoginAttempt = time.Now()
	}

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when updating login attempts", logbuch.Fields{"err": err, "id": user.ID})
	}
}

func loginRedirect(w http.ResponseWriter, r *http.Request, user *model.User) {
	redirect := rest.GetParam(r, "redirect")

	if user.ResetPassword {
		resetURL, _ := url.Parse("/auth/passwordreset")
		query := resetURL.Query()
		query.Add("redirect", redirect)
		resetURL.RawQuery = query.Encode()
		http.Redirect(w, r, resetURL.String(), http.StatusFound)
	} else if redirect != "" {
		http.Redirect(w, r, redirect, http.StatusFound)
	} else {
		renderLoginSuccessPage(w, r)
	}
}

func saveUserLogin(userId hide.ID) {
	login := &model.Login{UserId: userId}

	if err := model.SaveLogin(nil, login); err != nil {
		logbuch.Error("Error saving login entity on user login", logbuch.Fields{"err": err, "user_id": userId})
	}
}

func renderLoginPage(w http.ResponseWriter, r *http.Request, loginErr string, email string) {
	tpl := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars             map[string]template.HTML
		EndVars              map[string]template.HTML
		Vars                 map[string]template.HTML
		Error                string
		Email                string
		Redirect             string
		GitHubSSOClientId    string
		SlackSSOClientId     string
		GoogleSSOClientId    string
		MicrosoftSSOClientId string
		AuthHost             string
		WebsiteHost          string
	}{
		headI18n[langCode],
		endI18n[langCode],
		loginPageI18n[langCode],
		loginErr,
		email,
		getRedirect(r),
		githubSSOClientId,
		slackSSOClientId,
		googleSSOClientId,
		microsoftSSOClientId,
		authHost,
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, loginPageTemplate, data); err != nil {
		logbuch.Error("Error executing login template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func renderLoginSuccessPage(w http.ResponseWriter, r *http.Request) {
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
		loginSuccessPageI18n[langCode],
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, loginSuccessPageTemplate, data); err != nil {
		logbuch.Error("Error executing login success template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func loggedIn(r *http.Request) *model.User {
	cookie, err := r.Cookie(constants.AuthCookieName)

	if err != nil {
		if err != http.ErrNoCookie {
			logbuch.Warn("Error reading cookie for authentication", logbuch.Fields{"err": err})
		}

		return nil
	}

	if cookie.Value == "" {
		logbuch.Warn("Error reading cookie value for authentication (empty)")
		return nil
	}

	session := jwt.GetUserTokenClaims(cookie.Value)

	if session == nil {
		return nil
	}

	return model.GetUserById(session.UserId)
}
