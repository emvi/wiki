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
	resumeSubscriptionSubject = "resume_subscription"
)

var resumeSubscriptionMailI18n = i18n.Translation{
	"en": {
		"title":    "Your subscription has been resumed",
		"text-1":   "The subscription for your organization",
		"text-2":   "has been resumed. We will charge your account at the start of the next billing cycle again.",
		"greeting": "Your subscription has been resumed.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Abonnement wird fortgesetzt",
		"text-1":   "Das Abonnement für deine Organisation",
		"text-2":   "wird fortgesetzt. Wir werden dein Konto mit dem Beginn des nächsten Zahlungszeitraums wieder belasten.",
		"greeting": "Dein Abonnement wird fortgesetzt.",
		"goodbye":  "Dein Emvi Team",
	},
}

// ResumeSubscription resumes a subscription that has been cancelled before it runs out.
func ResumeSubscription(orga *model.Organization, userId hide.ID) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	if !orga.StripeSubscriptionID.Valid {
		logbuch.Error("Subscription not found while resuming subscription", logbuch.Fields{"orga_id": orga.ID, "user_id": userId, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotFound
	}

	if !orga.SubscriptionCancelled {
		logbuch.Error("Subscription not cancelled while resuming subscription", logbuch.Fields{"orga_id": orga.ID, "user_id": userId, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotCancelled
	}

	if err := client.ResumeSubscription(orga.StripeSubscriptionID.String); err != nil {
		logbuch.Error("Error resuming subscription", logbuch.Fields{"orga_id": orga.ID, "user_id": userId, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.ResumingSubscription
	}

	orga.SubscriptionCancelled = false

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while resuming subscription", logbuch.Fields{"err": err, "subscription_id": orga.StripeSubscriptionID.String, "orga_id": orga.ID, "user_id": userId})
		return errs.Saving
	}

	go func() {
		admins, err := getAdmins(orga)

		if err != nil {
			logbuch.Error("No administrators found while resuming subscription", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}

		for _, admin := range admins {
			sendResumeSubscriptionMail(orga, &admin)
		}
	}()

	return nil
}

func sendResumeSubscriptionMail(orga *model.Organization, user *model.User) {
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
		i18n.GetVars(langCode, resumeSubscriptionMailI18n),
	}

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.ResumeSubscriptionMailTemplate, data); err != nil {
		logbuch.Error("Error executing resume subscription mail template", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}

	subject := i18n.GetMailTitle(langCode)[resumeSubscriptionSubject]

	if err := mailProvider(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending resume subscription mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}
}
