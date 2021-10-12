package user

import (
	"bytes"
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"html/template"
)

const (
	changePasswordMailTitle = "change_password_mail"
)

var changePasswordMailI18n = i18n.Translation{
	"en": {
		"title":    "Your password at Emvi was changed",
		"text":     "In case this wasn't you, please change reset your password immediately. You can do this using the password reset function on our login page.",
		"greeting": "Your password was updated.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Passwort für Emvi wurde geändert",
		"text":     "Falls diese Änderung nicht durch dich gestartet wurde, setze bitte umgehend dein Passwort zurück. Du kannst das Passwort auf der Loginseite über die Passwort vergessen Funktion zurücksetzen.",
		"greeting": "Dein Passwort wurde geändert.",
		"goodbye":  "Dein Emvi Team",
	},
}

type PasswordData struct {
	OldPwd  string `json:"old_password"`
	NewPwd1 string `json:"new_password"`
	NewPwd2 string `json:"new_password_repeat"`
}

func (data *PasswordData) validate(user *model.User) error {
	oldpwd := util.Sha256Base64(data.OldPwd + user.PasswordSalt.String)

	if oldpwd != user.Password.String {
		return errs.PasswordWrong
	}

	if data.NewPwd1 != data.NewPwd2 {
		return errs.PasswordMatch
	}

	return nil
}

func ChangePassword(userId hide.ID, data PasswordData, lang string, mail mail.Sender) error {
	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	if user.AuthProvider != emviAuthProviderName {
		return errs.OperationNotAllowed
	}

	if err := data.validate(user); err != nil {
		return err
	}

	user.Password = null.NewString(util.Sha256Base64(data.NewPwd1+user.PasswordSalt.String), true)

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when updating password", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	go sendChangedPasswordMail(user, lang, mail)
	return nil
}

func sendChangedPasswordMail(user *model.User, lang string, mail mail.Sender) {
	if user.Language.Valid {
		lang = user.Language.String
	}

	t := mailTplCache.Get()
	data := struct {
		EndVars map[string]template.HTML
		Vars    map[string]template.HTML
	}{
		i18n.GetMailEndI18n(lang),
		i18n.GetVars(lang, changePasswordMailI18n),
	}
	var buffer bytes.Buffer

	if err := t.ExecuteTemplate(&buffer, changePasswordMail, data); err != nil {
		logbuch.Error("Error executing change password mail", logbuch.Fields{"err": err})
		return
	}

	subject := i18n.GetMailTitle(lang)[changePasswordMailTitle]

	if err := mail(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending change password mail", logbuch.Fields{"err": err, "user_id": user.ID, "email": user.Email})
	}
}
