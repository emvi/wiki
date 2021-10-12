package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"time"
)

const invoiceLimit = 10

type Invoice struct {
	ID      string    `json:"id"`
	Number  string    `json:"number"`
	Total   int64     `json:"total"`
	Paid    bool      `json:"paid"`
	Created time.Time `json:"created"`
	PDFLink string    `json:"pdf_link"`
}

// GetInvoices returns the last view invoices. The next batch of invoices can be accessed by passing the
// previously last invoice ID as a parameter or leave it empty for the first call.
func GetInvoices(orga *model.Organization, userId hide.ID, startInvoiceId string) ([]Invoice, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return nil, err
	}

	if !orga.StripeSubscriptionID.Valid {
		return []Invoice{}, nil
	}

	logbuch.Debug("Start invoice ID", logbuch.Fields{"start_invoice_id": startInvoiceId})
	stripeInvoices, err := client.GetInvoices(orga.StripeSubscriptionID.String, startInvoiceId, invoiceLimit)

	if err != nil {
		logbuch.Error("Error reading invoices", logbuch.Fields{"err": err, "orga_id": orga.ID, "subscription_id": orga.StripeSubscriptionID.String})
		return nil, errs.ReadingInvoices
	}

	invoices := make([]Invoice, 0, len(stripeInvoices))

	for _, invoice := range stripeInvoices {
		invoices = append(invoices, Invoice{
			ID:      invoice.ID,
			Number:  invoice.Number,
			Paid:    invoice.Paid,
			Total:   invoice.Total,
			Created: time.Unix(invoice.Created, 0),
			PDFLink: invoice.InvoicePDF,
		})
	}

	return invoices, nil
}
