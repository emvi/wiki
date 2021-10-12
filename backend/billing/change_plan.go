package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
)

// ChangePlan changes the billing interval (aka "plan") for given organization.
// This won't do anything if the plan won't change actually (monthly -> monthly for example).
func ChangePlan(orga *model.Organization, userId hide.ID, plan string) error {
	plan = strings.TrimSpace(strings.ToLower(plan))

	if plan != yearly && plan != monthly {
		return errs.BillingIntervalInvalid
	}

	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	if !orga.StripeSubscriptionID.Valid {
		return errs.SubscriptionNotFound
	}

	if plan == orga.SubscriptionPlan.String {
		return nil
	}

	sub, err := client.GetSubscription(orga.StripeSubscriptionID.String)

	if err != nil {
		logbuch.Error("Error reading subscription while changing plan", logbuch.Fields{"orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotFound
	}

	planId := sub.Items.Data[0].ID

	if err := client.UpdateSubscriptionPrice(orga.StripeSubscriptionID.String, planId, getStripePriceID(plan)); err != nil {
		logbuch.Error("Error updating subscription plan", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.UpdatingSubscription
	}

	orga.SubscriptionPlan.SetValid(plan)

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while changing subscription plan", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.Saving
	}

	return nil
}
