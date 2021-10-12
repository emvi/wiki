package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/tag"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"github.com/gorilla/mux"
	"net/http"
)

func AddTagHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := tag.AddTagData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := tag.AddTag(ctx.Organization, req); err != nil {
		return []error{err}
	}

	return nil
}

func RenameTagHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Id  hide.ID `json:"id"`
		Tag string  `json:"tag"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := tag.RenameTag(ctx.Organization, ctx.UserId, req.Id, req.Tag); err != nil {
		return []error{err}
	}

	return nil
}

func RemoveTagHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	articleId, err := rest.GetIdParam(r, "article_id")

	if err != nil {
		return []error{err}
	}

	t := r.URL.Query().Get("tag")

	if err := tag.RemoveTag(ctx.Organization, ctx.UserId, articleId, t); err != nil {
		return []error{err}
	}

	return nil
}

func GetTagByNameHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	params := mux.Vars(r)
	rest.WriteResponse(w, tag.GetTagByIdOrName(ctx.Organization, ctx.UserId, 0, params["name"]))
	return nil
}

func ValidateTagHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	data := tag.AddTagData{Tag: rest.GetParam(r, "tag")}

	if err := tag.ValidateTag(data); err != nil {
		return []error{err}
	}

	return nil
}

func DeleteTagHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	tagId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := tag.DeleteTag(ctx.Organization, ctx.UserId, tagId); err != nil {
		return []error{err}
	}

	return nil
}
