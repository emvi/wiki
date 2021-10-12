package util

import (
	"testing"
)

func TestSha256Base64(t *testing.T) {
	if len(Sha256Base64("string")) != 64 {
		t.Fatal("SHA256 string must be 64 chars long")
	}
}
