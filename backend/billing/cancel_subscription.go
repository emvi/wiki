package billing

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/perm"
	"emviwiki/shared/config"
	"emviwiki/shared/i18n"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html/template"
)

const (
	subscriptionCancelledSubject = "subscription_cancelled"
	cancelSubscriptionSubject    = "cancel_subscription"
)

var subscriptionCancelledMailI18n = i18n.Translation{
	"en": {
		"title":    "Your subscription has been cancelled",
		"text-1":   "The subscription for your organization",
		"text-2":   "has been cancelled. We won't charge your account anymore. Your organization will be downgraded at the end of the billing cycle.",
		"text-3":   "Thank you for using Emvi. In case you would like to tell use why you cancelled your subscription, please send a mail to:",
		"greeting": "Your subscription has been cancelled.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Abonnement wurde beendet",
		"text-1":   "Das Abonnement für deine Organisation",
		"text-2":   "wurde beendet. Wir werden dein Konto nicht weiter belasten. Deine Organisation wird am Ende des Rechnungszeitraums herabgestuft.",
		"text-3":   "Danke, dass du Emvi verwendest. Wenn du uns mitteilen möchtest, warum du dein Abonnement beendet hast, sende doch bitte eine E-Mail an:",
		"greeting": "Dein Abonnement wurde beendet.",
		"goodbye":  "Dein Emvi Team",
	},
}

var cancelSubscriptionMailI18n = i18n.Translation{
	"en": {
		"title":    "Your subscription has been cancelled",
		"text-1":   "The subscription for your organization",
		"text-2":   "has been cancelled. We won't charge your account anymore. Thank you for using Emvi.",
		"greeting": "Your subscription has been cancelled.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Abonnement wurde beendet",
		"text-1":   "Das Abonnement für deine Organisation",
		"text-2":   "wurde beendet. Wir werden dein Konto nicht weiter belasten. Danke, dass du Emvi verwendest.",
		"greeting": "Dein Abonnement wurde beendet.",
		"goodbye":  "Dein Emvi Team",
	},
}

// MarkSubscriptionCancelled cancels the subscription for given organization.
// This will just mark the subscription as cancelled at Stripe, but won't remove it from the organization,
// as the subscription might have been payed for longer.
func CancelSubscription(orga *model.Organization, userId hide.ID) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	if !orga.StripeSubscriptionID.Valid {
		logbuch.Error("Subscription not found while cancelling subscription", logbuch.Fields{"orga_id": orga.ID, "user_id": userId, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotFound
	}

	sub, err := client.GetSubscription(orga.StripeSubscriptionID.String)

	if err != nil {
		logbuch.Error("Subscription not found while cancelling subscription", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotFound
	}

	if sub.CancelAtPeriodEnd {
		logbuch.Error("Subscription has been been cancelled already", logbuch.Fields{"orga_id": orga.ID, "user_id": userId, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionCancelled
	}

	if err := client.MarkSubscriptionCancelled(orga.StripeSubscriptionID.String); err != nil {
		logbuch.Error("Error cancelling subscription", logbuch.Fields{"err": err, "subscription_id": orga.StripeSubscriptionID.String, "orga_id": orga.ID, "user_id": userId})
		return errs.CancellingSubscription
	}

	orga.SubscriptionCancelled = true

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while cancelling subscription", logbuch.Fields{"err": err, "subscription_id": orga.StripeSubscriptionID.String, "orga_id": orga.ID, "user_id": userId})
		return errs.Saving
	}

	go func() {
		admins, err := getAdmins(orga)

		if err != nil {
			logbuch.Error("No administrators found while cancelling subscription", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}

		for _, admin := range admins {
			sendSubscriptionCancelledMail(orga, &admin)
		}
	}()

	return nil
}

func sendSubscriptionCancelledMail(orga *model.Organization, user *model.User) {
	tpl := mailtpl.Cache.Get()
	var buffer bytes.Buffer
	langCode := util.DetermineSystemSupportedLangCode(orga.ID, user.ID)
	data := struct {
		Organization *model.Organization
		SupportMail  string
		EndVars      map[string]template.HTML
		Vars         map[string]template.HTML
	}{
		orga,
		config.Get().Mail.Sender,
		i18n.GetMailEndI18n(langCode),
		i18n.GetVars(langCode, subscriptionCancelledMailI18n),
	}

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.SubscriptionCancelledMailTemplate, data); err != nil {
		logbuch.Error("Error executing cancel subscription mail template", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}

	subject := i18n.GetMailTitle(langCode)[subscriptionCancelledSubject]

	if err := mailProvider(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending cancel subscription mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}
}

// DeleteSubscriptionAndCustomer will cancel the subscription immediately and delete the customer.
// This function must only be triggered by the system. Sending mails to administrators is optional.
func DeleteSubscriptionAndCustomer(orga *model.Organization, sendMails bool) error {
	subCancelled := false

	if orga.StripeSubscriptionID.Valid {
		subCancelled = true

		if err := client.CancelSubscription(orga.StripeSubscriptionID.String); err != nil {
			logbuch.Error("Error cancelling subscription immediately", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String, "customer_id": orga.StripeCustomerID.String})
			return errs.CancellingSubscription
		}
	}

	if orga.StripeCustomerID.Valid {
		if err := client.DeleteCustomer(orga.StripeCustomerID.String); err != nil {
			logbuch.Error("Error deleting customer while cancelling subscription immediately", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String, "customer_id": orga.StripeCustomerID.String})
			return errs.DeletingCustomer
		}
	}

	if subCancelled && sendMails {
		go func() {
			admins, err := getAdmins(orga)

			if err != nil {
				logbuch.Error("No administrators found while cancelling subscription immediately", logbuch.Fields{"err": err, "orga_id": orga.ID})
			}

			for _, admin := range admins {
				sendCancelSubscriptionMail(orga, &admin)
			}
		}()
	}

	return nil
}

func sendCancelSubscriptionMail(orga *model.Organization, user *model.User) {
	tpl := mailtpl.Cache.Get()
	var buffer bytes.Buffer
	langCode := util.DetermineSystemSupportedLangCode(orga.ID, user.ID)
	data := struct {
		Organization *model.Organization
		EndVars      map[string]template.HTML
		Vars         map[string]template.HTML
	}{
		orga,
		i18n.GetMailEndI18n(langCode),
		i18n.GetVars(langCode, cancelSubscriptionMailI18n),
	}

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.CancelSubscriptionMailTemplate, data); err != nil {
		logbuch.Error("Error executing cancel subscription mail template", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}

	subject := i18n.GetMailTitle(langCode)[cancelSubscriptionSubject]

	if err := mailProvider(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending cancel subscription mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}
}

// ResetSubscription resets an incomplete subscription before the organization has been upgraded.
// This must only be called by the system and not be triggered by the user.
func ResetSubscription(orga *model.Organization) error {
	if orga.Expert {
		return nil
	}

	if orga.StripeSubscriptionID.Valid {
		if err := client.CancelSubscription(orga.StripeSubscriptionID.String); err != nil {
			logbuch.Error("Error resetting subscription", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String, "customer_id": orga.StripeCustomerID.String})
			return errs.CancellingSubscription
		}
	}

	orga.StripeSubscriptionID.SetNil()
	orga.StripePaymentMethodID.SetNil()
	orga.StripePaymentIntentClientSecret.SetNil()
	orga.SubscriptionPlan.SetNil()
	orga.SubscriptionCancelled = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while resetting subscription", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String, "customer_id": orga.StripeCustomerID.String})
		return errs.Saving
	}

	return nil
}
