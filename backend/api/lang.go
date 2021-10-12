package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/lang"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
)

func GetLangsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	resp := lang.GetLangs(ctx.Organization)
	rest.WriteResponse(w, resp)
	return nil
}

func GetLangHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	language, err := lang.GetLang(ctx.Organization, id)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, language)
	return nil
}

func AddLangHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Code string `json:"code"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := lang.AddLang(ctx.Organization, ctx.UserId, req.Code); err != nil {
		return []error{err}
	}

	return nil
}

func SwitchDefaultLangHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		LanguageId hide.ID `json:"language_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := lang.SwitchDefaultLanguage(ctx.Organization, ctx.UserId, req.LanguageId); err != nil {
		return []error{err}
	}

	return nil
}
