package user

import (
	"bytes"
	"emviwiki/auth/errs"
	"emviwiki/auth/jwt"
	"emviwiki/auth/model"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html/template"
	"time"
)

const (
	registrationCompletedMailTitle = "registration_completed_mail"
	emviAuthProviderName           = "emvi"
)

var registrationCompleteMailI18n = i18n.Translation{
	"en": {
		"title":         "Thank you for signing up at Emvi!",
		"text":          "Your registration was completed successfully. You can now create your first organization or join an existing one.",
		"or":            "or",
		"action_create": "Create an Organization",
		"action_join":   "Join an Organization",
		"link":          "Or paste this link into your browser",
		"greeting":      "Your registration is complete.",
		"goodbye":       "Cheers, Emvi Team",
	},
	"de": {
		"title":         "Vielen Dank für deine Registrierung bei Evmi!",
		"text":          "Deine Registrierung wurde erfolgreich abgeschlossen. Du kannst nun deine erste Organisation erstellen oder einer Existierenden beitreten.",
		"or":            "oder",
		"action_create": "Erstelle eine Organisation",
		"action_join":   "Tritt einer Organisation bei",
		"link":          "Oder kopiere diesen Link in deinen Browser",
		"greeting":      "Deine Registrierung ist abgeschlossen.",
		"goodbye":       "Dein Emvi Team",
	},
}

func CompleteRegistration(data RegistrationCompletionData, lang string, mail mail.Sender) (string, time.Time, []error) {
	user, err := getUserByRegistrationCodeAndCheckStep(data.Code, StepCompletion)

	if err != nil {
		return "", time.Time{}, []error{err}
	}

	if err := data.validate(); err != nil {
		return "", time.Time{}, err
	}

	if user.RegistrationStep == StepCompletion {
		user.RegistrationStep++
	}

	user.AcceptMarketing = data.AcceptMarketing
	user.RegistrationCode.SetNil()
	user.RegistrationMailsSend = 0
	user.Active = true
	user.AuthProvider = emviAuthProviderName

	if err := model.SaveUser(nil, user); err != nil {
		logbuch.Error("Error saving user when completing registration", logbuch.Fields{"err": err, "user_id": user.ID})
		return "", time.Time{}, []error{errs.Saving}
	}

	if user.Language.Valid {
		lang = user.Language.String
	}

	token, expires, _ := loginUser(user.ID, lang)
	sendRegistrationCompletedMail(user.Email, lang, mail)
	return token, expires, nil
}

func loginUser(userId hide.ID, lang string) (string, time.Time, error) {
	token, expires, err := jwt.NewUserToken(&jwt.UserTokenClaims{UserId: userId, Language: lang})

	if err != nil {
		logbuch.Error("Error creating session when completing registration", logbuch.Fields{"err": err, "user_id": userId, "lang": lang})
		return "", time.Time{}, err
	}

	return token, expires, nil
}

func sendRegistrationCompletedMail(email, lang string, mail mail.Sender) {
	t := mailTplCache.Get()
	data := struct {
		EndVars map[string]template.HTML
		Vars    map[string]template.HTML
		URLnew  string
		URLjoin string
	}{
		i18n.GetMailEndI18n(lang),
		i18n.GetVars(lang, registrationCompleteMailI18n),
		registrationCompletedNewOrgaURI,
		registrationCompletedJoinOrgaURI,
	}
	var buffer bytes.Buffer

	if err := t.ExecuteTemplate(&buffer, registrationCompletedMail, data); err != nil {
		logbuch.Error("Error executing registration completed mail", logbuch.Fields{"err": err})
		return
	}

	subject := i18n.GetMailTitle(lang)[registrationCompletedMailTitle]

	if err := mail(subject, buffer.String(), email); err != nil {
		logbuch.Error("Error sending registration completed mail", logbuch.Fields{"err": err, "email": email})
	}
}
