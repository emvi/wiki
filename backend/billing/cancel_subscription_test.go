package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"errors"
	"github.com/stripe/stripe-go/v71"
	"strings"
	"testing"
)

func TestCancelSubscriptionByUser(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetSubscriptionResult = &stripe.Subscription{}
	client = mock
	mailProvider = mailMock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "nonadmin@test.com")
	orga.StripeSubscriptionID.SetValid("sub_id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := CancelSubscription(orga, nonAdmin.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	if err := CancelSubscription(orga, user.ID); err != nil {
		t.Fatalf("Subscription must have been cancelled, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if !orga.StripeSubscriptionID.Valid {
		t.Fatal("Stripe subscription ID must have been kept")
	}

	if !orga.SubscriptionCancelled {
		t.Fatal("Subscription must have been marked as cancelled")
	}

	if len(mock.MarkSubscriptionCancelledIDs) != 1 || mock.MarkSubscriptionCancelledIDs[0] != "sub_id" {
		t.Fatalf("Subscription must have been cancelled on stripe, but was: %v", mock.MarkSubscriptionCancelledIDs)
	}
}

func TestCancelSubscriptionByUserFailure(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetSubscriptionError = errors.New("not found")
	client = mock
	mailProvider = mailMock
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.StripeSubscriptionID.SetValid("sub_id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := CancelSubscription(orga, user.ID); err != errs.SubscriptionNotFound {
		t.Fatalf("Subscription must not have been found, but was: %v", err)
	}

	mock.GetSubscriptionError = nil
	mock.GetSubscriptionResult = &stripe.Subscription{CancelAtPeriodEnd: true}

	if err := CancelSubscription(orga, user.ID); err != errs.SubscriptionCancelled {
		t.Fatalf("Subscription have have been cancelled already, but was: %v", err)
	}
}

func TestSendSubscriptionCancelledMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	config.Get().Mail.Sender = "support-sender"
	var subject, body, to string
	mailProvider = func(mailSubject string, mailBody string, mailFrom string, mailTo ...string) error {
		subject = mailSubject
		body = mailBody
		to = mailFrom
		return nil
	}
	orga, user := testutil.CreateOrgaAndUser(t)
	sendSubscriptionCancelledMail(orga, user)

	if subject != "Your subscription has been cancelled" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if to != user.Email {
		t.Fatalf("Receiver not as expected: %v", to)
	}

	t.Log(body)

	if !strings.Contains(body, orga.Name) ||
		!strings.Contains(body, "mailto:support-sender") ||
		!strings.Contains(body, string(subscriptionCancelledMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(subscriptionCancelledMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(subscriptionCancelledMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(subscriptionCancelledMailI18n["en"]["text-3"])) ||
		!strings.Contains(body, string(subscriptionCancelledMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(subscriptionCancelledMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestCancelSubscriptionImmediately(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetSubscriptionResult = &stripe.Subscription{}
	client = mock
	mailProvider = mailMock
	orga, _ := testutil.CreateOrgaAndUser(t)
	orga.StripeSubscriptionID.SetValid("sub_id")
	orga.StripeCustomerID.SetValid("cust_id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := DeleteSubscriptionAndCustomer(orga, false); err != nil {
		t.Fatalf("Subscription must have been cancelled, but was: %v", err)
	}

	if len(mock.DeleteCustomerIDs) != 1 || mock.DeleteCustomerIDs[0] != "cust_id" {
		t.Fatalf("Customer must have been deleted on stripe, but was: %v", mock.MarkSubscriptionCancelledIDs)
	}

	if len(mock.CancelSubscriptionIDs) != 1 || mock.CancelSubscriptionIDs[0] != "sub_id" {
		t.Fatalf("Subscription must have been cancelled on stripe, but was: %v", mock.MarkSubscriptionCancelledIDs)
	}
}

func TestSendCancelSubscriptionMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	var subject, body, to string
	mailProvider = func(mailSubject string, mailBody string, mailFrom string, mailTo ...string) error {
		subject = mailSubject
		body = mailBody
		to = mailFrom
		return nil
	}
	orga, user := testutil.CreateOrgaAndUser(t)
	sendCancelSubscriptionMail(orga, user)

	if subject != "Your subscription has been cancelled" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if to != user.Email {
		t.Fatalf("Receiver not as expected: %v", to)
	}

	t.Log(body)

	if !strings.Contains(body, orga.Name) ||
		!strings.Contains(body, string(cancelSubscriptionMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(cancelSubscriptionMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(cancelSubscriptionMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(cancelSubscriptionMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(cancelSubscriptionMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func TestResetSubscription(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	client = mock
	orga, _ := testutil.CreateOrgaAndUser(t)

	if err := ResetSubscription(orga); err != nil {
		t.Fatalf("Resetting an expert organization must have no effect, but was: %v", err)
	}

	if len(mock.CancelSubscriptionIDs) != 0 {
		t.Fatal("Subscription must not have been cancelled")
	}

	orga.Expert = false
	orga.StripeCustomerID.SetValid("cust-id")
	orga.StripeSubscriptionID.SetValid("sub-id")
	orga.StripePaymentMethodID.SetValid("pm-id")
	orga.StripePaymentIntentClientSecret.SetValid("secret")
	orga.SubscriptionPlan.SetValid(yearly)
	orga.SubscriptionCancelled = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := ResetSubscription(orga); err != nil {
		t.Fatalf("Subscription must have been reset, but was: %v", err)
	}

	if len(mock.CancelSubscriptionIDs) != 1 {
		t.Fatal("Subscription must have been cancelled")
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.Expert ||
		!orga.StripeCustomerID.Valid ||
		orga.StripeSubscriptionID.Valid ||
		orga.SubscriptionPlan.Valid ||
		orga.StripePaymentMethodID.Valid ||
		orga.StripePaymentIntentClientSecret.Valid ||
		orga.SubscriptionCancelled {
		t.Fatalf("Organization not as expected after resetting subscription: %v", orga)
	}
}
