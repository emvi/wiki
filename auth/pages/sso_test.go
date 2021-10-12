package pages

import (
	"testing"
)

func TestSplitName(t *testing.T) {
	input := []string{"", "John", "John Doe", "John Jr Doe", "John Albert Jr Doe"}
	expected := []struct {
		firstname string
		lastname  string
	}{
		{"", ""},
		{"", "John"},
		{"John", "Doe"},
		{"John Jr", "Doe"},
		{"John Albert Jr", "Doe"},
	}

	for i, in := range input {
		firstname, lastname := splitName(in)

		if firstname != expected[i].firstname || lastname != expected[i].lastname {
			t.Fatalf("Expected %s %s, but was: %s %s", expected[i].firstname, expected[i].lastname, firstname, lastname)
		}
	}
}
