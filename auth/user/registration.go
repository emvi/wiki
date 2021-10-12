package user

import (
	"bytes"
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"html/template"
)

const (
	codeLength            = 40
	maxRegistrationMails  = 5
	registrationMailTitle = "registration_mail"
)

var registrationMailI18n = i18n.Translation{
	"en": {
		"title":    "Your registration at Emvi",
		"text":     "Thank you for your registration. To activate your account, please follow the link below.",
		"action":   "Activate account",
		"link":     "Or paste this link into your browser",
		"greeting": "Thank you for signing up!",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Deine Registrierung bei Emvi",
		"text":     "Vielen Dank für deine Registrierung bei Emvi. Um Dein Konto zu aktivieren, folge bitte dem untenstehenden Link.",
		"action":   "Konto aktivieren",
		"link":     "Oder kopiere diesen Link in deinen Browser",
		"greeting": "Danke für die Registrierung!",
		"goodbye":  "Dein Emvi Team",
	},
}

func Registration(data RegistrationData, lang string, mail mail.Sender) error {
	if err := data.validate(); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction for registration", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	user, code, err := createUser(tx, data.Email, lang)

	if err != nil {
		return err
	}

	if code != "" {
		if !sendRegistrationConfirmationMail(user.Email, code, lang, mail) {
			db.Rollback(tx)
			return errs.Mail
		}
	} else {
		if err := resendRegistrationMail(tx, user, mail); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction for registration", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	return nil
}

func createUser(tx *sqlx.Tx, email, lang string) (*model.User, string, error) {
	user := model.GetUserByEmailIgnoreActiveTx(tx, email)

	if user != nil {
		return user, "", nil
	}

	code := getRegistrationCode(tx)
	user = &model.User{Email: email,
		RegistrationCode: null.NewString(code, true),
		Language:         null.NewString(lang, true),
		RegistrationStep: StepPassword,
		AuthProvider:     emviAuthProviderName}

	if err := model.SaveUser(tx, user); err != nil {
		logbuch.Error("Error saving user on registration", logbuch.Fields{"err": err})
		return nil, "", errs.Saving
	}

	return user, code, nil
}

func getRegistrationCode(tx *sqlx.Tx) string {
	code := util.GenRandomString(codeLength)

	for model.GetUserByRegistrationCodeAndInactive(tx, code) != nil {
		code = util.GenRandomString(codeLength)
	}

	return code
}

func sendRegistrationConfirmationMail(email, code, lang string, mail mail.Sender) bool {
	t := mailTplCache.Get()
	data := struct {
		EndVars map[string]template.HTML
		Vars    map[string]template.HTML
		URL     string
		Code    string
	}{
		i18n.GetMailEndI18n(lang),
		i18n.GetVars(lang, registrationMailI18n),
		registrationConfirmationURI,
		code,
	}
	var buffer bytes.Buffer

	if err := t.ExecuteTemplate(&buffer, registrationMail, data); err != nil {
		logbuch.Error("Error executing registration mail", logbuch.Fields{"err": err})
		return false
	}

	subject := i18n.GetMailTitle(lang)[registrationMailTitle]

	if err := mail(subject, buffer.String(), email); err != nil {
		logbuch.Error("Error sending registration mail", logbuch.Fields{"err": err, "email": email})
		return false
	}

	return true
}

func resendRegistrationMail(tx *sqlx.Tx, user *model.User, mail mail.Sender) error {
	if user.Active || !user.RegistrationCode.Valid {
		db.Rollback(tx)
		return errs.UserRegisteredAlready
	}

	if user.RegistrationMailsSend >= maxRegistrationMails {
		db.Rollback(tx)
		return errs.MaxRegistrationMailsSend
	}

	if !sendRegistrationConfirmationMail(user.Email, user.RegistrationCode.String, user.Language.String, mail) {
		db.Rollback(tx)
		return errs.Mail
	}

	user.RegistrationMailsSend++

	if err := model.SaveUser(tx, user); err != nil {
		logbuch.Error("Error saving user when resending registration mail", logbuch.Fields{"err": err, "user_id": user.ID})
		return errs.Saving
	}

	return nil
}
