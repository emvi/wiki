package newsletter

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
)

func Unsubscribe(code string) error {
	newsletter := model.GetNewsletterSubscriptionByCode(code)

	if newsletter == nil {
		return errs.NewsletterNotFound
	}

	if err := model.DeleteNewsletterSubscriptionById(nil, newsletter.ID); err != nil {
		logbuch.Error("Error deleting newsletter while unsubscribing", logbuch.Fields{"err": err, "code": code})
		return errs.Saving
	}

	return nil
}
