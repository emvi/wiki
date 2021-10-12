package pages

import (
	"emviwiki/shared/config"
	"emviwiki/shared/tpl"
)

const (
	landingPageTemplate  = "landing_page.html"
	pricingPageTemplate  = "pricing_page.html"
	termsPageTemplate    = "terms_page.html"
	privacyPageTemplate  = "privacy_page.html"
	cookiePageTemplate   = "cookie_page.html"
	legalPageTemplate    = "legal_page.html"
	blogPageTemplate     = "blog_page.html"
	articlePageTemplate  = "article_page.html"
	notfoundPageTemplate = "404_page.html"
)

var (
	tplCache                                                                     *tpl.Cache
	backendHost, authHost, websiteHost, clientId, version                        string
	googleSSOClientId, githubSSOClientId, slackSSOClientId, microsoftSSOClientId string
	isIntegration                                                                bool
)

func InitTemplates() {
	c := config.Get()
	backendHost = c.Hosts.Backend
	authHost = c.Hosts.Auth
	websiteHost = c.Hosts.Website
	clientId = c.AuthClient.ID
	version = c.Version
	googleSSOClientId = c.SSO.Google.ID
	githubSSOClientId = c.SSO.GitHub.ID
	slackSSOClientId = c.SSO.Slack.ID
	microsoftSSOClientId = c.SSO.Microsoft.ID
	isIntegration = c.IsIntegration
	tplCache = tpl.NewCache(c.Template.TemplateDir, c.Template.HotReload)
}
