package tpl

import (
	"bytes"
	"emviwiki/shared/i18n"
	"html/template"
	"strings"
	"testing"
)

func TestRenderMail(t *testing.T) {
	cache := NewCache("../../template/mail/*.html", false)
	tpl := cache.Get()
	var buffer bytes.Buffer
	data := struct {
		EndVars map[string]template.HTML
	}{
		i18n.GetMailEndI18n("en"),
	}

	if err := tpl.ExecuteTemplate(&buffer, "mail_test.html", data); err != nil {
		t.Fatalf("Template must be executed, but was: %v", err)
	}

	content := buffer.String()

	if !strings.Contains(content, "Head Title") {
		t.Fatal("Mail must contain head title")
	}

	if !strings.Contains(content, "Preheader") {
		t.Fatal("Mail must contain preheader")
	}

	if !strings.Contains(content, "wordmark-black@1x.png") {
		t.Fatal("Mail must contain the logo")
	}

	if !strings.Contains(content, "Custom Greeting!") {
		t.Fatal("Mail must contain the greeting")
	}

	if !strings.Contains(content, "Custom Goodbye!") {
		t.Fatal("Mail must contain goodbye")
	}

	if !strings.Contains(content, "Paragraph 1") {
		t.Fatal("Mail must contain the first paragraph")
	}

	if !strings.Contains(content, "Paragraph 2") {
		t.Fatal("Mail must contain the second paragraph")
	}

	if !strings.Contains(content, "button_url") ||
		!strings.Contains(content, "Button Label") ||
		!strings.Contains(content, "Or follow this link") {
		t.Fatal("Mail must contain a button")
	}
}
