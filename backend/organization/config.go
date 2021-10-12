package organization

import (
	"emviwiki/shared/config"
	"github.com/emvi/logbuch"
	"net/url"
	"strings"
)

const (
	defaultFrontendHost = "https://emvi.com"
)

var (
	frontendHostWithoutProtocol, frontendHostProtocol string
)

func LoadConfig() {
	frontendHost := config.Get().Hosts.Frontend

	if frontendHost == "" {
		frontendHost = defaultFrontendHost
	}

	u, err := url.Parse(frontendHost)

	if err != nil {
		logbuch.Fatal("Error parsing host frontend host URL", logbuch.Fields{"err": err, "url": frontendHost})
	}

	frontendHostWithoutProtocol = u.Host
	parts := strings.Split(frontendHost, "://")

	if len(parts) < 2 || (strings.ToLower(parts[0]) != "http" && strings.ToLower(parts[0]) != "https") {
		logbuch.Fatal("Error parsing protocol from frontend host URL", logbuch.Fields{"err": err, "url": frontendHost, "parts": parts})
	}

	frontendHostProtocol = strings.ToLower(parts[0])
}
