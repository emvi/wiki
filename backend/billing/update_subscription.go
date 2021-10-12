package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
)

// UpdateSubscription updates the subscription quantity and limits for given organization.
// This must be called if the member count changed. It won't do anything if the organization is non-expert.
func UpdateSubscription(orga *model.Organization) error {
	if !orga.StripeSubscriptionID.Valid {
		return nil
	}

	member := model.CountOrganizationMemberByOrganizationIdAndActiveAndNotReadOnly(orga.ID)
	orga.MaxStorageGB = int64(constants.StorageGBPerUser * member)

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while upgrading", logbuch.Fields{"err": err, "orga_id": orga.ID})
		return errs.Saving
	}

	sub, err := client.GetSubscription(orga.StripeSubscriptionID.String)

	if err != nil {
		logbuch.Error("Subscription not found while updating subscription", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotFound
	}

	planId := sub.Items.Data[0].ID

	if err := client.UpdateSubscriptionQuantity(sub.ID, planId, int64(member)); err != nil {
		logbuch.Error("Error updating subscription quantity", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.UpdatingSubscription
	}

	logbuch.Debug("Next invoice", logbuch.Fields{"timestamp": sub.NextPendingInvoiceItemInvoice})
	return nil
}
