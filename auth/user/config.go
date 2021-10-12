package user

import (
	"emviwiki/shared/config"
	"emviwiki/shared/recaptcha"
)

var (
	authHost                         string
	registrationConfirmationURI      string
	registrationCompletedNewOrgaURI  string
	registrationCompletedJoinOrgaURI string
	recaptchaValidator               recaptcha.Recaptcha
)

func LoadConfig() {
	c := config.Get()
	authHost = c.Hosts.Auth
	registrationConfirmationURI = c.Registration.ConfirmationURI
	registrationCompletedNewOrgaURI = c.Registration.CompletedNewOrgaURI
	registrationCompletedJoinOrgaURI = c.Registration.CompletedJoinOrgaURI
	recaptchaValidator = recaptcha.NewRecaptchaValidator()
}
