package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestUpdateCustomer(t *testing.T) {
	config.Get().Stripe.TaxIDDE = "tax-id"
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetSubscriptionResult = &stripe.Subscription{
		ID: "sub-id",
		Items: &stripe.SubscriptionItemList{
			Data: []*stripe.SubscriptionItem{
				{ID: "plan-id"},
			},
		},
	}
	mock.GetCustomerResult = &stripe.Customer{
		ID: "cust-id",
		InvoiceSettings: &stripe.CustomerInvoiceSettings{
			DefaultPaymentMethod: &stripe.PaymentMethod{
				ID: "pm-id1",
			},
		},
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "non-admin@user.com")
	order := Order{
		Name:            "New Name",
		Email:           "new@mail.com",
		Country:         "DE",
		AddressLine1:    "221b Baker Street",
		PostalCode:      "NW1 6XE",
		City:            "London",
		Phone:           "01+44+207 224 3688",
		TaxNumber:       "DE123456789",
		PaymentMethodId: "pm-id2",
	}

	if err := UpdateCustomer(orga, nonAdmin.ID, order); len(err) != 1 || err[0] != errs.PermissionDenied {
		t.Fatalf("Permission must have been denied, but was: %v", err)
	}

	if err := UpdateCustomer(orga, user.ID, order); len(err) != 1 || err[0] != errs.CustomerNotFound {
		t.Fatalf("Customer must not have been found, but was: %v", err)
	}

	orga.StripeCustomerID.SetValid("cust-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := UpdateCustomer(orga, user.ID, order); err != nil {
		t.Fatalf("Customer must have been updated, but was: %v", err)
	}

	if len(mock.UpdateCustomerParams) != 1 ||
		mock.UpdateCustomerParams[0].ID != "cust-id" ||
		mock.UpdateCustomerParams[0].TaxNumber != "DE123456789" {
		t.Fatalf("Tax ID must have been updated, but was: %v", mock.UpdateSubscriptionTaxIDParams)
	}

	if len(mock.UpdateSubscriptionTaxIDParams) != 1 ||
		mock.UpdateSubscriptionTaxIDParams[0].ID != "sub-id" ||
		mock.UpdateSubscriptionTaxIDParams[0].PlanID != "plan-id" ||
		mock.UpdateSubscriptionTaxIDParams[0].TaxID != "tax-id" {
		t.Fatalf("Tax ID must have been updated, but was: %v", mock.UpdateSubscriptionTaxIDParams)
	}
}
