package pages

import (
	"emviwiki/auth/sso"
	"emviwiki/shared/config"
	"emviwiki/shared/mail"
)

var (
	mailProvider                                                                 mail.Sender
	githubSSOProvider, slackSSOProvider, googleSSOProvider, microsoftSSOProvider sso.SSOProvider
	githubSSOClientId, slackSSOClientId, googleSSOClientId, microsoftSSOClientId string
	authHost, websiteHost                                                        string
)

func LoadConfig() {
	c := config.Get()
	mailProvider = mail.SelectMailSender()
	githubSSOClientId = c.SSO.GitHub.ID
	slackSSOClientId = c.SSO.Slack.ID
	googleSSOClientId = c.SSO.Google.ID
	microsoftSSOClientId = c.SSO.Microsoft.ID
	authHost = c.Hosts.Auth
	websiteHost = c.Hosts.Website
	githubSSOProvider = sso.NewGitHubSSOProvider(githubSSOClientId, c.SSO.GitHub.Secret)
	slackSSOProvider = sso.NewSlackSSOProvider(slackSSOClientId, c.SSO.Slack.Secret)
	googleSSOProvider = sso.NewGoogleSSOProvider(googleSSOClientId, c.SSO.Google.Secret)
	microsoftSSOProvider = sso.NewMicrosoftSSOProvider(microsoftSSOClientId, c.SSO.Microsoft.Secret)
}
