package pages

import (
	"emviwiki/shared/i18n"
)

var newsletterComponentI18n = i18n.Translation{
	"en": {
		"name@company.com":   "name@company.com",
		"your_email":         "Your email address",
		"newsletter":         "Subscribe to our newsletter",
		"subscribe":          "Subscribe",
		"newsletter_success": "Thank you for subscribing to our newsletter. Please check your inbox to verify your subscription.",
		"newsletter_error":   "Email address invalid.",
	},
	"de": {
		"name@company.com":   "name@firma.de",
		"your_email":         "Ihre E-Mail-Adresse",
		"newsletter":         "Abonnieren Sie unseren Newsletter",
		"subscribe":          "Abonnieren",
		"newsletter_success": "Vielen Dank, dass Sie sich für unseren Newsletter angemeldet haben. Bitte überprüfe Sie Ihren Posteingang, um das Abonnement zu bestätigen.",
		"newsletter_error":   "Die E-Mail-Adresse ist ungültig.",
	},
}
