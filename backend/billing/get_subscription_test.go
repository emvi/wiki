package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestGetPaymentMethod(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetCustomerResult = &stripe.Customer{
		ID: "cust-id",
		InvoiceSettings: &stripe.CustomerInvoiceSettings{
			DefaultPaymentMethod: &stripe.PaymentMethod{
				ID: "pm-id",
			},
		},
		Address: stripe.Address{
			Country:    "DE",
			Line1:      "line 1",
			Line2:      "line 2",
			PostalCode: "postal code",
			City:       "city",
		},
		Phone: "123456",
		TaxIDs: &stripe.TaxIDList{
			Data: []*stripe.TaxID{
				{Type: stripe.TaxIDTypeEUVAT, Value: "DE123456789"},
			},
		},
	}
	mock.GetPaymentMethodResult = &stripe.PaymentMethod{
		ID: "pm-id",
		BillingDetails: &stripe.BillingDetails{
			Email: "billing@mail.com",
			Name:  "Billing Name",
		},
		Card: &stripe.PaymentMethodCard{
			Brand:    stripe.PaymentMethodCardBrandVisa,
			ExpMonth: 7,
			ExpYear:  2023,
			Last4:    "1234",
		},
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "non-admin@user.com")

	if _, _, err := GetSubscription(orga, nonAdmin.ID); err != errs.PermissionDenied {
		t.Fatalf("Access must have been denied, but was: %v", err)
	}

	if _, _, err := GetSubscription(orga, user.ID); err != errs.CustomerNotFound {
		t.Fatalf("Customer must not have been found, but was: %v", err)
	}

	orga.StripeCustomerID.SetValid("cust-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	customer, _, err := GetSubscription(orga, user.ID)

	if err != nil {
		t.Fatalf("Subscription must have been returned, but was: %v", err)
	}

	if len(mock.GetPaymentMethodIDs) != 1 || mock.GetPaymentMethodIDs[0] != "pm-id" {
		t.Fatalf("Subscription call not as expected: %v", mock.GetPaymentMethodIDs)
	}

	if customer.ID != "cust-id" ||
		customer.Country != "DE" ||
		customer.AddressLine1 != "line 1" ||
		customer.AddressLine2 != "line 2" ||
		customer.PostalCode != "postal code" ||
		customer.City != "city" ||
		customer.Phone != "123456" ||
		customer.TaxNumber != "DE123456789" {
		t.Fatalf("Customer not as expected: %v", customer)
	}
}
