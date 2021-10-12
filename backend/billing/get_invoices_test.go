package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestGetInvoices(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetInvoicesResult = []*stripe.Invoice{
		{
			Number:     "123",
			Paid:       false,
			Total:      5686,
			Created:    54321,
			InvoicePDF: "",
		},
		{
			Number:     "122",
			Paid:       true,
			Total:      4567,
			Created:    44321,
			InvoicePDF: "pdf-link",
		},
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "non-admin@user.com")

	if _, err := GetInvoices(orga, nonAdmin.ID, ""); err != errs.PermissionDenied {
		t.Fatalf("Permission must have been denied, but was: %v", err)
	}

	invoices, err := GetInvoices(orga, user.ID, "start-id")

	if len(invoices) != 0 || err != nil {
		t.Fatalf("No invoices must have been returned, but was: %v %v", len(invoices), err)
	}

	orga.StripeSubscriptionID.SetValid("sub-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	invoices, err = GetInvoices(orga, user.ID, "start-id")

	if len(invoices) != 2 || err != nil {
		t.Fatalf("Invoices must have been returned, but was: %v %v", len(invoices), err)
	}

	if len(mock.GetInvoicesParams) != 1 ||
		mock.GetInvoicesParams[0].SubscriptionID != "sub-id" ||
		mock.GetInvoicesParams[0].Limit != invoiceLimit ||
		mock.GetInvoicesParams[0].InvoiceID != "start-id" {
		t.Fatalf("Invoice call not as expected: %v", mock.GetInvoicesParams)
	}
}
