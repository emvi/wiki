package user

import (
	"emviwiki/backend/content"
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"net/http"
)

const (
	userPicturePath = "user/pictures"
)

func UploadUserPicture(r *http.Request, userId hide.ID) error {
	// read user
	user := model.GetUserById(userId)

	if user == nil {
		return errs.PermissionDenied
	}

	if err := deleteUserPicture(user); err != nil {
		logbuch.Debug("Error deleting old user picture", logbuch.Fields{"err": err})
	}

	uniqueName, err := content.UploadFile(&content.File{
		Request:       r,
		UserId:        userId,
		Path:          userPicturePath,
		RequiresImage: true,
	})

	if err != nil {
		return errs.UploadingFile
	}

	user.Picture = null.NewString(uniqueName, true)

	if err := model.SaveUser(nil, user, false); err != nil {
		logbuch.Error("Error saving user when uploading picture", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	return nil
}

func DeleteUserPicture(userId hide.ID) error {
	user := model.GetUserById(userId)

	if user == nil {
		return errs.PermissionDenied
	}

	if err := deleteUserPicture(user); err != nil {
		return err
	}

	return nil
}

func deleteUserPicture(user *model.User) error {
	if user.Picture.Valid {
		if err := content.DeleteFile(nil, user.ID, user.Picture.String); err != nil {
			logbuch.Error("Error deleting user picture", logbuch.Fields{"err": err, "user_id": user.ID})
		}

		// save user
		user.Picture = null.String{}

		if err := model.SaveUser(nil, user, false); err != nil {
			logbuch.Error("Error saving user when deleting picture", logbuch.Fields{"err": err, "user_id": user.ID})
			return errs.Saving
		}
	}

	return nil
}
