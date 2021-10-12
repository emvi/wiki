package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func UpdateColorMode(userId hide.ID, colorMode int) error {
	if colorMode < 0 || colorMode > 2 {
		return errs.ColorModeInvalid
	}

	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	user.ColorMode = colorMode

	if err := model.SaveUser(nil, user, false); err != nil {
		logbuch.Error("Error saving user while changing color mode", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	return nil
}
