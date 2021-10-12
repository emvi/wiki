package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"github.com/emvi/logbuch"
)

func CancelRegistration(code string) error {
	user := model.GetUserByRegistrationCodeAndInactive(nil, code)

	if user == nil {
		return errs.UserNotFound
	}

	if err := model.DeleteUserById(nil, user.ID); err != nil {
		logbuch.Error("Error deleting user when cancelling registration", logbuch.Fields{"err": err, "user_id": user.ID})
		return errs.Saving
	}

	return nil
}
