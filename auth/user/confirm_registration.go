package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"github.com/emvi/logbuch"
)

func ConfirmRegistration(code string) (int, error) {
	user := model.GetUserByRegistrationCodeAndInactive(nil, code)

	if user == nil {
		return 0, errs.UserNotFound
	}

	if user.RegistrationStep == StepInitial {
		user.RegistrationStep++

		if err := model.SaveUser(nil, user); err != nil {
			logbuch.Error("Error saving user when confirming email address", logbuch.Fields{"err": err, "user_id": user.ID})
			return 0, errs.Saving
		}
	}

	return user.RegistrationStep, nil
}
