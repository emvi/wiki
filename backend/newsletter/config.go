package newsletter

import (
	"emviwiki/shared/config"
)

var (
	newsletterConfirmationURI, newsletterUnsubscribeURI string
)

func LoadConfig() {
	newsletterConfirmationURI = config.Get().Newsletter.ConfirmationURI
	newsletterUnsubscribeURI = config.Get().Newsletter.UnsubscribeURI
}
