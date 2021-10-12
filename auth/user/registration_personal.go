package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
)

func RegistrationPersonal(data RegistrationPersonalData) []error {
	user, err := getUserByRegistrationCodeAndCheckStep(data.Code, StepPersonalData)

	if err != nil {
		return []error{err}
	}

	if err := data.validate(); err != nil {
		return err
	}

	if user.RegistrationStep == StepPersonalData {
		user.RegistrationStep++
	}

	user.Firstname = null.NewString(data.Firstname, true)
	user.Lastname = null.NewString(data.Lastname, true)
	user.AuthProvider = emviAuthProviderName

	// don't overwrite setting from step 0 (initial registration) if not set
	if data.Language != "" {
		user.Language = null.NewString(data.Language, true)
	}

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when setting registration personal data", logbuch.Fields{"err": err, "user_id": user.ID})
		return []error{errs.Saving}
	}

	return nil
}
