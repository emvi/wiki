package util

import (
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"fmt"
	"net/url"
)

func InjectSubdomain(urlStr, subdomain string) string {
	fullUrl, err := url.Parse(urlStr)

	if err != nil {
		return urlStr
	}

	fullUrl.Host = subdomain + "." + fullUrl.Host
	return fullUrl.String()
}

// GetOrganizationURL returns the URL for given path and organization.
// Example: /my/file.txt -> https://orga.emvi.com/my/file.txt
func GetOrganizationURL(orga *model.Organization, path string) string {
	urlStr := fmt.Sprintf("%s%s", config.Get().Hosts.Frontend, path)
	return InjectSubdomain(urlStr, orga.NameNormalized)
}
