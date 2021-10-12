package content

import (
	"strings"
)

const (
	imageMimeType = "image"
)

func IsImage(mimeType string) bool {
	parts := strings.Split(mimeType, "/")

	if len(parts) == 0 {
		return false
	}

	return strings.ToLower(parts[0]) == imageMimeType
}
