package api

import (
	"emviwiki/backend/newsletter"
	"emviwiki/shared/rest"
	"net/http"
)

func SubscribeNewsletterHandler(w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Email string `json:"email"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := newsletter.Subscribe(req.Email, newsletter.NewsletterList, rest.GetSupportedLangCode(r), mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func SubscribeOnPremiseNewsletterHandler(w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Email string `json:"email"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := newsletter.Subscribe(req.Email, newsletter.NewsletterOnPremiseList, rest.GetLangCode(r), mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func ConfirmNewsletterHandler(w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Code string `json:"code"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := newsletter.Confirm(req.Code); err != nil {
		return []error{err}
	}

	return nil
}

func UnsubscribeNewsletterHandler(w http.ResponseWriter, r *http.Request) []error {
	code := rest.GetParam(r, "code")

	if err := newsletter.Unsubscribe(code); err != nil {
		return []error{err}
	}

	return nil
}
