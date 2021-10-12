package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/support"
	"emviwiki/shared/rest"
	"net/http"
)

func ContactSupportHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Type    string `json:"type"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := support.ContactSupport(ctx.Organization, ctx.UserId, req.Type, req.Subject, req.Message, mailProvider); err != nil {
		return err
	}

	return nil
}
