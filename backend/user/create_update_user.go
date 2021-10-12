package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/auth"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"strings"
)

type AuthUser struct {
	auth.UserResponse
}

func (data *AuthUser) validate() error {
	data.Email = strings.TrimSpace(data.Email)
	data.Firstname = strings.TrimSpace(data.Firstname)
	data.Lastname = strings.TrimSpace(data.Lastname)
	data.Info = strings.TrimSpace(data.Info)

	if data.Id == 0 || data.Email == "" || data.Firstname == "" || data.Lastname == "" {
		return errs.UserDataInvalid
	}

	return nil
}

func CreateOrUpdateUser(data AuthUser) (*model.User, error) {
	if err := data.validate(); err != nil {
		return nil, err
	}

	user := model.GetUserById(data.Id)
	isnew := user == nil

	if isnew || userRequiresUpdate(data, user) {
		if isnew {
			user = new(model.User)
			user.Picture = null.NewString(data.Picture, data.Picture != "")
			user.Introduction = true
		}

		logbuch.Debug("Updating user with response from authentication server", logbuch.Fields{"id": user.ID})
		user.ID = data.Id
		user.Email = data.Email
		user.Firstname = data.Firstname
		user.Lastname = data.Lastname
		user.Language = null.NewString(data.Language, data.Language != "")
		user.Info = null.NewString(data.Info, data.Info != "")
		user.AcceptMarketing = data.AcceptMarketing

		if err := model.SaveUser(nil, user, isnew); err != nil {
			return nil, errs.Saving
		}
	}

	user.IsSSOUser = data.IsSSOUser
	return user, nil
}

func userRequiresUpdate(data AuthUser, user *model.User) bool {
	return data.Id != user.ID ||
		data.Email != user.Email ||
		data.Firstname != user.Firstname ||
		data.Lastname != user.Lastname ||
		data.Language != user.Language.String ||
		data.Info != user.Info.String ||
		data.AcceptMarketing != user.AcceptMarketing
}
