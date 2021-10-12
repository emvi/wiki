package recaptcha

import (
	"emviwiki/shared/config"
	"github.com/emvi/logbuch"
)

var (
	recaptchaSecret string
)

func LoadConfig() {
	recaptchaSecret = config.Get().RecaptchaClientSecret

	if recaptchaSecret == "" {
		logbuch.Fatal("recaptcha secret must be configured")
	}
}
