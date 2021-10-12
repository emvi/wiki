package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func UpdatePaymentMethod(orga *model.Organization, userId hide.ID, paymentMethodId string) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	if !orga.StripeCustomerID.Valid {
		return errs.CustomerNotFound
	}

	customer, err := client.GetCustomer(orga.StripeCustomerID.String)

	if err != nil {
		logbuch.Error("Error getting existing customer while updating payment method", logbuch.Fields{"orga_id": orga.ID, "customer_id": orga.StripeCustomerID.String})
		return errs.CustomerNotFound
	}

	if err := client.DetachPaymentMethod(customer.InvoiceSettings.DefaultPaymentMethod.ID); err != nil {
		logbuch.Error("Error detaching payment method", logbuch.Fields{"orga_id": orga.ID, "customer_id": customer.ID})
		return errs.UpdatingSubscription
	}

	pm, err := client.AttachPaymentMethod(customer, paymentMethodId)

	if err != nil {
		logbuch.Error("Error attaching payment method", logbuch.Fields{"orga_id": orga.ID, "customer_id": customer.ID})
		return errs.UpdatingSubscription
	}

	if _, err := client.UpdateDefaultPaymentMethod(customer, paymentMethodId); err != nil {
		logbuch.Error("Error updating default payment method", logbuch.Fields{"orga_id": orga.ID, "customer_id": customer.ID})
		return errs.UpdatingSubscription
	}

	orga.StripePaymentMethodID.SetValid(pm.ID)

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error updating organization while updating payment method", logbuch.Fields{"orga_id": orga.ID, "customer_id": customer.ID})
		return errs.Saving
	}

	return nil
}
