package pages

import (
	"emviwiki/dashboard/auth"
	"emviwiki/dashboard/model"
	"errors"
	"github.com/emvi/hide"
	"github.com/gorilla/mux"
	"net/http"
)

func NewsletterPageHandler(claims *auth.UserTokenClaims, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		if err := deleteNewsletter(r); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	data := struct {
		Newsletter []model.Newsletter
	}{
		model.FindNewsletter(),
	}

	RenderPage(w, newsletterPageTemplate, claims, &data)
}

func deleteNewsletter(r *http.Request) error {
	params := mux.Vars(r)
	id, err := hide.FromString(params["id"])

	if err != nil {
		return err
	}

	newsletter := model.GetNewsletterById(id)

	if newsletter == nil {
		return errors.New("newsletter not found")
	}

	if err := model.DeleteNewsletterById(nil, id); err != nil {
		return err
	}

	return nil
}
