package api

import (
	"emviwiki/shared/auth"
	"emviwiki/shared/config"
	"emviwiki/shared/mail"
)

var (
	contentHost  string
	authProvider auth.AuthClient
	mailProvider mail.Sender
)

func LoadConfig() {
	c := config.Get()
	contentHost = c.Hosts.Backend
	authProvider = auth.NewEmviAuthClient(c.AuthClient.ID, c.AuthClient.Secret)
	mailProvider = mail.SelectMailSender()
}
