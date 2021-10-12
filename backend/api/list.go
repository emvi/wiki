package api

import (
	"emviwiki/backend/articlelist"
	"emviwiki/backend/context"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
)

func SaveArticleListHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := articlelist.SaveArticleListData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	id, err := articlelist.SaveArticleList(ctx.Organization, ctx.UserId, req)

	if err != nil {
		return err
	}

	rest.WriteResponse(w, struct {
		Id hide.ID `json:"id"`
	}{id})
	return nil
}

func DeleteArticleListHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := articlelist.DeleteArticleList(ctx.Organization, ctx.UserId, id); err != nil {
		return []error{err}
	}

	return nil
}

func AddArticleListMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		UserIds  []hide.ID `json:"user_ids"`
		GroupIds []hide.ID `json:"group_ids"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	newMembers, err := articlelist.AddArticleListMember(ctx.Organization, ctx.UserId, listId, req.UserIds, req.GroupIds)

	if err != nil {
		return []error{err}
	}

	for i := range newMembers {
		if newMembers[i].User != nil && newMembers[i].User.Picture.Valid {
			newMembers[i].User.Picture.SetValid(getResourceURL(newMembers[i].User.Picture.String))
		}
	}

	rest.WriteResponse(w, newMembers)
	return nil
}

func RemoveArticleListMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	memberIds, err := rest.GetIdParams(r, "member_ids")

	if err != nil {
		return []error{err}
	}

	if err := articlelist.RemoveArticleListMember(ctx.Organization, ctx.UserId, listId, memberIds); err != nil {
		return []error{err}
	}

	return nil
}

func ToggleArticleListModeratorHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		MemberId hide.ID `json:"member_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := articlelist.ToggleArticleListModerator(ctx.Organization, ctx.UserId, listId, req.MemberId); err != nil {
		return []error{err}
	}

	return nil
}

func GetArticleListHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "lang")

	if err != nil {
		return []error{err}
	}

	list, moderator, observed, bookmarked, err := articlelist.ReadArticleList(ctx, langId, listId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		List       *model.ArticleList `json:"list"`
		Moderator  bool               `json:"moderator"`
		Observed   bool               `json:"observed"`
		Bookmarked bool               `json:"bookmarked"`
	}{list, moderator, observed, bookmarked})
	return nil
}

func AddArticleListEntryHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		ArticleIds []hide.ID `json:"article_ids"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	newArticles, err := articlelist.AddArticleListEntry(ctx.Organization, ctx.UserId, listId, req.ArticleIds)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, newArticles)
	return nil
}

func RemoveArticleListEntryHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	articleIds, err := rest.GetIdParams(r, "article_ids")

	if err != nil {
		return []error{err}
	}

	if err := articlelist.RemoveArticleListEntry(ctx.Organization, ctx.UserId, listId, articleIds); err != nil {
		return []error{err}
	}

	return nil
}

func GetArticleListEntriesHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	langId, err := rest.GetIdParam(r, "lang")

	if err != nil {
		return []error{err}
	}

	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	authorUserIds, err := rest.GetIdParams(r, "authors")

	if err != nil {
		return []error{err}
	}

	centerArticleId, err := rest.GetIdParam(r, "center_article_id")

	if err != nil {
		return []error{err}
	}

	centerBefore, err := rest.GetIntParam(r, "center_before")

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchArticleListEntryFilter{
		baseFilter,
		rest.GetBoolParam(r, "client_access"),
		rest.GetBoolParam(r, "archived"),
		rest.GetParam(r, "title"),
		rest.GetParam(r, "content"),
		rest.GetParam(r, "tags"),
		authorUserIds,
		rest.GetParam(r, "commits"),
		centerArticleId,
		centerBefore,
		rest.GetParam(r, "sort_position"),
		rest.GetParam(r, "sort_title"),
	}

	entries, count, pos, err := articlelist.ReadArticleListEntries(ctx, langId, listId, filter)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Entries  []model.Article `json:"entries"`
		Count    int             `json:"count"`
		StartPos int             `json:"start_pos"`
	}{entries, count, pos})
	return nil
}

func GetArticleListMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	baseFilter, err := getBaseFilter(w, r)

	if err != nil {
		return []error{err}
	}

	userIds, err := rest.GetIdParams(r, "user")

	if err != nil {
		return []error{err}
	}

	filter := &model.SearchArticleListMemberFilter{
		baseFilter,
		userIds,
		rest.GetParam(r, "sort_username"),
		rest.GetParam(r, "sort_email"),
		rest.GetParam(r, "sort_firstname"),
		rest.GetParam(r, "sort_lastname"),
	}

	member, count, err := articlelist.ReadArticleListMember(ctx.Organization, ctx.UserId, listId, filter)

	if err != nil {
		return []error{err}
	}

	for i := range member {
		if member[i].User.Picture.Valid {
			member[i].User.Picture.SetValid(getResourceURL(member[i].User.Picture.String))
		}
	}

	rest.WriteResponse(w, struct {
		Member []model.ArticleListMember `json:"member"`
		Count  int                       `json:"count"`
	}{member, count})
	return nil
}

func SortArticleListEntryHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	listId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	// either sort an entry up/down (depending on direction is positive or negative) or swap two entries (A and B)
	req := struct {
		ArticleId  hide.ID `json:"article_id"`
		Direction  int     `json:"direction"`
		ArticleIdA hide.ID `json:"article_id_b"`
		ArticleIdB hide.ID `json:"article_id_a"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if req.ArticleId != 0 {
		if err := articlelist.SortArticleListEntry(ctx.Organization, ctx.UserId, listId, req.ArticleId, req.Direction); err != nil {
			return []error{err}
		}
	} else {
		if err := articlelist.SwapArticleListEntries(ctx.Organization, ctx.UserId, listId, req.ArticleIdA, req.ArticleIdB); err != nil {
			return []error{err}
		}
	}

	return nil
}

func ReadPrivateArticleListsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	offset, err := rest.GetIntParam(r, "offset")

	if err != nil {
		return []error{err}
	}

	lists := articlelist.ReadPrivateArticleLists(ctx.Organization, ctx.UserId, offset)
	rest.WriteResponse(w, lists)
	return nil
}
