package article

import (
	"emviwiki/shared/config"
)

var (
	frontendHost string
)

func LoadConfig() {
	frontendHost = config.Get().Hosts.Frontend
}
