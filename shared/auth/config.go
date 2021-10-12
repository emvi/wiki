package auth

import (
	"emviwiki/shared/config"
)

var (
	// Host must be set to the corresponding authentication server by environment variable if used.
	Host string
)

func LoadConfig() {
	Host = config.Get().Hosts.Auth
}
