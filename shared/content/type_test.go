package content

import (
	"testing"
)

func TestIsImage(t *testing.T) {
	input := []string{
		"",
		"wrong",
		"image",
		"image/png",
	}
	expected := []bool{
		false,
		false,
		true,
		true,
	}

	for i, in := range input {
		if IsImage(in) != expected[i] {
			t.Fatalf("Expected %v to be an image: %v", in, expected[i])
		}
	}
}
