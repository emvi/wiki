package api

import (
	"emviwiki/auth/client"
	"emviwiki/shared/rest"
	"net/http"
)

func RegisterClientHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Name   string            `json:"name"`
		Scopes map[string]string `json:"scopes"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	c, err := client.NewClient(req.Name, req.Scopes)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, c)
	return nil
}

func DeleteClientHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	clientId := rest.GetParam(r, "client_id")
	clientSecret := rest.GetParam(r, "client_secret")

	if err := client.DeleteClient(clientId, clientSecret); err != nil {
		return []error{err}
	}

	return nil
}
