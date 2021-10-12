package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/stripe/stripe-go/v71"
)

// Customer is the customer as stored on Stripe, but only the part we show in the UI.
type Customer struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Country      string `json:"country"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	PostalCode   string `json:"postal_code"`
	City         string `json:"city"`
	Phone        string `json:"phone"`
	TaxNumber    string `json:"tax_number"`
	Balance      int64  `json:"balance"`
}

// PaymentMethod is the payment method as stored on Stripe, but only the part we show in the UI.
type PaymentMethod struct {
	ID       string `json:"id"`
	Brand    string `json:"brand"`
	ExpMonth uint64 `json:"exp_month"`
	ExpYear  uint64 `json:"exp_year"`
	Last4    string `json:"last_4"`
}

// GetSubscription returns the customer and default payment method details for given organization.
func GetSubscription(orga *model.Organization, userId hide.ID) (*Customer, *PaymentMethod, error) {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return nil, nil, err
	}

	if !orga.StripeCustomerID.Valid {
		return nil, nil, errs.CustomerNotFound
	}

	customer, err := client.GetCustomer(orga.StripeCustomerID.String)

	if err != nil {
		logbuch.Error("Error reading customer", logbuch.Fields{"err": err, "orga_id": orga.ID, "customer_id": orga.StripeCustomerID.String})
		return nil, nil, errs.CustomerNotFound
	}

	paymentMethod, err := client.GetPaymentMethod(customer.InvoiceSettings.DefaultPaymentMethod.ID)

	if err != nil {
		logbuch.Error("Error reading payment method", logbuch.Fields{"err": err, "orga_id": orga.ID, "payment_method_id": customer.InvoiceSettings.DefaultPaymentMethod.ID})
		return nil, nil, errs.PaymentMethodNotFound
	}

	c := &Customer{
		ID:           customer.ID,
		Name:         customer.Name,
		Email:        customer.Email,
		Country:      customer.Address.Country,
		AddressLine1: customer.Address.Line1,
		AddressLine2: customer.Address.Line2,
		PostalCode:   customer.Address.PostalCode,
		City:         customer.Address.City,
		Phone:        customer.Phone,
		TaxNumber:    getCustomerTaxNumber(customer),
		Balance:      customer.Balance * -1, // our debit is the customers credit and the other way around
	}
	pm := &PaymentMethod{
		ID:       paymentMethod.ID,
		Brand:    string(paymentMethod.Card.Brand),
		ExpMonth: paymentMethod.Card.ExpMonth,
		ExpYear:  paymentMethod.Card.ExpYear,
		Last4:    paymentMethod.Card.Last4,
	}
	return c, pm, nil
}

func getCustomerTaxNumber(customer *stripe.Customer) string {
	if customer.TaxIDs != nil && len(customer.TaxIDs.Data) > 0 {
		return customer.TaxIDs.Data[0].Value
	}

	return ""
}
