package pages

import (
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	"emviwiki/auth/scope"
	"emviwiki/shared/constants"
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"fmt"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var authorizationPageI18n = i18n.Translation{
	"en": {
		"headline":      "Authorize",
		"text":          "The application will gain the following access rights",
		"noscopes_text": "Authorize application",
		"submit_button": "Authorize",
	},
	"de": {
		"headline":      "Autorisierung",
		"text":          "Der Anwendung werden folgende Zugriffsrechte gewährt",
		"noscopes_text": "Anwendung autorisieren",
		"submit_button": "Autorisieren",
	},
}

type Scope struct {
	Config *scope.ScopeConfig
	Value  string
}

func AuthorizationPageHandler(w http.ResponseWriter, r *http.Request) {
	client, scopes := validateAuthParams(r)

	if client == nil {
		renderClientUnknownPage(w, r)
		return
	}

	user := loggedIn(r)

	if user == nil {
		redirectToLogin(w, r)
		return
	} else if user.ResetPassword {
		loginURL, _ := url.Parse("/auth/passwordreset")
		query := loginURL.Query()
		query.Add("redirect", r.URL.String())
		loginURL.RawQuery = query.Encode()
		http.Redirect(w, r, loginURL.String(), http.StatusFound)
		return
	}

	if client.Trusted || model.GetAccessGrantByUserIdAndClientId(user.ID, client.ID) != nil {
		handleAuth(w, r, client, user, scopes, true)
		return
	}

	if r.Method == http.MethodGet {
		renderAuthPage(w, r, client, scopes)
	} else if r.Method == http.MethodPost {
		handleAuth(w, r, client, user, scopes, false)
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request, client *model.Client, user *model.User, scopes []Scope, skipAccessGrant bool) {
	if !skipAccessGrant && !saveAccessGrant(client, user, scopes) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, ttl, err := jwt.NewUserToken(&jwt.UserTokenClaims{UserId: user.ID,
		Language: user.Language.String,
		Scopes:   scopesToString(scopes)})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redirectURL, err := url.Parse(client.RedirectURI.String)

	if err != nil {
		logbuch.Error("Error parsing client redirect URI", logbuch.Fields{"err": err, "redirect_uri": client.RedirectURI})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, getRedirectURL(r, client, redirectURL, token, ttl), http.StatusFound)
}

func saveAccessGrant(client *model.Client, user *model.User, scopes []Scope) bool {
	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction on authorization", logbuch.Fields{"err": err, "user_id": user.ID, "client_id": client.ID})
		return false
	}

	access := &model.AccessGrant{UserId: user.ID, ClientId: client.ID}

	if err := model.SaveAccessGrant(tx, access); err != nil {
		logbuch.Error("Error saving access grant on authorization", logbuch.Fields{"err": err, "user_id": user.ID, "client_id": client.ID})
		return false
	}

	for _, s := range scopes {
		s := &model.Scope{ClientId: client.ID, Key: s.Config.Name, Value: s.Value}

		if err := model.SaveScope(tx, s); err != nil {
			logbuch.Error("Error saving scope on authorization", logbuch.Fields{"err": err, "user_id": user.ID, "client_id": client.ID})
			return false
		}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction on authorization", logbuch.Fields{"err": err, "user_id": user.ID, "client_id": client.ID})
		return false
	}

	return true
}

func scopesToString(scopes []Scope) []string {
	scopesStr := make([]string, len(scopes))

	for i := range scopes {
		scopesStr[i] = fmt.Sprintf("%s:%s", scopes[i].Config.Name, scopes[i].Value)
	}

	return scopesStr
}

func validateAuthParams(r *http.Request) (*model.Client, []Scope) {
	respType := strings.ToLower(rest.GetParam(r, "response_type"))
	clientID := rest.GetParam(r, "client_id")
	redirectURI := rest.GetParam(r, "redirect_uri")
	s := rest.GetParam(r, "scope")

	if respType != "token" || clientID == "" || redirectURI == "" {
		return nil, nil
	}

	redirectURL, err := url.Parse(redirectURI)

	if err != nil {
		logbuch.Warn("Could not parse redirect URI", logbuch.Fields{"err": err, "url": redirectURI})
		return nil, nil
	}

	client := model.GetClientByClientIdAndRedirectURI(clientID, redirectURL.String())

	if client == nil {
		return nil, nil
	}

	scopes := getScopes(s)

	if scopes == nil {
		return nil, nil
	}

	return client, scopes
}

func getScopes(scopesStr string) []Scope {
	if scopesStr == "" {
		return []Scope{}
	}

	scopes := strings.Split(scopesStr, " ")
	cfgScopes := make([]Scope, len(scopes))

	for i, s := range scopes {
		s := strings.Split(s, ":")

		if len(s) != 1 && len(s) != 2 {
			logbuch.Warn("Scope invalid", logbuch.Fields{"scope": s})
			return nil
		}

		cfg := scope.GetScope(s[0])

		if cfg == nil {
			logbuch.Warn("Scope not found", logbuch.Fields{"scope": s})
			return nil
		}

		if len(s) == 2 && !cfg.AllowValue {
			logbuch.Warn("Scope does not allow value", logbuch.Fields{"scope": s})
			return nil
		}

		value := ""

		if len(s) == 2 {
			value = s[1]
		}

		cfgScopes[i] = Scope{cfg, value}
	}

	return cfgScopes
}

func getRedirectURL(r *http.Request, client *model.Client, redirectURL *url.URL, token string, ttl time.Time) string {
	query := redirectURL.Query()

	// don't need return the token and everything else to trusted clients, as they're ours
	if !client.Trusted {
		query.Add("token_type", constants.AuthTokenType)
		query.Add("expires_in", strconv.Itoa(int(ttl.Sub(time.Now()).Seconds())))
		query.Add("access_token", token)
	}

	state := rest.GetParam(r, "state")

	if state != "" {
		query.Add("state", state)
	}

	redirectURL.RawQuery = query.Encode()
	return redirectURL.String()
}

func renderAuthPage(w http.ResponseWriter, r *http.Request, client *model.Client, scopes []Scope) {
	tpl := tplCache.Get()
	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		HeadVars    map[string]template.HTML
		EndVars     map[string]template.HTML
		Vars        map[string]template.HTML
		Client      *model.Client
		Scopes      []Scope
		Lang        string
		WebsiteHost string
	}{
		headI18n[langCode],
		endI18n[langCode],
		authorizationPageI18n[langCode],
		client,
		scopes,
		langCode,
		websiteHost,
	}

	if err := tpl.ExecuteTemplate(w, authorizationPageTemplate, data); err != nil {
		logbuch.Error("Error executing authorization template", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
