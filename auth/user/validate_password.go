package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
)

func ValidatePassword(id hide.ID, pwd string) error {
	user := model.GetUserById(id)

	if user == nil {
		return errs.UserNotFound
	}

	pwd = util.Sha256Base64(pwd + user.PasswordSalt.String)

	if model.GetUserByIdAndPassword(id, pwd) == nil {
		return errs.PasswordWrong
	}

	return nil
}
