package pages

import (
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	"emviwiki/auth/sso"
	user2 "emviwiki/auth/user"
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const (
	githubProviderName    = "github"
	slackProviderName     = "slack"
	googleProviderName    = "google"
	microsoftProviderName = "microsoft"
)

var ssoPageI18n = i18n.Translation{
	"en": {
		"headline": "Error on Authorization",
		"text":     "An error occurred while authorizing Emvi for your account. Please try again.",
	},
	"de": {
		"headline": "Fehler bei der Autorisierung",
		"text":     "Während der Autorisierung deines Account für Emvi ist ein Fehler aufgetreten. Bitte versuche es erneut.",
	},
}

func SSOPageHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	providerName := strings.ToLower(params["provider"])
	code := rest.GetParam(r, "code")
	var provider sso.SSOProvider

	if providerName == githubProviderName {
		provider = githubSSOProvider
	} else if providerName == slackProviderName {
		provider = slackSSOProvider
	} else if providerName == googleProviderName {
		provider = googleSSOProvider
	} else if providerName == microsoftProviderName {
		provider = microsoftSSOProvider
	} else {
		logbuch.Error("Unknown provider", logbuch.Fields{"provider": providerName})
		renderSSOErrorPage(w, r)
		return
	}

	token, err := provider.GetToken(code)

	if err != nil {
		logbuch.Error("Error authenticating SSO user", logbuch.Fields{"err": err})
		renderSSOErrorPage(w, r)
		return
	}

	user, err := provider.GetUser(token.AccessToken)

	if err != nil {
		logbuch.Error("Error obtaining SSO user information", logbuch.Fields{"err": err})
		renderSSOErrorPage(w, r)
		return
	}

	userEntity, err := createSSOUser(providerName, user.Id, user.Email, user.Name, user.Picture)

	if err != nil {
		logbuch.Error("Error creating/updating SSO user", logbuch.Fields{"err": err})
		renderSSOErrorPage(w, r)
		return
	}

	userToken, expires, err := loginSSOUser(userEntity)

	if err != nil {
		logbuch.Error("Error logging in SSO user", logbuch.Fields{"err": err})
		renderSSOErrorPage(w, r)
		return
	}

	logbuch.Debug("SSO user logged in successfully", logbuch.Fields{"id": userEntity.ID, "email": userEntity.Email, "provider": userEntity.AuthProvider})
	user2.SetSessionCookie(w, userToken, expires)
	http.Redirect(w, r, fmt.Sprintf("%s/organizations", websiteHost), http.StatusFound)
}

func createSSOUser(providerName, id, email, name, picture string) (*model.User, error) {
	user := model.GetUserByAuthProviderAndUserId(providerName, id)

	if user == nil {
		logbuch.Debug("Creating new SSO user", logbuch.Fields{"email": email, "name": name})
		user = &model.User{
			AcceptMarketing:    false,
			Active:             true,
			AuthProvider:       providerName,
			AuthProviderUserId: null.NewString(id, true),
		}
	}

	firstname, lastname := splitName(name)
	user.Firstname = null.NewString(firstname, firstname != "")
	user.Lastname = null.NewString(lastname, lastname != "")
	user.Email = email
	user.PictureURL = null.NewString(picture, picture != "")
	user.LastLogin = time.Now()

	if err := model.SaveUser(nil, user); err != nil {
		return nil, err
	}

	logbuch.Debug("SSO user created/updated", logbuch.Fields{"provider": providerName, "id": id, "email": email, "name": name, "picture": picture})
	return user, nil
}

func splitName(name string) (string, string) {
	parts := strings.Split(name, " ")

	if len(parts) <= 1 {
		return "", name
	} else if len(parts) == 2 {
		return parts[0], parts[1]
	}

	return strings.Join(parts[:len(parts)-1], " "), parts[len(parts)-1]
}

func loginSSOUser(user *model.User) (string, time.Time, error) {
	logbuch.Debug("Logging in SSO user", logbuch.Fields{"id": user.ID, "email": user.Email})
	token, expires, err := jwt.NewUserToken(&jwt.UserTokenClaims{UserId: user.ID, Language: user.Language.String, IsSSOUser: true})

	if err != nil {
		return "", expires, err
	}

	go saveUserLogin(user.ID)
	return token, expires, nil
}

func renderSSOErrorPage(w http.ResponseWriter, r *http.Request) {
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
		ssoPageI18n[langCode],
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, ssoErrorPageTemplate, data); err != nil {
		logbuch.Error("Error executing sso error template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
