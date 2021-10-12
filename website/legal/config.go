package legal

import (
	"emviwiki/shared/config"
)

var (
	privacyPolicyURL, cookiePolicyURL, termsAndConditionsURL, cookiesNote map[string]string
)

func LoadConfig() {
	c := config.Get()
	privacyPolicyURL = c.Legal.PrivacyPolicyURL
	cookiePolicyURL = c.Legal.CookiePolicyURL
	termsAndConditionsURL = c.Legal.TermsAndConditionsURL
	cookiesNote = c.Legal.CookiesNote
}
