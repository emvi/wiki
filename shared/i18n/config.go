package i18n

import "emviwiki/shared/config"

var (
	websiteHost string
)

func LoadConfig() {
	websiteHost = config.Get().Hosts.Website
}
