package api

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/rest"
	"net/http"
)

func ReadClientHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if id != 0 {
		c, err := client.ReadClient(ctx.Organization, ctx.UserId, id)

		if err != nil {
			return []error{err}
		}

		rest.WriteResponse(w, c)
	} else {
		clients, err := client.ReadClients(ctx.Organization, ctx.UserId)

		if err != nil {
			return []error{err}
		}

		rest.WriteResponse(w, clients)
	}

	return nil
}

func SaveClientHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := new(client.SaveClientData)

	if err := rest.DecodeJSON(r, req); err != nil {
		return []error{err}
	}

	if err := client.SaveClient(ctx.Organization, ctx.UserId, req, authProvider); err != nil {
		return err
	}

	return nil
}

func DeleteClientHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := client.DeleteClient(ctx.Organization, ctx.UserId, id, authProvider); err != nil {
		return []error{err}
	}

	return nil
}
