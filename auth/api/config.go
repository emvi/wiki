package api

import (
	"emviwiki/shared/mail"
)

var (
	mailProvider mail.Sender
)

func LoadConfig() {
	mailProvider = mail.SelectMailSender()
}
