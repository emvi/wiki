package payment

import (
	"github.com/emvi/hide"
	"github.com/stripe/stripe-go/v71"
)

// BillingClient is an interface used for billing using Stripe.
type BillingClient interface {
	// CreateCustomer creates a new customer for given email, name and organization ID.
	CreateCustomer(string, string, string, string, string, string, string, string, string, string, string, hide.ID) (*stripe.Customer, error)

	// UpdateCustomer updates the customers email and name for given customer ID.
	UpdateCustomer(string, string, string, string, string, string, string, string, string, string, string, string) (*stripe.Customer, error)

	// GetCustomer returns the customer for given ID.
	GetCustomer(string) (*stripe.Customer, error)

	// DeleteCustomer deletes the customer for given ID.
	DeleteCustomer(string) error

	// CreateSubscription creates a new subscription for given customer, plan ID quantity and tax ID.
	CreateSubscription(*stripe.Customer, string, int64, string) (*stripe.Subscription, error)

	// GetSubscription returns the subscription for given ID.
	GetSubscription(string) (*stripe.Subscription, error)

	// UpdateSubscriptionQuantity updates the quantity for given subscription and plan ID.
	UpdateSubscriptionQuantity(string, string, int64) error

	// UpdateSubscriptionPrice updates the price ID for given subscription and plan ID.
	UpdateSubscriptionPrice(string, string, string) error

	// UpdateSubscriptionPrice updates the tax ID for given subscription and plan ID.
	UpdateSubscriptionTaxID(string, string, string) error

	// MarkSubscriptionCancelled marks the subscription for given ID as cancelled.
	// This won't directly cancel the subscription, as the user might have time left.
	MarkSubscriptionCancelled(string) error

	// ResumeSubscription resumes the subscription for given ID.
	ResumeSubscription(string) error

	// CancelSubscription cancels the subscription for given ID.
	CancelSubscription(string) error

	// AttachPaymentMethod attaches the given payment method ID to given customer.
	AttachPaymentMethod(*stripe.Customer, string) (*stripe.PaymentMethod, error)

	// UpdateDefaultPaymentMethod updates the default payment method to given payment method ID for given customer.
	UpdateDefaultPaymentMethod(*stripe.Customer, string) (*stripe.Customer, error)

	// GetInvoices returns all invoices for given subscription ID, start and limit.
	// Start must be the last invoice read before.
	GetInvoices(string, string, int) ([]*stripe.Invoice, error)

	// DetachPaymentMethod deletes the payment method for given ID.
	DetachPaymentMethod(string) error

	// GetPaymentMethod returns the payment method for given ID.
	GetPaymentMethod(string) (*stripe.PaymentMethod, error)

	// GetPaymentIntent returns the payment intent for given ID.
	GetPaymentIntent(string) (*stripe.PaymentIntent, error)

	// AddBalance adds given amount to given customer ID balance.
	// Additionally the currency and a description must be passed.
	AddBalance(string, int64, string, string) error
}
