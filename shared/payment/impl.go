package payment

import (
	"emviwiki/shared/config"
	"github.com/emvi/hide"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
	"strconv"
)

type Client struct {
	*client.API
}

func NewStripeClient() BillingClient {
	return &Client{
		API: client.New(config.Get().Stripe.PrivateKey, nil),
	}
}

func (client *Client) CreateCustomer(email, name, country, addressLine1, addressLine2, postalCode, city, phone, taxNumber, taxNumberType, taxExempt string, orgaId hide.ID) (*stripe.Customer, error) {
	var taxIdData []*stripe.CustomerTaxIDDataParams

	if taxNumber != "" {
		taxIdData = []*stripe.CustomerTaxIDDataParams{
			{Type: stripe.String(taxNumberType), Value: stripe.String(taxNumber)},
		}
	}

	return client.Customers.New(&stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Phone: stripe.String(phone),
		Address: &stripe.AddressParams{
			City:       stripe.String(city),
			Country:    stripe.String(country),
			Line1:      stripe.String(addressLine1),
			Line2:      stripe.String(addressLine2),
			PostalCode: stripe.String(postalCode),
		},
		Params: stripe.Params{
			Metadata: map[string]string{"organization_id": strconv.Itoa(int(orgaId))},
		},
		TaxIDData: taxIdData,
		TaxExempt: stripe.String(taxExempt),
	})
}

func (client *Client) UpdateCustomer(id, email, name, country, addressLine1, addressLine2, postalCode, city, phone, taxNumber, taxNumberType, taxExempt string) (*stripe.Customer, error) {
	iter := client.TaxIDs.List(&stripe.TaxIDListParams{
		Customer: stripe.String(id),
	})
	var taxId *stripe.TaxID

	for iter.Next() {
		if err := iter.Err(); err != nil {
			return nil, err
		}

		taxId = iter.TaxID()
		break
	}

	if taxId == nil && taxNumber != "" {
		// attach a new tax ID
		_, err := client.TaxIDs.New(&stripe.TaxIDParams{
			Customer: stripe.String(id),
			Type:     stripe.String(taxNumberType),
			Value:    stripe.String(taxNumber),
		})

		if err != nil {
			return nil, err
		}
	} else if taxId != nil && taxNumber == "" {
		// delete tax ID
		if _, err := client.TaxIDs.Del(taxId.ID, &stripe.TaxIDParams{Customer: stripe.String(id)}); err != nil {
			return nil, err
		}
	} else if taxId != nil && taxNumber != "" && taxId.Value != taxNumber {
		// update tax ID
		if _, err := client.TaxIDs.Del(taxId.ID, &stripe.TaxIDParams{Customer: stripe.String(id)}); err != nil {
			return nil, err
		}

		_, err := client.TaxIDs.New(&stripe.TaxIDParams{
			Customer: stripe.String(id),
			Type:     stripe.String(taxNumberType),
			Value:    stripe.String(taxNumber),
		})

		if err != nil {
			return nil, err
		}
	}

	return client.Customers.Update(id, &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Phone: stripe.String(phone),
		Address: &stripe.AddressParams{
			City:       stripe.String(city),
			Country:    stripe.String(country),
			Line1:      stripe.String(addressLine1),
			Line2:      stripe.String(addressLine2),
			PostalCode: stripe.String(postalCode),
		},
		TaxExempt: stripe.String(taxExempt),
	})
}

func (client *Client) GetCustomer(id string) (*stripe.Customer, error) {
	return client.Customers.Get(id, nil)
}

func (client *Client) DeleteCustomer(id string) error {
	_, err := client.Customers.Del(id, nil)
	return err
}

func (client *Client) CreateSubscription(customer *stripe.Customer, planId string, quantity int64, taxId string) (*stripe.Subscription, error) {
	taxRates := stripe.StringSlice([]string{})

	if taxId != "" {
		taxRates = append(taxRates, stripe.String(taxId))
	}

	subscriptionParams := &stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan:     stripe.String(planId),
				Quantity: stripe.Int64(quantity),
				TaxRates: taxRates,
			},
		},
		PendingInvoiceItemInterval: &stripe.SubscriptionPendingInvoiceItemIntervalParams{
			Interval:      stripe.String(string(stripe.SubscriptionPendingInvoiceItemIntervalIntervalMonth)),
			IntervalCount: stripe.Int64(1),
		},
	}
	subscriptionParams.AddExpand("latest_invoice.payment_intent")
	return client.Subscriptions.New(subscriptionParams)
}

func (client *Client) GetSubscription(id string) (*stripe.Subscription, error) {
	return client.Subscriptions.Get(id, nil)
}

func (client *Client) UpdateSubscriptionQuantity(id, planId string, quantity int64) error {
	_, err := client.Subscriptions.Update(id, &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{ID: stripe.String(planId), Quantity: stripe.Int64(quantity)},
		},
	})
	return err
}

func (client *Client) UpdateSubscriptionPrice(id, planId, priceId string) error {
	_, err := client.Subscriptions.Update(id, &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{ID: stripe.String(planId), Price: stripe.String(priceId)},
		},
	})
	return err
}

func (client *Client) UpdateSubscriptionTaxID(id, planId, taxId string) error {
	// an empty list will remove the tax rate
	taxRates := stripe.StringSlice([]string{})

	if taxId != "" {
		taxRates = append(taxRates, stripe.String(taxId))
	}

	_, err := client.Subscriptions.Update(id, &stripe.SubscriptionParams{
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:       stripe.String(planId),
				TaxRates: taxRates,
			},
		},
	})
	return err
}

func (client *Client) MarkSubscriptionCancelled(id string) error {
	_, err := client.Subscriptions.Update(id, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	})
	return err
}

func (client *Client) ResumeSubscription(id string) error {
	_, err := client.Subscriptions.Update(id, &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	})
	return err
}

func (client *Client) CancelSubscription(id string) error {
	_, err := client.Subscriptions.Cancel(id, nil)
	return err
}

func (client *Client) AttachPaymentMethod(customer *stripe.Customer, paymentMethodId string) (*stripe.PaymentMethod, error) {
	return client.PaymentMethods.Attach(paymentMethodId, &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customer.ID),
	})
}

func (client *Client) UpdateDefaultPaymentMethod(customer *stripe.Customer, paymentMethodId string) (*stripe.Customer, error) {
	return client.Customers.Update(customer.ID, &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethodId),
		},
	})
}

func (client *Client) GetInvoices(subId string, invoiceId string, limit int) ([]*stripe.Invoice, error) {
	params := &stripe.InvoiceListParams{
		Subscription: stripe.String(subId),
	}
	params.Filters.AddFilter("limit", "", strconv.Itoa(limit))

	if invoiceId != "" {
		params.Filters.AddFilter("starting_after", "", invoiceId)
	}

	iterator := client.Invoices.List(params)
	invoices := make([]*stripe.Invoice, 0)

	for iterator.Next() {
		invoices = append(invoices, iterator.Invoice())
	}

	if err := iterator.Err(); err != nil {
		return nil, err
	}

	return invoices, nil
}

func (client *Client) DetachPaymentMethod(id string) error {
	_, err := client.PaymentMethods.Detach(id, nil)
	return err
}

func (client *Client) GetPaymentMethod(id string) (*stripe.PaymentMethod, error) {
	return client.PaymentMethods.Get(id, nil)
}

func (client *Client) GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	return client.PaymentIntents.Get(id, nil)
}

func (client *Client) AddBalance(customerId string, amount int64, currency, description string) error {
	_, err := client.CustomerBalanceTransactions.New(&stripe.CustomerBalanceTransactionParams{
		Customer:    stripe.String(customerId),
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(currency),
		Description: stripe.String(description),
	})
	return err
}
