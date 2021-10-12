package member

import (
	"emviwiki/shared/config"
)

var (
	websiteHost, authHost, frontendHost string
)

func LoadConfig() {
	c := config.Get()
	websiteHost = c.Hosts.Website
	authHost = c.Hosts.Auth
	frontendHost = c.Hosts.Frontend
}
