package api

import (
	"emviwiki/backend/client"
	"emviwiki/backend/context"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"net/http"
)

const (
	headerOrg = "Organization"
)

type AuthHandler func(context.EmviContext, http.ResponseWriter, *http.Request) []error

func AuthMiddleware(next AuthHandler, getOrga, requireExpert, requireWritePermissions bool, scopes ...string) http.Handler {
	scopeList := make([]client.Scope, 0, len(scopes))

	for _, scope := range scopes {
		scopeList = append(scopeList, client.ScopeFromString(scope))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenResp, orga := AuthenticateUser(w, r, getOrga)

		if tokenResp == nil {
			return
		}

		ctx := context.NewEmviContext(orga, tokenResp.UserId, tokenResp.Scopes, tokenResp.Trusted)

		if !ctx.HasScopes(scopeList...) ||
			(requireExpert && !orga.Expert) ||
			(requireWritePermissions && !memberHasWritePermissions(orga.ID, tokenResp.UserId)) {
			rest.WriteErrorResponse(w, http.StatusForbidden)
			return
		}

		updateMemberSeen(orga, ctx.UserId)
		err := next(ctx, w, r)
		rest.HandleErrors(w, r, err...)
	})
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request, getOrga bool) (*auth.TokenResponse, *model.Organization) {
	tokenResp, err := authProvider.ValidateToken(r)

	if err != nil {
		logbuch.Debug("Error obtaining token response", logbuch.Fields{"err": err, "method": r.Method, "url": r.URL})
		w.WriteHeader(http.StatusUnauthorized)
		return nil, nil
	}

	var org *model.Organization

	if getOrga {
		org = getOrganization(r, tokenResp)

		if org == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return nil, nil
		}
	}

	return tokenResp, org
}

func getOrganization(r *http.Request, tokenResp *auth.TokenResponse) *model.Organization {
	name := r.Header.Get(headerOrg)

	// check if client has access to organization
	// the secret does not need to be checked, since this was done by the authorization server before
	if tokenResp.IsClient() {
		orga := model.GetOrganizationByNameNormalized(name)

		if orga == nil || (!tokenResp.Trusted && model.GetClientByOrganizationIdAndClientId(orga.ID, tokenResp.ClientId) == nil) {
			return nil
		}

		return orga
	}

	// return organization for user
	return model.GetOrganizationByUserIdAndNameNormalized(tokenResp.UserId, name)
}

func memberHasWritePermissions(orgaId, userId hide.ID) bool {
	if userId == 0 {
		return true
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orgaId, userId)
	return member != nil && !member.ReadOnly
}

// update the last seen date for Expert organization members
// this is necessary for the fair pricing model
func updateMemberSeen(orga *model.Organization, userId hide.ID) {
	if orga != nil && orga.Expert {
		member := model.GetOrganizationMemberByOrganizationIdAndUserIdAndLastSeenBeforeToday(orga.ID, userId)

		if member != nil {
			if err := model.UpdateOrganizationMemberLastSeenById(member.ID); err != nil {
				logbuch.Error("Error updating last seen for organization member", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "member_id": member.ID})
			}
		}
	}
}
