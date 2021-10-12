package pages

import (
	"emviwiki/shared/config"
	"emviwiki/shared/tpl"
)

const (
	loginPageTemplate           = "login_page.html"
	logoutPageTemplate          = "logout_page.html"
	authorizationPageTemplate   = "authorization_page.html"
	clientUnknownPageTemplate   = "client_unknown_page.html"
	loginSuccessPageTemplate    = "login_success_page.html"
	passwordPageTemplate        = "password_page.html"
	passwordSuccessPageTemplate = "password_success_page.html"
	passwordResetPageTemplate   = "password_reset_page.html"
	notFoundPageTemplate        = "not_found_page.html"
	updateEmailPageTemplate     = "update_email_page.html"
	ssoErrorPageTemplate        = "sso_error_page.html"
	passwordMailTemplate        = "mail_password.html"
)

var (
	tplCache     *tpl.Cache
	mailTplCache *tpl.Cache
)

func InitTemplates() {
	c := config.Get()
	tplCache = tpl.NewCache(c.Template.TemplateDir, c.Template.HotReload)
	mailTplCache = tpl.NewCache(c.Template.MailTemplateDir, c.Template.HotReload)
}
