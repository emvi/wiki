package newsletter

import (
	"emviwiki/shared/config"
	"emviwiki/shared/mail"
	"emviwiki/shared/tpl"
)

var (
	mailProvider             mail.Sender
	tplCache                 *tpl.Cache
	newsletterUnsubscribeURI string
)

func LoadConfig() {
	mailProvider = mail.SelectMailSender()
	tplCache = tpl.NewCache(config.Get().Template.MailTemplateDir, false)
	newsletterUnsubscribeURI = config.Get().Newsletter.UnsubscribeURI
}
