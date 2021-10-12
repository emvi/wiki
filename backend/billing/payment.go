package billing

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/perm"
	"emviwiki/shared/i18n"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html/template"
)

const (
	paymentActionRequiredSubject = "payment_action_required"
)

var paymentActionRequiredMailI18n = i18n.Translation{
	"en": {
		"title":    "Your subscription at Emvi requires you to take action",
		"text-1":   "The subscription for your organization",
		"text-2":   "requires you to take action. Unfortunately we could not charge your account. Please go to the billing page of your organization and confirm your payment details.",
		"greeting": "Your subscription at Emvi requires you to take action",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Abonnement bei Emvi erfordert deine Aufmerksamkeit",
		"text-1":   "Das Abonnement für deine Organisation",
		"text-2":   "erfordert deine Aufmerksamkeit. Leider konnten wir deinen Account nicht belasten. Bitte rufe die Rechnungsseite deiner Organization auf und bestätige deine Zahlung.",
		"greeting": "Dein Abonnement bei Emvi erfordert deine Aufmerksamkeit",
		"goodbye":  "Dein Emvi Team",
	},
}

// PaymentActionRequired saves the payment intent client secret for later to be confirmed by the customer.
// It won't be saved if the organization has not been upgraded yet (on first subscription).
func PaymentActionRequired(orga *model.Organization, paymentIntentId string) error {
	// it's the first time we're charging the customer if the organization has not been upgraded yet
	// don't do anything in that case
	if !orga.Expert {
		return nil
	}

	pm, err := client.GetPaymentIntent(paymentIntentId)

	if err != nil {
		logbuch.Error("Payment intent not found while saving payment intent client secret", logbuch.Fields{"err": err, "orga_id": orga.ID, "pm_id": paymentIntentId})
		return errs.PaymentIntentNotFound
	}

	orga.StripePaymentIntentClientSecret.SetValid(pm.ClientSecret)

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error updating organization while saving payment intent client secret", logbuch.Fields{"err": err, "orga_id": orga.ID, "pm_id": paymentIntentId, "client_secret": pm.ClientSecret})
		return errs.Saving
	}

	go func() {
		admins, err := getAdmins(orga)

		if err != nil {
			logbuch.Error("No administrators found while sending payment action required mail", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}

		for _, admin := range admins {
			sendPaymentActionRequiredMail(orga, &admin)
		}
	}()

	return nil
}

func sendPaymentActionRequiredMail(orga *model.Organization, user *model.User) {
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
		i18n.GetVars(langCode, paymentActionRequiredMailI18n),
	}

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.PaymentActionRequiredMailTemplate, data); err != nil {
		logbuch.Error("Error executing payment action required mail template", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}

	subject := i18n.GetMailTitle(langCode)[paymentActionRequiredSubject]

	if err := mailProvider(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending payment action required mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}
}

// RemovePaymentIntentClientSecret removes the payment intent client secret.
func RemovePaymentIntentClientSecret(orga *model.Organization, userId hide.ID) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	if !orga.StripePaymentIntentClientSecret.Valid {
		return nil
	}

	orga.StripePaymentIntentClientSecret.SetNil()

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error updating organization while removing payment intent client secret", logbuch.Fields{"err": err, "orga_id": orga.ID})
		return errs.Saving
	}

	return nil
}
