package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

// UpdateCustomer updates the customer and tax ID if required. It uses the same data as Subscribe.
// The payment method won't be changed.
func UpdateCustomer(orga *model.Organization, userId hide.ID, order Order) []error {
	if err := order.validate(); err != nil && len(err) > 1 && err[0] != errs.BillingIntervalInvalid {
		return err
	}

	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return []error{err}
	}

	if !orga.StripeCustomerID.Valid {
		return []error{errs.CustomerNotFound}
	}

	customer, err := client.GetCustomer(orga.StripeCustomerID.String)

	if err != nil {
		logbuch.Error("Error getting existing customer while updating customer details", logbuch.Fields{"orga_id": orga.ID, "customer_id": orga.StripeCustomerID.String})
		return []error{errs.CustomerNotFound}
	}

	_, err = client.UpdateCustomer(customer.ID,
		order.Email,
		order.Name,
		order.Country,
		order.AddressLine1,
		order.AddressLine2,
		order.PostalCode,
		order.City,
		order.Phone,
		order.TaxNumber,
		getTaxIdType(order.Country),
		getTaxExempt(order.Country))

	if err != nil {
		logbuch.Error("Error updating existing customer while updating customer details", logbuch.Fields{"err": err, "orga_id": orga.ID, "customer_id": customer.ID})
		return []error{errs.UpdatingCustomer}
	}

	if customer.Address.Country != order.Country ||
		getCustomerTaxNumber(customer) != order.TaxNumber {
		if err := updateTaxId(orga, order); err != nil {
			return []error{err}
		}
	}

	return nil
}

func updateTaxId(orga *model.Organization, order Order) error {
	sub, err := client.GetSubscription(orga.StripeSubscriptionID.String)

	if err != nil {
		logbuch.Error("Subscription not found to update tax ID", logbuch.Fields{"orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return errs.SubscriptionNotFound
	}

	planId := sub.Items.Data[0].ID
	taxId := getTaxId(order.Country, order.TaxNumber)
	logbuch.Debug("Using tax ID for update", logbuch.Fields{"orga_id": orga.ID, "country": order.Country, "tax_number": order.TaxNumber, "tax_id": taxId})

	if (taxId == "" && len(sub.Items.Data[0].TaxRates) != 0) ||
		(taxId != "" && (len(sub.Items.Data[0].TaxRates) == 0 || sub.Items.Data[0].TaxRates[0].ID != taxId)) {
		if err := client.UpdateSubscriptionTaxID(sub.ID, planId, taxId); err != nil {
			logbuch.Error("Error updating tax ID", logbuch.Fields{"orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String, "tax_id": taxId})
			return errs.UpdatingSubscription
		}
	}

	return nil
}
