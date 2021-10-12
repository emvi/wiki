package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestUpdatePaymentMethod(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetCustomerResult = &stripe.Customer{
		ID: "cust-id",
		InvoiceSettings: &stripe.CustomerInvoiceSettings{
			DefaultPaymentMethod: &stripe.PaymentMethod{
				ID: "pm-id1",
			},
		},
	}
	mock.AttachPaymentMethodResult = &stripe.PaymentMethod{
		ID: "pm-id",
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "non-admin@user.com")

	if err := UpdatePaymentMethod(orga, nonAdmin.ID, ""); err != errs.PermissionDenied {
		t.Fatalf("Permission must have been denied, but was: %v", err)
	}

	if err := UpdatePaymentMethod(orga, user.ID, ""); err != errs.CustomerNotFound {
		t.Fatalf("Customer must not be found, but was: %v", err)
	}

	orga.StripeCustomerID.SetValid("cust-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := UpdatePaymentMethod(orga, user.ID, "pm-id2"); err != nil {
		t.Fatalf("Payment method must have been updated, but was: %v", err)
	}

	if len(mock.GetCustomerIDs) != 1 || mock.GetCustomerIDs[0] != "cust-id" {
		t.Fatalf("Customer request not as expected: %v", mock.GetCustomerIDs)
	}

	if len(mock.DetachPaymentMethodIDs) != 1 || mock.DetachPaymentMethodIDs[0] != "pm-id1" {
		t.Fatalf("Detach payment method request not as expected: %v", mock.DetachPaymentMethodIDs)
	}

	if len(mock.AttachPaymentMethodParams) != 1 ||
		mock.AttachPaymentMethodParams[0].PaymentMethodID != "pm-id2" ||
		mock.AttachPaymentMethodParams[0].Customer.ID != "cust-id" {
		t.Fatalf("Attach payment method request not as expected: %v", mock.DetachPaymentMethodIDs)
	}

	if len(mock.UpdateDefaultPaymentMethodParams) != 1 ||
		mock.UpdateDefaultPaymentMethodParams[0].PaymentMethodID != "pm-id2" ||
		mock.UpdateDefaultPaymentMethodParams[0].Customer.ID != "cust-id" {
		t.Fatalf("Update payment method request not as expected: %v", mock.UpdateDefaultPaymentMethodParams)
	}
}
