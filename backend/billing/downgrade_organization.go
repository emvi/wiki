package billing

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/config"
	"emviwiki/shared/constants"
	"emviwiki/shared/i18n"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"html/template"
)

const (
	downgradeSubject = "organization_downgrade"
)

var downgradeMailI18n = i18n.Translation{
	"en": {
		"title":    "Your Expert organization has expired",
		"text-1":   "The Expert subscription plan for your organization",
		"text-2":   "has expired. But don't worry, no data was lost. Your subscription might have been stopped because you cancelled the subscription, or because we couldn't charge your account. In case of the later, please check your payment settings and start a new subscription. In the other case you can ignore this mail.",
		"text-3":   "This mail was send to all administrators of the organization.",
		"greeting": "Your Expert organization has expired.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Deine Export Organisation ist abgelaufen",
		"text-1":   "Das Expert Abonnement deiner Organisation",
		"text-2":   "ist abgelaufen. Keine Sorge, dabei sind keine Daten verloren gegangen. Dein Abonnement wurde entweder beendet, weil du das Abonnement gekündigt hast, oder weil wir deinen Account nicht belasten konnten. Im letzteren Fall überprüfe bitte die Zahlungsinformationen in den Einstellungen und beginne ein neues Abonnement. Ansonsten kannst du diese Email ignorieren.",
		"text-3":   "Diese Email wurde an alle Administratoren der Organisation gesendet.",
		"greeting": "Deine Export Organisation ist abgelaufen.",
		"goodbye":  "Dein Emvi Team",
	},
}

// Downgrade downgrades the organization for given subscription ID to Entry and sends a mail to all administrators.
func Downgrade(subscriptionId string) error {
	orga := model.GetOrganizationByStripeSubscriptionID(subscriptionId)

	if orga == nil {
		logbuch.Error("Organization for subscription not found while downgrading")
		return errs.SubscriptionNotFound
	}

	// do not delete the customer ID
	orga.Expert = false
	orga.MaxStorageGB = constants.DefaultMaxStorageGb
	orga.StripeSubscriptionID.SetNil()
	orga.StripePaymentMethodID.SetNil()
	orga.StripePaymentIntentClientSecret.SetNil()
	orga.SubscriptionPlan.SetNil()
	orga.SubscriptionCancelled = false
	orga.SubscriptionCycle.SetNil()

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while downgrading", logbuch.Fields{"err": err, "orga_id": orga.ID})
		return err
	}

	go func() {
		admins, err := getAdmins(orga)

		if err != nil {
			logbuch.Error("No administrators found while downgrading subscription", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}

		for _, admin := range admins {
			sendDowngradeMail(orga, &admin)
		}
	}()

	return nil
}

func sendDowngradeMail(orga *model.Organization, user *model.User) {
	tpl := mailtpl.Cache.Get()
	var buffer bytes.Buffer
	langCode := util.DetermineSystemSupportedLangCode(orga.ID, user.ID)
	data := struct {
		Organization *model.Organization
		OrgaURL      string
		EndVars      map[string]template.HTML
		Vars         map[string]template.HTML
	}{
		orga,
		util.InjectSubdomain(config.Get().Hosts.Frontend, orga.NameNormalized),
		i18n.GetMailEndI18n(langCode),
		i18n.GetVars(langCode, downgradeMailI18n),
	}

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.DowngradeMailTemplate, data); err != nil {
		logbuch.Error("Error executing mail template", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
		return
	}

	subject := i18n.GetMailTitle(langCode)[downgradeSubject]

	if err := mailProvider(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending downgrade mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}
}
