package pages

import (
	"emviwiki/shared/config"
	"emviwiki/shared/mail"
	"emviwiki/shared/tpl"
)

const (
	loginPageTemplate          = "login_page.html"
	startPageTemplate          = "start_page.html"
	newsletterPageTemplate     = "newsletter_page.html"
	newsletterEditPageTemplate = "newsletter_edit_page.html"
	newsletterMailTemplate     = "mail_newsletter.html"
	filesPageTemplate          = "files_page.html"
)

var (
	SecureCookies bool
	tplCache      *tpl.Cache
	mailTplCache  *tpl.Cache
	mailProvider  mail.Sender
)

func InitTemplates() {
	SecureCookies = config.Get().Server.HTTP.SecureCookies
	tplCache = tpl.NewCache(config.Get().Template.TemplateDir, config.Get().Template.HotReload)
	mailTplCache = tpl.NewCache(config.Get().Template.MailTemplateDir, config.Get().Template.HotReload)
	mailProvider = mail.SelectMailSender()
}
