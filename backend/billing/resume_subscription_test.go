package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"strings"
	"testing"
)

func TestResumeSubscription(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	client = mock
	mailProvider = mailMock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "nonadmin@test.com")

	if err := ResumeSubscription(orga, nonAdmin.ID); err != errs.PermissionDenied {
		t.Fatalf("Permission must be denied, but was: %v", err)
	}

	if err := ResumeSubscription(orga, user.ID); err != errs.SubscriptionNotFound {
		t.Fatalf("Subscription must not have been found, but was: %v", err)
	}

	orga.StripeSubscriptionID.SetValid("sub_id")
	orga.SubscriptionCancelled = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := ResumeSubscription(orga, user.ID); err != nil {
		t.Fatalf("Subscription must have been resumed, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if !orga.StripeSubscriptionID.Valid || orga.SubscriptionCancelled {
		t.Fatalf("Subscription ID must still be set and the subscription must not be marked as canelled anymore, but was: %v %v", orga.StripeSubscriptionID.String, orga.SubscriptionCancelled)
	}

	if len(mock.ResumeSubscriptionIDs) != 1 || mock.ResumeSubscriptionIDs[0] != "sub_id" {
		t.Fatalf("Subscription must have been resumed, but was: %v", mock.ResumeSubscriptionIDs)
	}
}

func TestSendResumeSubscriptionMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	var subject, body, to string
	mailProvider = func(mailSubject string, mailBody string, mailFrom string, mailTo ...string) error {
		subject = mailSubject
		body = mailBody
		to = mailFrom
		return nil
	}
	orga, user := testutil.CreateOrgaAndUser(t)
	sendResumeSubscriptionMail(orga, user)

	if subject != "Your subscription has been resumed" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if to != user.Email {
		t.Fatalf("Receiver not as expected: %v", to)
	}

	t.Log(body)

	if !strings.Contains(body, orga.Name) ||
		!strings.Contains(body, string(resumeSubscriptionMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(resumeSubscriptionMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(resumeSubscriptionMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(resumeSubscriptionMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(resumeSubscriptionMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}
