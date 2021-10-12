package api

import (
	"emviwiki/auth/user"
	"emviwiki/shared/config"
	"emviwiki/shared/rest"
	"net/http"
	"time"
)

func RegistrationHandler(w http.ResponseWriter, r *http.Request) []error {
	req := user.RegistrationData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.Registration(req, rest.GetLangCode(r), mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func RegistrationPasswordHandler(w http.ResponseWriter, r *http.Request) []error {
	req := user.RegistrationPasswordData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.RegistrationPassword(req); err != nil {
		return []error{err}
	}

	return nil
}

func RegistrationPersonalHandler(w http.ResponseWriter, r *http.Request) []error {
	req := user.RegistrationPersonalData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.RegistrationPersonal(req); err != nil {
		return err
	}

	return nil
}

func RegistrationCompletionHandler(w http.ResponseWriter, r *http.Request) []error {
	req := user.RegistrationCompletionData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	lang := rest.GetLangCode(r)
	token, expires, err := user.CompleteRegistration(req, lang, mailProvider)

	if err != nil {
		return err
	}

	rest.WriteResponse(w, struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Domain      string `json:"domain"`
		Secure      bool   `json:"secure"`
	}{
		token,
		int64(expires.Sub(time.Now()).Seconds()),
		config.Get().JWT.CookieDomainName,
		config.Get().Server.HTTP.SecureCookies,
	})
	return nil
}

func CancelRegistrationHandler(w http.ResponseWriter, r *http.Request) []error {
	code := rest.GetParam(r, "code")

	if err := user.CancelRegistration(code); err != nil {
		return []error{err}
	}

	return nil
}

func ConfirmRegistrationHandler(w http.ResponseWriter, r *http.Request) []error {
	code := rest.GetParam(r, "code")
	step, err := user.ConfirmRegistration(code)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Step int `json:"step"`
	}{step})
	return nil
}
