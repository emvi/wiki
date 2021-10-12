package notification

import (
	"emviwiki/shared/config"
	"emviwiki/shared/mail"
	"emviwiki/shared/tpl"
)

var (
	mailProvider mail.Sender
	frontendHost string
	tplCache     *tpl.Cache
)

func LoadConfig() {
	mailProvider = mail.SelectMailSender()
	frontendHost = config.Get().Hosts.Frontend
	tplCache = tpl.NewCache(config.Get().Template.MailTemplateDir, false)
}
