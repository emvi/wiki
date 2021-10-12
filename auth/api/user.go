package api

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/auth/user"
	"emviwiki/shared/rest"
	"errors"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultUpdateEmailUri = "/auth/email"
)

type UserDataResponse struct {
	Id              hide.ID     `json:"id"`
	Email           string      `json:"email"`
	Firstname       null.String `json:"firstname"`
	Lastname        null.String `json:"lastname"`
	Language        *string     `json:"language"`
	Picture         string      `json:"picture"`
	AcceptMarketing bool        `json:"accept_marketing"`
	Active          bool        `json:"active"`
	Created         time.Time   `json:"created"`
	Updated         time.Time   `json:"updated"`
	IsSSOUser       bool        `json:"is_sso_user"`
}

func UpdateUserEmailHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	req := user.EmailData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.ChangeUserEmail(ctx.UserId, req, rest.GetSupportedLangCode(r), mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func UpdateUserEmailConfirmationHandler(w http.ResponseWriter, r *http.Request) []error {
	email := strings.ToLower(rest.GetParam(r, "email"))
	code := rest.GetParam(r, "code")
	redirectURI := rest.GetParam(r, "redirect_uri")

	if redirectURI == "" {
		redirectURI = defaultUpdateEmailUri
	}

	if err := user.ConfirmUpdateEmail(email, code); err != nil {
		redirectURL, urlerr := url.Parse(redirectURI)

		if urlerr != nil {
			logbuch.Error("Error parsing redirect URI on update email confirmation", logbuch.Fields{"err": urlerr})
			return []error{errors.New("Error parsing redirect URI")}
		}

		query := redirectURL.Query()
		query.Add("error", err.Error())
		redirectURL.RawQuery = query.Encode()
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
		return nil
	}

	http.Redirect(w, r, redirectURI, http.StatusFound)
	return nil
}

func UpdateUserDataHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	req := user.UserData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.ChangeUserData(ctx.UserId, req); err != nil {
		return err
	}

	return nil
}

func UpdatePasswordHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	req := user.PasswordData{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := user.ChangePassword(ctx.UserId, req, rest.GetSupportedLangCode(r), mailProvider); err != nil {
		return []error{err}
	}

	return nil
}

func GetUserHandler(ctx *AuthContext, w http.ResponseWriter, r *http.Request) []error {
	userEntity := model.GetUserById(ctx.UserId)

	if userEntity == nil {
		return []error{errs.UserNotFound}
	}

	rest.WriteResponse(w, struct {
		TokenResponse
		User *UserDataResponse `json:"user"`
	}{
		TokenResponse: ctx.TokenResponse,
		User:          userToJSON(userEntity, ctx.TokenResponse),
	})
	return nil
}

func userToJSON(user *model.User, tokenResponse TokenResponse) *UserDataResponse {
	var lang *string

	if user.Language.Valid {
		lang = &user.Language.String
	}

	data := &UserDataResponse{Id: user.ID,
		Email:           user.Email,
		Firstname:       user.Firstname,
		Lastname:        user.Lastname,
		Language:        lang,
		Picture:         user.PictureURL.String,
		AcceptMarketing: user.AcceptMarketing,
		Active:          user.Active,
		Created:         user.DefTime,
		Updated:         user.ModTime,
		IsSSOUser:       tokenResponse.IsSSOUser}

	return data
}
