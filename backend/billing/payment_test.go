package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"strings"
	"testing"
)

func TestPaymentActionRequired(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailProvider = mailMock
	mock := payment.NewMockClient()
	mock.GetPaymentIntentResult = &stripe.PaymentIntent{
		ClientSecret: "secret",
	}
	client = mock
	orga, _ := testutil.CreateOrgaAndUser(t)

	if err := PaymentActionRequired(orga, "pi-id"); err != nil {
		t.Fatalf("Client secret must have been saved, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.StripePaymentIntentClientSecret.String != "secret" {
		t.Fatalf("Client secret not as expected: %v", orga.StripePaymentIntentClientSecret.String)
	}
}

func TestSendPaymentActionRequiredMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	var subject, body, to string
	mailProvider = func(mailSubject string, mailBody string, mailFrom string, mailTo ...string) error {
		subject = mailSubject
		body = mailBody
		to = mailFrom
		return nil
	}
	orga, user := testutil.CreateOrgaAndUser(t)
	sendPaymentActionRequiredMail(orga, user)

	if subject != "Your subscription at Emvi requires you to take action" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if to != user.Email {
		t.Fatalf("Receiver not as expected: %v", to)
	}

	t.Log(body)

	if !strings.Contains(body, orga.Name) ||
		!strings.Contains(body, string(paymentActionRequiredMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(paymentActionRequiredMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(paymentActionRequiredMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(paymentActionRequiredMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(paymentActionRequiredMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestRemovePaymentIntentClientSecret(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "non-admin@user.com")

	if err := RemovePaymentIntentClientSecret(orga, nonAdmin.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must have been denied, but was: %v", err)
	}

	if err := RemovePaymentIntentClientSecret(orga, user.ID); err != nil {
		t.Fatalf("Nothing must have been changed if client secret is not set, but was: %v", err)
	}

	orga.StripePaymentIntentClientSecret.SetValid("secret")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := RemovePaymentIntentClientSecret(orga, user.ID); err != nil {
		t.Fatalf("Client secret must have been removed, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.StripePaymentIntentClientSecret.Valid {
		t.Fatalf("Client secret must have been removed, but was: %v", orga.StripePaymentIntentClientSecret.String)
	}
}
