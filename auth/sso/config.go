package sso

import (
	"emviwiki/shared/config"
)

var (
	authHost string
)

func LoadConfig() {
	authHost = config.Get().Hosts.Auth
}
