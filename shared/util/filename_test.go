package util

import (
	"github.com/emvi/hide"
	"testing"
)

func TestGetPictureFilenameFromId(t *testing.T) {
	a, err := GetPictureFilenameFromId(hide.ID(123))

	if err != nil {
		t.Fatal(err)
	}

	b, err := GetPictureFilenameFromId(hide.ID(321))

	if err != nil {
		t.Fatal(err)
	}

	if a == b {
		t.Fatal("File names must be unique")
	}

	if len(a) != hashidMinLength || len(b) != hashidMinLength {
		t.Fatal("File names must have fixed length")
	}
}
