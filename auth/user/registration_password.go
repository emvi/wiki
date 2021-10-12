package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
)

const (
	saltLength = 20
)

func RegistrationPassword(data RegistrationPasswordData) error {
	user, err := getUserByRegistrationCodeAndCheckStep(data.Code, StepPassword)

	if err != nil {
		return err
	}

	if err := data.validate(); err != nil {
		return err
	}

	if user.RegistrationStep == StepPassword {
		user.RegistrationStep++
	}

	user.PasswordSalt = null.NewString(util.GenRandomString(saltLength), true)
	user.Password = null.NewString(util.Sha256Base64(data.Password+user.PasswordSalt.String), true)
	user.AuthProvider = emviAuthProviderName

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when setting registration password", logbuch.Fields{"err": err, "user_id": user.ID})
		return errs.Saving
	}

	return nil
}

func getUserByRegistrationCodeAndCheckStep(code string, step int) (*model.User, error) {
	user := model.GetUserByRegistrationCodeAndInactive(nil, code)

	if user == nil {
		return nil, errs.UserNotFound
	}

	if step > user.RegistrationStep {
		return nil, errs.StepInvalid
	}

	return user, nil
}
