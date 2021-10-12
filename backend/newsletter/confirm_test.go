package newsletter

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestConfirm(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		return nil
	}

	if err := Subscribe("test@user.com", NewsletterList, "en", mailMock); err != nil {
		t.Fatal(err)
	}

	if err := Confirm(""); err != errs.NewsletterNotFound {
		t.Fatalf("Newsletter must not be found, but was: %v", err)
	}

	newsletter := model.GetNewsletterSubscriptionByEmailAndList("test@user.com", NewsletterList)

	if err := Confirm(newsletter.Code); err != nil {
		t.Fatalf("Newsletter must be confirmed, but was: %v", err)
	}

	// again, should work too
	if err := Confirm(newsletter.Code); err != nil {
		t.Fatalf("Newsletter must be confirmed, but was: %v", err)
	}

	newsletter = model.GetNewsletterSubscriptionByEmailAndList("test@user.com", NewsletterList)

	if !newsletter.Confirmed {
		t.Fatal("Newsletter must be confirmed")
	}
}
