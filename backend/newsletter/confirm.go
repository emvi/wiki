package newsletter

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
)

func Confirm(code string) error {
	newsletter := model.GetNewsletterSubscriptionByCode(code)

	if newsletter == nil {
		return errs.NewsletterNotFound
	}

	if newsletter.Confirmed {
		return nil
	}

	newsletter.Confirmed = true

	if err := model.SaveNewsletterSubscription(nil, newsletter); err != nil {
		logbuch.Error("Error saving newsletter when confirming email address", logbuch.Fields{"err": err, "code": code})
		return errs.Saving
	}

	return nil
}
