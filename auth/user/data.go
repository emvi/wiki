package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"strings"
	"unicode/utf8"
)

type UserData struct {
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Language        string `json:"language"`
	AcceptMarketing bool   `json:"accept_marketing"`
}

func (data *UserData) Validate() []error {
	data.Firstname = strings.TrimSpace(data.Firstname)
	data.Lastname = strings.TrimSpace(data.Lastname)
	err := make([]error, 0)

	if utf8.RuneCountInString(data.Firstname) == 0 || utf8.RuneCountInString(data.Firstname) > firstnameMaxLength {
		err = append(err, errs.FirstnameInvalid)
	}

	if utf8.RuneCountInString(data.Lastname) == 0 || utf8.RuneCountInString(data.Lastname) > lastnameMaxLength {
		err = append(err, errs.LastnameInvalid)
	}

	if utf8.RuneCountInString(data.Language) != 0 && utf8.RuneCountInString(data.Language) != 2 {
		err = append(err, errs.LanguageInvalid)
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func ChangeUserData(userId hide.ID, data UserData) []error {
	user := model.GetUserById(userId)

	if user == nil {
		return []error{errs.UserNotFound}
	}

	if err := data.Validate(); err != nil {
		return err
	}

	user.Firstname = null.NewString(data.Firstname, true)
	user.Lastname = null.NewString(data.Lastname, true)
	user.Language = null.NewString(data.Language, data.Language != "")
	user.AcceptMarketing = data.AcceptMarketing

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when updating password", logbuch.Fields{"err": err, "user_id": userId})
		return []error{errs.Saving}
	}

	return nil
}
