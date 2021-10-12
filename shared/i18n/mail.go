package i18n

import (
	"emviwiki/shared/util"
	"html/template"
)

var (
	mailTitle = map[string]map[string]string{
		"en": {
			"password_mail":                          "Your password at Emvi has been reset",
			"registration_mail":                      "Your registration at Emvi",
			"registration_completed_mail":            "Thank you for signing up at Emvi!",
			"change_email_mail":                      "Your email address at Emvi was changed",
			"change_password_mail":                   "Your password at Emvi was changed",
			"invite":                                 "You have been invited to join an organization at Emvi",
			"joined":                                 "%s %s joined your organization at Emvi",
			"recommend_article":                      "You've got an article recommendation on Emvi",
			"invite_article":                         "You've got an invitation to edit an article on Emvi",
			"mail_notifications":                     "Your unread notifications on Emvi",
			"newsletter_confirmation_mail":           "Your newsletter subscription at Emvi",
			"newsletter_onpremise_confirmation_mail": "Your newsletter subscription at Emvi",
			"subscription":                           "Your subscription at Emvi",
			"organization_downgrade":                 "Your subscription has expired",
			"subscription_cancelled":                 "Your subscription has been cancelled",
			"resume_subscription":                    "Your subscription has been resumed",
			"cancel_subscription":                    "Your subscription has been cancelled",
			"payment_action_required":                "Your subscription at Emvi requires you to take action",
		},
		"de": {
			"password_mail":                          "Dein Passwort bei Emvi wurde zurückgesetzt",
			"registration_mail":                      "Deine Registrierung bei Emvi",
			"registration_completed_mail":            "Vielen Dank für deine Registrierung bei Emvi!",
			"change_email_mail":                      "Deine E-Mail-Adresse für Emvi wurde geändert",
			"change_password_mail":                   "Dein Passwort für Emvi wurde geändert",
			"invite":                                 "Du wurdest in eine Organisation bei Emvi eingeladen",
			"joined":                                 "Jemand ist deiner Organisation bei Emvi beigetreten",
			"recommend_article":                      "Du hast einen Lesevorschlag auf Emvi erhalten",
			"invite_article":                         "Du hast eine Einladung einen Artikel auf Emvi zu bearbeiten",
			"mail_notifications":                     "Deine ungelesenen Benachrichtigungen auf Emvi",
			"newsletter_confirmation_mail":           "Dein Newsletter Abo bei Emvi",
			"newsletter_onpremise_confirmation_mail": "Dein Newsletter Abo bei Emvi",
			"subscription":                           "Dein Abonnement bei Emvi",
			"organization_downgrade":                 "Dein Abonnement ist abgelaufen",
			"subscription_cancelled":                 "Dein Abonnement wurde beendet",
			"resume_subscription":                    "Dein Abonnement wird fortgesetzt",
			"cancel_subscription":                    "Dein Abonnement wurde beendet",
			"payment_action_required":                "Dein Abonnement bei Emvi erfordet deine Aufmerksamkeit",
		},
	}

	mailEndI18n = Translation{
		"en": {
			"end_terms":     "Terms",
			"end_privacy":   "Privacy",
			"end_legal":     "Legal",
			"end_copyright": "Copyright 2020 Emvi Software GmbH",
			"terms_url":     "/terms",
			"privacy_url":   "/privacy",
			"legal_url":     "/legal",
		},
		"de": {
			"end_terms":     "Nutzungsbedingungen",
			"end_privacy":   "Datenschutz",
			"end_legal":     "Impressum",
			"end_copyright": "Copyright 2020 Emvi Software GmbH",
			"terms_url":     "/terms",
			"privacy_url":   "/privacy",
			"legal_url":     "/legal",
		},
	}
)

// GetMailTitle returns the mail subject for given language code or the default, if not available.
func GetMailTitle(langCode string) map[string]string {
	_, ok := util.SupportedLangs[langCode]

	if ok {
		return mailTitle[langCode]
	}

	return mailTitle[util.DefaultSupportedLang]
}

// GetMailEndI18n returns the mail end variables for given language code or the default, if not available.
func GetMailEndI18n(langCode string) map[string]template.HTML {
	vars := GetVars(langCode, mailEndI18n)
	out := make(map[string]template.HTML)

	for k, v := range vars {
		out[k] = v
	}

	out["WebsiteHost"] = template.HTML(websiteHost)
	return out
}
