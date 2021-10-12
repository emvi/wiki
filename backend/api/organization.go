package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/organization"
	"emviwiki/shared/rest"
	"net/http"
)

const (
	maxOrganizationPictureMem = 2097152 // 2 MBi
)

func GetOrganizationsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	orgas := organization.ReadOrganizations(ctx.UserId)

	for i := range orgas {
		if orgas[i].Picture.Valid {
			orgas[i].Picture.SetValid(getResourceURL(orgas[i].Picture.String))
		}
	}

	rest.WriteResponse(w, orgas)
	return nil
}

func GetOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	if ctx.Organization != nil {
		if ctx.Organization.Picture.Valid {
			ctx.Organization.Picture.SetValid(getResourceURL(ctx.Organization.Picture.String))
		}

		rest.WriteResponse(w, ctx.Organization)
		return nil
	}

	orga, err := organization.ReadOrganization(ctx)

	if err != nil {
		return []error{err}
	}

	if orga.Picture.Valid {
		orga.Picture.SetValid(getResourceURL(orga.Picture.String))
	}

	rest.WriteResponse(w, orga)
	return nil
}

func CreateOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := organization.CreateOrganizationData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if errs := organization.CreateOrganization(ctx.UserId, req); errs != nil {
		return errs
	}

	return nil
}

func ValidateOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		organization.CreateOrganizationData

		ValidateNameDomain   bool `json:"validate_name_domain"`
		ValidateUsernameLang bool `json:"validate_username_lang"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	req.Trim()
	errs := make([]error, 0)

	if req.ValidateNameDomain {
		if err := req.CheckDomainValid(req.Domain); err != nil {
			errs = append(errs, err)
		}

		if err := req.CheckNameValid(req.Name); err != nil {
			errs = append(errs, err)
		}
	} else if req.ValidateUsernameLang {
		if err := req.CheckUsernameValid(req.Username); err != nil {
			errs = append(errs, err)
		}

		if err := req.CheckLanguage(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func DeleteOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	name := r.URL.Query().Get("name")

	if err := organization.DeleteOrganization(ctx.Organization.ID, ctx.UserId, name, authProvider); err != nil {
		return []error{err}
	}

	return nil
}

func UploadOrganizationPictureHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	r.Body = http.MaxBytesReader(w, r.Body, maxOrganizationPictureMem)

	if err := organization.UploadOrganizationPicture(r, ctx.Organization.ID, ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func DeleteOrganizationPictureHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	if err := organization.DeleteOrganizationPicture(ctx.Organization.ID, ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func UpdateOrganizationHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := organization.UpdateOrganization(ctx.Organization, ctx.UserId, req.Name, req.Domain); err != nil {
		return err
	}

	return nil
}

func UpdateOrganizationPermissionsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		CreateGroupAdmin bool `json:"create_group_admin"`
		CreateGroupMod   bool `json:"create_group_mod"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := organization.UpdateOrganizationPermissions(ctx.Organization, ctx.UserId, req.CreateGroupAdmin, req.CreateGroupMod); err != nil {
		return []error{err}
	}

	return nil
}

func GetOrganizationStatisticsHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	statistics, err := organization.GetOrganizationStatistics(ctx.Organization, ctx.UserId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, statistics)
	return nil
}

func GenerateInvitationCodeHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		ReadOnly bool `json:"read_only"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	code, err := organization.GenerateInvitationCode(ctx.Organization, ctx.UserId, req.ReadOnly)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Code string `json:"code"`
	}{code})
	return nil
}

func GetInvitationCodeHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	code, readOnly, err := organization.GetInvitationCode(ctx.Organization, ctx.UserId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Code     string `json:"code"`
		ReadOnly bool   `json:"read_only"`
	}{code, readOnly})
	return nil
}
