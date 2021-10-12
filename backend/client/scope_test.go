package client

import "testing"

func TestScope_String(t *testing.T) {
	scope := Scope{"name", false, false}

	if scope.String() != "name" {
		t.Fatalf("Expected 'name', but was: %v", scope.String())
	}

	scope.Read = true

	if scope.String() != "name:r" {
		t.Fatalf("Expected 'name:r', but was: %v", scope.String())
	}

	scope.Write = true

	if scope.String() != "name:rw" {
		t.Fatalf("Expected 'name:rw', but was: %v", scope.String())
	}
}

func TestScopeFromString(t *testing.T) {
	input := []string{
		"",
		"invalid",
		"articles:test",
		"articles:r",
		"articles:rw",
	}
	expected := []Scope{
		{"", false, false},
		{"", false, false},
		{"articles", false, false},
		{"articles", true, false},
		{"articles", true, true},
	}

	for i, in := range input {
		scope := ScopeFromString(in)

		if scope.Name != expected[i].Name || scope.Read != expected[i].Read || scope.Write != expected[i].Write {
			t.Fatalf("Expected '%v' %v %v, but was: %v %v %v", expected[i].Name, expected[i].Read, expected[i].Write, scope.Name, scope.Read, scope.Write)
		}
	}
}
