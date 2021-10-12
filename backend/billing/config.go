package billing

import (
	"emviwiki/shared/mail"
	"emviwiki/shared/payment"
)

var (
	mailProvider mail.Sender
	client       payment.BillingClient
)

func LoadConfig() {
	mailProvider = mail.SelectMailSender()
	client = payment.NewStripeClient()
}
