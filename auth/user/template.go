package user

import (
	"emviwiki/shared/config"
	"emviwiki/shared/tpl"
)

const (
	registrationCompletedMail = "mail_registration_completed.html"
	registrationMail          = "mail_registration.html"
	changeEmailMail           = "mail_change_email.html"
	changePasswordMail        = "mail_change_password.html"
)

var (
	mailTplCache *tpl.Cache
)

func InitTemplates() {
	mailTplCache = tpl.NewCache(config.Get().Template.MailTemplateDir, config.Get().Template.HotReload)
}
