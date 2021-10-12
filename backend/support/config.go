package support

import "emviwiki/shared/config"

var (
	supportMailAddress string
)

func LoadConfig() {
	supportMailAddress = config.Get().Mail.SupportSender
}
