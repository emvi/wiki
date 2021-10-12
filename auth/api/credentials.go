package api

import (
	"emviwiki/auth/client"
	"emviwiki/shared/rest"
	"net/http"
)

func ClientCredentialsHandler(w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		GrantType    string `json:"grant_type"`
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	resp, err := client.ValidateClientCredentials(req.GrantType, req.ClientId, req.ClientSecret)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, resp)
	return nil
}
