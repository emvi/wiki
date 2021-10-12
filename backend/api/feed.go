package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/feed"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
)

func GetFilteredFeedHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	userIds, err := rest.GetIdParams(r, "user")

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchFeedFilter{
		baseFilter,
		rest.GetBoolParam(r, "notifications"),
		rest.GetBoolParam(r, "unread"),
		rest.GetParams(r, "reason"),
		userIds,
	}

	results, count := feed.GetFilteredFeed(ctx.Organization, ctx.UserId, filter)

	for i := range results {
		if results[i].TriggeredByUser.Picture.Valid {
			results[i].TriggeredByUser.Picture.SetValid(getResourceURL(results[i].TriggeredByUser.Picture.String))
		}
	}

	rest.WriteResponse(w, struct {
		Feed  []model.Feed `json:"feed"`
		Count int          `json:"count"`
	}{results, count})
	return nil
}

func ToggleNotificationReadHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Id hide.ID `json:"id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := feed.ToggleNotificationRead(nil, ctx.Organization, ctx.UserId, req.Id); err != nil {
		return []error{err}
	}

	return nil
}
