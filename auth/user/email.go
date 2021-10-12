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
	"strings"
)

const (
	emailCodeLen         = 40
	changeEmailMailTitle = "change_email_mail"
	confirmMailEndpoint  = "/api/v1/auth/user/email"
)

var changeEmailMailI18n = i18n.Translation{
	"en": {
		"title":    "Your email address at Emvi was changed",
		"text-1":   "Your email address was changed:",
		"text-2":   "To confirm this change, please follow the link below.",
		"action":   "Confirm new email address",
		"link":     "Or paste this link into your browser",
		"greeting": "Your email address was updated.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Deine E-Mail-Adresse für Emvi wurde geändert",
		"text-1":   "Deine E-Mail-Adresse wurde geändert:",
		"text-2":   "Um die Änderung zu bestätigen, folge bitte dem untenstehenden Link.",
		"action":   "Neue E-Mail-Adresse bestätigen",
		"link":     "Oder kopiere diesen Link in deinen Browser",
		"greeting": "Deine E-Mail-Adresse wurde geändert.",
		"goodbye":  "Dein Emvi Team",
	},
}

type EmailData struct {
	Email string `json:"email"`
}

func (data *EmailData) validate(user *model.User) error {
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))

	if data.Email == "" {
		return errs.EmailEmpty
	}

	if user != nil && data.Email == user.Email {
		return errs.EmailNotChanged
	}

	if !mail.EmailValid(data.Email) {
		return errs.EmailInvalid
	}

	if user != nil {
		if found := model.GetUserByEmail(data.Email); found != nil && found.ID != user.ID {
			return errs.EmailInUse
		}
	}

	return nil
}

func ChangeUserEmail(userId hide.ID, data EmailData, lang string, mail mail.Sender) error {
	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	if user.AuthProvider != emviAuthProviderName {
		return errs.OperationNotAllowed
	}

	if user.Language.Valid {
		lang = user.Language.String
	}

	if err := data.validate(user); err != nil {
		return err
	}

	code := util.GenRandomString(emailCodeLen)

	if !sendChangeEmailConfirmationMail(user, data.Email, code, lang, mail) {
		return errs.Mail
	}

	user.NewEmail = null.NewString(data.Email, true)
	user.NewEmailCode = null.NewString(code, true)

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when updating email", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	return nil
}

func sendChangeEmailConfirmationMail(user *model.User, email, code, lang string, mail mail.Sender) bool {
	t := mailTplCache.Get()
	url := authHost + confirmMailEndpoint
	data := struct {
		EndVars map[string]template.HTML
		Vars    map[string]template.HTML
		Email   string
		URL     string
		Code    string
	}{
		i18n.GetMailEndI18n(lang),
		i18n.GetVars(lang, changeEmailMailI18n),
		email,
		url,
		code,
	}
	var buffer bytes.Buffer

	if err := t.ExecuteTemplate(&buffer, changeEmailMail, data); err != nil {
		logbuch.Error("Error executing change email mail", logbuch.Fields{"err": err})
		return false
	}

	subject := i18n.GetMailTitle(lang)[changeEmailMailTitle]

	if err := mail(subject, buffer.String(), email); err != nil {
		logbuch.Error("Error sending change email mail", logbuch.Fields{"err": err, "user_id": user.ID, "email": email})
		return false
	}

	return true
}

func ConfirmUpdateEmail(email, code string) error {
	user := model.GetUserByNewEmailAndNewEmailCode(email, code)

	if user == nil {
		return errs.UserNotFound
	}

	user.Email = user.NewEmail.String
	user.NewEmail.Valid = false
	user.NewEmailCode.Valid = false

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when updating email", logbuch.Fields{"err": err})
		return errs.Saving
	}

	return nil
}
