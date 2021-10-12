package balance

import (
	"emviwiki/shared/payment"
)

var (
	client payment.BillingClient
)

func LoadConfig() {
	client = payment.NewStripeClient()
}
