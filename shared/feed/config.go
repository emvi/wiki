package feed

import (
	"emviwiki/shared/config"
	"github.com/gosimple/slug"
)

var (
	websiteHost, authHost, frontendHost string
)

func LoadConfig() {
	c := config.Get()
	websiteHost = c.Hosts.Website
	authHost = c.Hosts.Auth
	frontendHost = c.Hosts.Frontend

	// keep uppercase characters in slugs
	slug.Lowercase = false
}
