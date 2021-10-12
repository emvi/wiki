package api

import (
	"emviwiki/backend/context"
	"emviwiki/backend/user"
	"emviwiki/shared/rest"
	"net/http"
)

const (
	maxUserPictureMem = 2097152 // 2 MBi
)

// Does nothing. The authentication is handled by middleware.
// This function is just a dummy to have an endpoint.
func AuthenticateUserHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	return nil
}

func GetUserHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	userData, err := authProvider.GetActiveUser(w, r)

	if err != nil {
		return []error{err}
	}

	resp, err := user.CreateOrUpdateUser(user.AuthUser{*userData})

	if err != nil {
		return []error{err}
	}

	if resp.Picture.Valid {
		resp.Picture.SetValid(getResourceURL(resp.Picture.String))
	}

	rest.WriteResponse(w, resp)
	return nil
}

func GetProfileByIdHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	userId, err := rest.IdParam(r, "id")

	if err != nil {
		return []error{err}
	}

	resp, err := user.ReadUserProfile(ctx.Organization, userId, "")

	if err != nil {
		return []error{err}
	}

	if resp.Picture.Valid {
		resp.Picture.SetValid(getResourceURL(resp.Picture.String))
	}

	rest.WriteResponse(w, resp)
	return nil
}

func GetProfileByNameHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	profile, err := user.ReadUserProfile(ctx.Organization, 0, rest.GetParam(r, "name"))

	if err != nil {
		return []error{err}
	}

	if profile.Picture.Valid {
		profile.Picture.SetValid(getResourceURL(profile.Picture.String))
	}

	rest.WriteResponse(w, profile)
	return nil
}

func UploadUserPictureHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	r.Body = http.MaxBytesReader(w, r.Body, maxUserPictureMem)

	if err := user.UploadUserPicture(r, ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func DeleteUserPictureHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	if err := user.DeleteUserPicture(ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func UpdateUserColorModeHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		ColorMode int `json:"color_mode"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.UpdateColorMode(ctx.UserId, req.ColorMode); err != nil {
		return []error{err}
	}

	return nil
}

func UpdateUserIntroductionHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Introduction bool `json:"introduction"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.UpdateIntroduction(ctx.UserId, req.Introduction); err != nil {
		return []error{err}
	}

	return nil
}
