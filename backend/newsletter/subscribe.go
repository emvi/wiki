package newsletter

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"html/template"
	"strings"
	"unicode/utf8"
)

const (
	NewsletterList                           = "" // empty so the default will be null in database
	NewsletterOnPremiseList                  = "onpremise"
	maxEmailLength                           = 255
	codeLength                               = 20
	newsletterConfirmationMailTitle          = "newsletter_confirmation_mail"
	newsletterOnPremiseConfirmationMailTitle = "newsletter_onpremise_confirmation_mail"
)

var newsletterConfirmationMailI18n = i18n.Translation{
	"en": {
		"title":          "Your newsletter subscription at Emvi",
		"text":           "To confirm your email address, please follow the link bellow.",
		"text_confirmed": "Your email address was confirmed already, you don't need to do anything.",
		"action":         "Confirm your email address",
		"link":           "Or paste this link into your browser",
		"cancel":         "If you didn't subscribe to our newsletter or would like to cancel your subscription, please follow this link:",
		"greeting":       "Thank you for subscribing to our newsletter!",
		"goodbye":        "Cheers, Emvi Team",
	},
	"de": {
		"title":          "Dein Newsletter Abo bei Emvi",
		"text":           "Um deine E-Mail-Adresse zu bestätigen, folge bitte dem untenstehenden Link.",
		"text_confirmed": "Deine E-Mail-Adresse war bereits bestätigt, du brauchst nichts weiter zu tun.",
		"action":         "Bestätige deine E-Mail-Adresse",
		"link":           "Oder kopiere diesen Link in deinen Browser",
		"cancel":         "Wenn du keinen Newsletter bestellt hast oder dein Abo beenden möchtest, folge bitte diesem Link:",
		"greeting":       "Danke für dein Newsletter Abo bei Emvi!",
		"goodbye":        "Dein Emvi Team",
	},
}

var newsletterOnpremiseConfirmationMailI18n = i18n.Translation{
	"en": {
		"title":          "Your on-premises newsletter subscription at Emvi",
		"text":           "We will notify you when we have updates on our on-premises solution only. To confirm your email address, please follow the link bellow.",
		"text_confirmed": "We will notify you when we have updates on our on-premises solution only. Your email address was confirmed already, you don't need to do anything.",
		"action":         "Confirm your email address",
		"link":           "Or paste this link into your browser",
		"cancel":         "If you didn't subscribe to our newsletter or would like to cancel your subscription, please follow this link:",
		"greeting":       "Thank you for your newsletter subscription at Emvi!",
		"goodbye":        "Cheers, Emvi Team",
	},
	"de": {
		"title":          "Dein On-Premises Newsletter Abo bei Emvi",
		"text":           "Wir werden dich ausschließlich bei Neuigkeiten zu unserer On-Premises-Lösung benachrichtigen. Um deine E-Mail-Adresse zu bestätigen, folge bitte dem untenstehenden Link.",
		"text_confirmed": "Wir werden dich ausschließlich bei Neuigkeiten zu unserer On-Premises-Lösung benachrichtigen. Deine E-Mail-Adresse war bereits bestätigt, du brauchst nichts weiter zu tun.",
		"action":         "Bestätige deine E-Mail-Adresse",
		"link":           "Oder kopiere diesen Link in deinen Browser",
		"cancel":         "Wenn du keinen Newsletter bestellt hast oder dein Abo beenden möchtest, folge bitte diesem Link:",
		"greeting":       "Danke für dein Newsletter Abo bei Emvi!",
		"goodbye":        "Dein Emvi Team",
	},
}

func Subscribe(email, list, lang string, mail mail.Sender) error {
	email = strings.TrimSpace(email)

	if err := validEmail(email); err != nil {
		return err
	}

	newsletter := model.GetNewsletterSubscriptionByEmailAndList(email, list)

	if newsletter == nil {
		newsletter = &model.NewsletterSubscription{Email: email,
			List: null.NewString(list, list != ""),
			Code: getUniqueCode()}
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction while registering to newsletter", logbuch.Fields{"err": err, "email": email, "list": list})
		return errs.TxBegin
	}

	if err := model.SaveNewsletterSubscription(tx, newsletter); err != nil {
		logbuch.Error("Error saving new newsletter", logbuch.Fields{"err": err, "email": email, "list": list})
		return errs.Saving
	}

	templateName := mailtpl.NewsletterConfirmationMailTemplate
	subjectConfig := newsletterConfirmationMailTitle

	if list == NewsletterOnPremiseList {
		templateName = mailtpl.NewsletterOnPremiseConfirmationMailTemplate
		subjectConfig = newsletterOnPremiseConfirmationMailTitle
	}

	if !sendNewsletterConfirmationMail(templateName, subjectConfig, email, newsletter.Code, lang, newsletter.Confirmed, mail) {
		db.Rollback(tx)
		logbuch.Error("Error sending newsletter confirmation mail", logbuch.Fields{"email": email, "list": list})
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction while registering to newsletter", logbuch.Fields{"err": err, "email": email, "list": list})
		return errs.TxCommit
	}

	return nil
}

func validEmail(email string) error {
	if email == "" || utf8.RuneCountInString(email) > maxEmailLength || !mail.EmailValid(email) {
		return errs.EmailInvalid
	}

	return nil
}

func getUniqueCode() string {
	code := util.GenRandomString(codeLength)

	for model.GetNewsletterSubscriptionByCode(code) != nil {
		code = util.GenRandomString(codeLength)
	}

	return code
}

func sendNewsletterConfirmationMail(templateName, subjectConfig, email, code, langCode string, confirmed bool, mail mail.Sender) bool {
	vars := newsletterConfirmationMailI18n

	if templateName == mailtpl.NewsletterOnPremiseConfirmationMailTemplate {
		vars = newsletterOnpremiseConfirmationMailI18n
	}

	t := mailtpl.Cache.Get()
	data := struct {
		EndVars    map[string]template.HTML
		Vars       map[string]template.HTML
		URLconfirm string
		URLunsub   string
		Code       string
		Confirmed  bool
	}{
		i18n.GetMailEndI18n(langCode),
		i18n.GetVars(langCode, vars),
		newsletterConfirmationURI,
		newsletterUnsubscribeURI,
		code,
		confirmed,
	}
	var buffer bytes.Buffer

	if err := t.ExecuteTemplate(&buffer, templateName, data); err != nil {
		logbuch.Error("Error executing newsletter confirmation mail", logbuch.Fields{"err": err})
		return false
	}

	subject := i18n.GetMailTitle(langCode)[subjectConfig]

	if err := mail(subject, buffer.String(), email); err != nil {
		logbuch.Error("Error sending newsletter confirmation mail", logbuch.Fields{"err": err, "email": email})
		return false
	}

	return true
}
