package newsletter

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestUnsubscribe(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailMock := func(subject, msgHTML, from string, to ...string) error {
		return nil
	}

	if err := Subscribe("test@user.com", NewsletterList, "en", mailMock); err != nil {
		t.Fatal(err)
	}

	if err := Unsubscribe(""); err != errs.NewsletterNotFound {
		t.Fatalf("Newsletter must not be found, but was: %v", err)
	}

	if err := Subscribe("test@user.com", NewsletterList, "en", mailMock); err != nil {
		t.Fatal(err)
	}

	newsletter := model.GetNewsletterSubscriptionByEmailAndList("test@user.com", NewsletterList)

	if err := Unsubscribe(newsletter.Code); err != nil {
		t.Fatalf("Newsletter must be deleted, but was: %v", err)
	}

	if model.GetNewsletterSubscriptionByEmailAndList("test@user.com", NewsletterList) != nil {
		t.Fatal("Newsletter must have been deleted")
	}
}
