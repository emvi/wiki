package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/usergroup"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"net/http"
)

func SaveUserGroupHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := usergroup.SaveUserGroupData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	id, err := usergroup.SaveUserGroup(ctx.Organization, ctx.UserId, req)

	if err != nil {
		return err
	}

	rest.WriteResponse(w, &struct {
		Id hide.ID `json:"id"`
	}{id})
	return nil
}

func DeleteUserGroupHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := usergroup.DeleteUserGroup(ctx.Organization, ctx.UserId, id); err != nil {
		return []error{err}
	}

	return nil
}

func AddUserGroupMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	groupId, err := rest.IdParam(r, "id")

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

	newMember, err := usergroup.AddUserGroupMember(ctx.Organization, ctx.UserId, groupId, req.UserIds, req.GroupIds)

	if err != nil {
		return []error{err}
	}

	for i := range newMember {
		if newMember[i].User.Picture.Valid {
			newMember[i].User.Picture.SetValid(getResourceURL(newMember[i].User.Picture.String))
		}
	}

	rest.WriteResponse(w, newMember)
	return nil
}

func RemoveUserGroupMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	groupId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	memberIds, err := rest.GetIdParams(r, "member_ids")

	if err != nil {
		return []error{err}
	}

	if err := usergroup.RemoveUserGroupMember(ctx.Organization, ctx.UserId, groupId, memberIds); err != nil {
		return []error{err}
	}

	return nil
}

func ToggleUserGroupModeratorHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	groupId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	req := struct {
		MemberId hide.ID `json:"member_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := usergroup.ToggleUserGroupModerator(ctx.Organization, ctx.UserId, groupId, req.MemberId); err != nil {
		return []error{err}
	}

	return nil
}

func GetUserGroupHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	groupId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	group, moderator, observed, err := usergroup.ReadUserGroup(ctx.Organization, ctx.UserId, groupId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Group     *model.UserGroup `json:"group"`
		Moderator bool             `json:"moderator"`
		Observed  bool             `json:"observed"`
	}{group, moderator, observed})
	return nil
}

func GetUserGroupMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	groupId, err := rest.IdParam(r, "id")

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

	filter := &model.SearchUserGroupMemberFilter{
		baseFilter,
		userIds,
		rest.GetParam(r, "sort_username"),
		rest.GetParam(r, "sort_email"),
		rest.GetParam(r, "sort_firstname"),
		rest.GetParam(r, "sort_lastname"),
	}

	member, count, err := usergroup.ReadUserGroupMember(ctx.Organization, groupId, filter)

	if err != nil {
		return []error{err}
	}

	for i := range member {
		if member[i].User.Picture.Valid {
			member[i].User.Picture.SetValid(getResourceURL(member[i].User.Picture.String))
		}
	}

	rest.WriteResponse(w, struct {
		Member []model.UserGroupMember `json:"member"`
		Count  int                     `json:"count"`
	}{member, count})
	return nil
}
