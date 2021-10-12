package newsletter

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"strings"
	"testing"
)

func TestSubscribe(t *testing.T) {
	testutil.CleanBackendDb(t)
	input := []struct {
		email    string
		list     string
		langCode string
	}{
		{"", NewsletterList, "en"},
		{"   ", NewsletterList, "en"},
		{"invalid", NewsletterList, "en"},
		{"0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345", NewsletterList, "en"},
		{"okay@test.com", NewsletterList, "en"},
		{"okay@test.com", NewsletterList, "en"}, // resend confirmation
		{"okay@test.com", NewsletterOnPremiseList, "en"},
		{"different@test.com", NewsletterOnPremiseList, "en"},
	}
	expected := []error{
		errs.EmailInvalid,
		errs.EmailInvalid,
		errs.EmailInvalid,
		errs.EmailInvalid,
		nil,
		nil,
		nil,
		nil,
	}
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		return nil
	}

	for i, in := range input {
		if err := Subscribe(in.email, in.list, in.langCode, mailMock); err != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}
}

func TestSubscribeResend(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailCount := 0
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		mailCount++
		return nil
	}

	if err := Subscribe("test@user.com", NewsletterList, "de", mailMock); err != nil {
		t.Fatalf("Mail must have been send, but was: %v", err)
	}

	if err := Subscribe("test@user.com", NewsletterList, "de", mailMock); err != nil {
		t.Fatalf("Mail must have been send, but was: %v", err)
	}

	if mailCount != 2 {
		t.Fatalf("Two mails must have been send, but was: %v", mailCount)
	}

	n := model.GetNewsletterSubscriptionByEmailAndList("test@user.com", NewsletterList)

	if n == nil || n.Confirmed || len(n.Code) != 20 || n.List.Valid {
		t.Fatalf("Newsletter not as expected: %v", n)
	}
}

func TestSubscribeResendOnPremise(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailCount := 0
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		mailCount++
		return nil
	}

	if err := Subscribe("test@user.com", NewsletterOnPremiseList, "de", mailMock); err != nil {
		t.Fatalf("Mail must have been send, but was: %v", err)
	}

	if err := Subscribe("test@user.com", NewsletterOnPremiseList, "de", mailMock); err != nil {
		t.Fatalf("Mail must have been send, but was: %v", err)
	}

	if mailCount != 2 {
		t.Fatalf("Two mails must have been send, but was: %v", mailCount)
	}

	n := model.GetNewsletterSubscriptionByEmailAndList("test@user.com", NewsletterOnPremiseList)

	if n == nil || n.Confirmed || len(n.Code) != 20 || !n.List.Valid || n.List.String != NewsletterOnPremiseList {
		t.Fatalf("Newsletter not as expected: %v", n)
	}
}

func TestSendNewsletterConfirmationMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	newsletterConfirmationURI = "confirm-url"
	newsletterUnsubscribeURI = "unsub-url"
	var subject, body string
	mailMock := func(mailSubject, mailBody, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if !sendNewsletterConfirmationMail(mailtpl.NewsletterConfirmationMailTemplate, newsletterConfirmationMailTitle, "new@sub.com", "sub-code", "en", false, mailMock) {
		t.Fatal("Mail must have been send")
	}

	if subject != "Your newsletter subscription at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, newsletterConfirmationURI) ||
		!strings.Contains(body, newsletterUnsubscribeURI) ||
		!strings.Contains(body, "sub-code") ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["text"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["action"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["link"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["cancel"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestSendNewsletterConfirmationMailConfirmed(t *testing.T) {
	testutil.CleanBackendDb(t)
	newsletterConfirmationURI = "confirm-url"
	newsletterUnsubscribeURI = "unsub-url"
	var subject, body string
	mailMock := func(mailSubject, mailBody, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if !sendNewsletterConfirmationMail(mailtpl.NewsletterConfirmationMailTemplate, newsletterConfirmationMailTitle, "new@sub.com", "sub-code", "en", true, mailMock) {
		t.Fatal("Mail must have been send")
	}

	if subject != "Your newsletter subscription at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if strings.Contains(body, newsletterConfirmationURI) ||
		!strings.Contains(body, newsletterUnsubscribeURI) ||
		!strings.Contains(body, "sub-code") ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["text_confirmed"])) ||
		strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["action"])) ||
		strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["link"])) ||
		!strings.Contains(body, string(newsletterConfirmationMailI18n["en"]["cancel"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestSendNewsletterConfirmationMailOnPremise(t *testing.T) {
	testutil.CleanBackendDb(t)
	newsletterConfirmationURI = "confirm-url"
	newsletterUnsubscribeURI = "unsub-url"
	var subject, body string
	mailMock := func(mailSubject, mailBody, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if !sendNewsletterConfirmationMail(mailtpl.NewsletterOnPremiseConfirmationMailTemplate, newsletterOnPremiseConfirmationMailTitle, "new@sub.com", "sub-code", "en", false, mailMock) {
		t.Fatal("Mail must have been send")
	}

	if subject != "Your newsletter subscription at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if !strings.Contains(body, newsletterConfirmationURI) ||
		!strings.Contains(body, newsletterUnsubscribeURI) ||
		!strings.Contains(body, "sub-code") ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["text"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["action"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["link"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["cancel"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestSendNewsletterConfirmationMailOnPremiseConfirmed(t *testing.T) {
	testutil.CleanBackendDb(t)
	config.Get().Newsletter.ConfirmationURI = "confirm-url"
	config.Get().Newsletter.UnsubscribeURI = "unsub-url"
	LoadConfig()
	var subject, body string
	mailMock := func(mailSubject, mailBody, from string, to ...string) error {
		subject = mailSubject
		body = mailBody
		return nil
	}

	if !sendNewsletterConfirmationMail(mailtpl.NewsletterOnPremiseConfirmationMailTemplate, newsletterOnPremiseConfirmationMailTitle, "new@sub.com", "sub-code", "en", true, mailMock) {
		t.Fatal("Mail must have been send")
	}

	if subject != "Your newsletter subscription at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	t.Log(body)

	if strings.Contains(body, newsletterConfirmationURI) ||
		!strings.Contains(body, newsletterUnsubscribeURI) ||
		!strings.Contains(body, "sub-code") ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["text_confirmed"])) ||
		strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["action"])) ||
		strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["link"])) ||
		!strings.Contains(body, string(newsletterOnpremiseConfirmationMailI18n["en"]["cancel"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
