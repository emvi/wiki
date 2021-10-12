package api

import (
	"fmt"
	"strings"
)

const (
	contentEndpoint = "/api/v1/content"
)

func getResourceURL(filename string) string {
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") {
		return filename
	}

	return fmt.Sprintf("%s%s/%s", contentHost, contentEndpoint, filename)
}
