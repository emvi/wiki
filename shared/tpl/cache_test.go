package tpl

import (
	"bytes"
	"testing"
)

func TestCache(t *testing.T) {
	cache := NewCache("template/*", false)
	tpl := cache.Get()
	var buffer bytes.Buffer
	data := "World!"

	if err := tpl.ExecuteTemplate(&buffer, "test.html", data); err != nil {
		t.Fatalf("Template must be executed, but was: %v", err)
	}

	if buffer.String() != "<p>Hello <span>World!</span>\n</p>\n" {
		t.Fatalf("Output not as expected: %v", buffer.String())
	}

	cache.Clear()

	if cache.loaded {
		t.Fatal("Cache must have been cleared")
	}

	cache.Disable()

	if !cache.disabled {
		t.Fatal("Cache must have been disabled")
	}

	cache.Enable()

	if cache.disabled {
		t.Fatal("Cache must have been enabled")
	}
}
