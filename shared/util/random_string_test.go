package util

import (
	"testing"
)

func TestGenRandomString(t *testing.T) {
	if len(GenRandomString(20)) != 20 {
		t.Fatal("Random string must be 20 chars long")
	}

	if GenRandomString(10) == GenRandomString(10) {
		t.Fatal("Two generated strings must not be equal")
	}
}
