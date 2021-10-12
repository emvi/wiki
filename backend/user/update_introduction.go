package user

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func UpdateIntroduction(userId hide.ID, introduction bool) error {
	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	user.Introduction = introduction

	if err := model.SaveUser(nil, user, false); err != nil {
		logbuch.Error("Error saving user while updating introduction", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	return nil
}
