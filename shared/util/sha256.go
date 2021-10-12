package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// Sha256Base64 hashes a string to SHA256 and encodes it to base64.
// If it cannot hash it, an empty string is returned.
func Sha256Base64(str string) string {
	hash := sha256.New()

	if _, err := io.WriteString(hash, str); err != nil {
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}
