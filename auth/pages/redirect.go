package pages

import (
	"emviwiki/shared/rest"
	"net/http"
	"net/url"
)

func getRedirect(r *http.Request) string {
	redirect := rest.GetParam(r, "redirect")

	if redirect != "" {
		return "?redirect=" + url.QueryEscape(redirect)
	}

	return ""
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	loginURL, _ := url.Parse("/auth/login")
	query := loginURL.Query()
	query.Add("redirect", r.URL.String())
	loginURL.RawQuery = query.Encode()
	http.Redirect(w, r, loginURL.String(), http.StatusFound)
}
