package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"strings"
	"testing"
	"time"
)

func TestDowngrade(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailProvider = mailMock
	orga, _ := testutil.CreateOrgaAndUser(t)
	orga.Expert = true
	orga.MaxStorageGB = 123
	orga.SubscriptionPlan.SetValid(monthly)
	orga.SubscriptionCancelled = true
	orga.StripeSubscriptionID.SetValid("sub-id")
	orga.StripeCustomerID.SetValid("cust-id")
	orga.StripePaymentMethodID.SetValid("pm-id")
	orga.StripePaymentIntentClientSecret.SetValid("secret")
	orga.SubscriptionCycle.SetValid(time.Now())

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := Downgrade("not-found"); err != errs.SubscriptionNotFound {
		t.Fatalf("Subscription must not have been found, but was: %v", err)
	}

	if err := Downgrade("sub-id"); err != nil {
		t.Fatalf("Organization must have been downgraded, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.Expert ||
		orga.MaxStorageGB != constants.DefaultMaxStorageGb ||
		orga.SubscriptionPlan.Valid ||
		orga.SubscriptionCancelled ||
		orga.StripeSubscriptionID.Valid ||
		orga.StripePaymentMethodID.Valid ||
		orga.StripePaymentIntentClientSecret.Valid ||
		!orga.StripeCustomerID.Valid ||
		orga.SubscriptionCycle.Valid {
		t.Fatalf("Organization not as expected: %v", orga)
	}
}

func TestDowngradeMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	var subject, body, to string
	mailProvider = func(mailSubject string, mailBody string, mailFrom string, mailTo ...string) error {
		subject = mailSubject
		body = mailBody
		to = mailFrom
		return nil
	}
	orga, user := testutil.CreateOrgaAndUser(t)
	sendDowngradeMail(orga, user)

	if subject != "Your subscription has expired" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if to != user.Email {
		t.Fatalf("Receiver not as expected: %v", to)
	}

	t.Log(body)

	if !strings.Contains(body, orga.Name) ||
		!strings.Contains(body, string(downgradeMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(downgradeMailI18n["en"]["goodbye"])) ||
		!strings.Contains(body, string(downgradeMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(downgradeMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(downgradeMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(downgradeMailI18n["en"]["text-3"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
