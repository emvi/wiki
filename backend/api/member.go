package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/member"
	"emviwiki/backend/organization"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func InviteMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := member.InviteMemberData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := member.InviteMember(ctx.Organization, ctx.UserId, req, mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func JoinOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := member.JoinOrganizationData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := member.JoinOrganization(ctx.UserId, req, mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func LeaveOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Name string `json:"name"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := member.LeaveOrganization(ctx.Organization, ctx.UserId, req.Name); err != nil {
		return []error{err}
	}

	return nil
}

func GetInvitationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	params := mux.Vars(r)
	orga, err := member.GetInvitation(ctx.UserId, params["code"])

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Organization *model.Organization `json:"organization"`
	}{orga})
	return nil
}

func GetInvitationOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	orga, err := organization.GetInvitationCodeOrganization(ctx.UserId, rest.GetParam(r, "code"))

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Organization *model.Organization `json:"organization"`
	}{orga})
	return nil
}

func ReadInvitationsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	invitations, err := member.ReadInvitations(ctx.UserId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Invitations []model.Invitation `json:"invitations"`
	}{invitations})
	return nil
}

func ReadOrganizationInvitationsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	invitations, err := member.ReadOrganizationInvitations(ctx.Organization, ctx.UserId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, invitations)
	return nil
}

func CancelInvitationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	invitationId, err := rest.GetIdParam(r, "invitation_id")

	if err != nil {
		return []error{err}
	}

	if err := member.CancelInvitation(ctx.Organization, ctx.UserId, invitationId); err != nil {
		return []error{err}
	}

	return nil
}

func DeleteInvitationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	id, err := rest.GetIdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := member.DeleteInvitation(ctx.UserId, id); err != nil {
		return []error{err}
	}

	return nil
}

func ToggleMemberModeratorHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	userId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := member.ToggleModerator(ctx.Organization, ctx.UserId, userId); err != nil {
		return []error{err}
	}

	return nil
}

func ToggleMemberAdminHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	userId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := member.ToggleAdmin(ctx.Organization, ctx.UserId, userId); err != nil {
		return []error{err}
	}

	return nil
}

func ToggleReadOnlyHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	userId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	if err := member.ToggleReadOnly(ctx.Organization, ctx.UserId, userId); err != nil {
		return []error{err}
	}

	return nil
}

func RemoveMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	userId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	removePermissions := rest.GetBoolParam(r, "remove_permissions")

	if err := member.RemoveMember(ctx.Organization, ctx.UserId, userId, removePermissions); err != nil {
		return []error{err}
	}

	return nil
}

func GetMemberHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	orgaMember, err := member.GetMember(ctx.Organization, ctx.UserId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, orgaMember)
	return nil
}

func SaveSettingsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := member.Settings{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := member.SaveSettings(ctx.Organization, ctx.UserId, req); err != nil {
		return err
	}

	return nil
}
