package payment

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/stripe/stripe-go/v71"
)

type MockCreateCustomer struct {
	Email          string
	Name           string
	Country        string
	AddressLine1   string
	AddressLine2   string
	PostalCode     string
	City           string
	Phone          string
	TaxNumber      string
	TaxNumberType  string
	TaxExempt      string
	OrganizationID hide.ID
}

type MockUpdateCustomer struct {
	ID            string
	Email         string
	Name          string
	Country       string
	AddressLine1  string
	AddressLine2  string
	PostalCode    string
	City          string
	Phone         string
	TaxNumber     string
	TaxNumberType string
	TaxExempt     string
}

type MockCreateSubscription struct {
	Customer *stripe.Customer
	PlanID   string
	Quantity int64
	TaxID    string
}

type MockUpdateSubscriptionQuantity struct {
	ID       string
	PlanID   string
	Quantity int64
}

type MockUpdateSubscriptionPrice struct {
	ID      string
	PlanID  string
	PriceID string
}

type MockUpdateSubscriptionTaxID struct {
	ID     string
	PlanID string
	TaxID  string
}

type MockNewCheckoutSession struct {
	Organization *model.Organization
	PriceID      string
	Quantity     int64
	Email        string
}

type MockCreateInvoice struct {
	SubscriptionID string
	CustomerID     string
}

type MockAttachPaymentMethod struct {
	Customer        *stripe.Customer
	PaymentMethodID string
}

type MockUpdateDefaultPaymentMethod struct {
	Customer        *stripe.Customer
	PaymentMethodID string
}

type MockGetInvoices struct {
	SubscriptionID string
	InvoiceID      string
	Limit          int
}

type MockAddBalance struct {
	CustomerID  string
	Amount      int64
	Currency    string
	Description string
}

type MockClient struct {
	CreateCustomerParams             []MockCreateCustomer
	CreateCustomerResult             *stripe.Customer
	CreateCustomerError              error
	UpdateCustomerParams             []MockUpdateCustomer
	UpdateCustomerResult             *stripe.Customer
	UpdateCustomerError              error
	GetCustomerIDs                   []string
	GetCustomerResult                *stripe.Customer
	GetCustomerError                 error
	DeleteCustomerIDs                []string
	DeleteCustomerError              error
	CreateSubscriptionParams         []MockCreateSubscription
	CreateSubscriptionResult         *stripe.Subscription
	CreateSubscriptionError          error
	GetSubscriptionIDs               []string
	GetSubscriptionResult            *stripe.Subscription
	GetSubscriptionError             error
	UpdateSubscriptionQuantityParams []MockUpdateSubscriptionQuantity
	UpdateSubscriptionQuantityError  error
	UpdateSubscriptionPriceParams    []MockUpdateSubscriptionPrice
	UpdateSubscriptionPriceError     error
	UpdateSubscriptionTaxIDParams    []MockUpdateSubscriptionTaxID
	UpdateSubscriptionTaxIDError     error
	MarkSubscriptionCancelledIDs     []string
	MarkSubscriptionCancelledError   error
	ResumeSubscriptionIDs            []string
	ResumeSubscriptionError          error
	CancelSubscriptionIDs            []string
	CancelSubscriptionError          error
	CreateInvoiceParams              []MockCreateInvoice
	CreateInvoiceResult              *stripe.Invoice
	CreateInvoiceError               error
	AttachPaymentMethodParams        []MockAttachPaymentMethod
	AttachPaymentMethodResult        *stripe.PaymentMethod
	AttachPaymentMethodError         error
	UpdateDefaultPaymentMethodParams []MockUpdateDefaultPaymentMethod
	UpdateDefaultPaymentMethodResult *stripe.Customer
	UpdateDefaultPaymentMethodError  error
	GetInvoicesParams                []MockGetInvoices
	GetInvoicesResult                []*stripe.Invoice
	GetInvoicesError                 error
	DetachPaymentMethodIDs           []string
	DetachPaymentMethodError         error
	GetPaymentMethodIDs              []string
	GetPaymentMethodResult           *stripe.PaymentMethod
	GetPaymentMethodError            error
	GetPaymentIntentIDs              []string
	GetPaymentIntentResult           *stripe.PaymentIntent
	GetPaymentIntentError            error
	AddBalanceParams                 []MockAddBalance
	AddBalanceError                  error
}

func NewMockClient() *MockClient {
	return &MockClient{
		CreateCustomerParams:             make([]MockCreateCustomer, 0),
		UpdateCustomerParams:             make([]MockUpdateCustomer, 0),
		GetCustomerIDs:                   make([]string, 0),
		DeleteCustomerIDs:                make([]string, 0),
		CreateSubscriptionParams:         make([]MockCreateSubscription, 0),
		GetSubscriptionIDs:               make([]string, 0),
		UpdateSubscriptionQuantityParams: make([]MockUpdateSubscriptionQuantity, 0),
		UpdateSubscriptionPriceParams:    make([]MockUpdateSubscriptionPrice, 0),
		UpdateSubscriptionTaxIDParams:    make([]MockUpdateSubscriptionTaxID, 0),
		MarkSubscriptionCancelledIDs:     make([]string, 0),
		ResumeSubscriptionIDs:            make([]string, 0),
		CancelSubscriptionIDs:            make([]string, 0),
		CreateInvoiceParams:              make([]MockCreateInvoice, 0),
		AttachPaymentMethodParams:        make([]MockAttachPaymentMethod, 0),
		UpdateDefaultPaymentMethodParams: make([]MockUpdateDefaultPaymentMethod, 0),
		GetInvoicesParams:                make([]MockGetInvoices, 0),
		DetachPaymentMethodIDs:           make([]string, 0),
		GetPaymentMethodIDs:              make([]string, 0),
		GetPaymentIntentIDs:              make([]string, 0),
		AddBalanceParams:                 make([]MockAddBalance, 0),
	}
}

func (client *MockClient) CreateCustomer(email, name, country, addressLine1, addressLine2, postalCode, city, phone, taxNumber, taxNumberType, taxExempt string, orgaId hide.ID) (*stripe.Customer, error) {
	client.CreateCustomerParams = append(client.CreateCustomerParams, MockCreateCustomer{
		email,
		name,
		country,
		addressLine1,
		addressLine2,
		postalCode,
		city,
		phone,
		taxNumber,
		taxNumberType,
		taxExempt,
		orgaId,
	})
	return client.CreateCustomerResult, client.CreateCustomerError
}

func (client *MockClient) UpdateCustomer(id, email, name, country, addressLine1, addressLine2, postalCode, city, phone, taxNumber, taxNumberType, taxExempt string) (*stripe.Customer, error) {
	client.UpdateCustomerParams = append(client.UpdateCustomerParams, MockUpdateCustomer{
		id,
		email,
		name,
		country,
		addressLine1,
		addressLine2,
		postalCode,
		city,
		phone,
		taxNumber,
		taxNumberType,
		taxExempt,
	})
	return client.UpdateCustomerResult, client.UpdateCustomerError
}

func (client *MockClient) GetCustomer(id string) (*stripe.Customer, error) {
	client.GetCustomerIDs = append(client.GetCustomerIDs, id)
	return client.GetCustomerResult, client.GetCustomerError
}

func (client *MockClient) DeleteCustomer(id string) error {
	client.DeleteCustomerIDs = append(client.DeleteCustomerIDs, id)
	return client.DeleteCustomerError
}

func (client *MockClient) CreateSubscription(customer *stripe.Customer, planId string, quantity int64, taxId string) (*stripe.Subscription, error) {
	client.CreateSubscriptionParams = append(client.CreateSubscriptionParams, MockCreateSubscription{
		customer,
		planId,
		quantity,
		taxId,
	})
	return client.CreateSubscriptionResult, client.CreateSubscriptionError
}

func (client *MockClient) GetSubscription(id string) (*stripe.Subscription, error) {
	client.GetSubscriptionIDs = append(client.GetSubscriptionIDs, id)
	return client.GetSubscriptionResult, client.GetSubscriptionError
}

func (client *MockClient) UpdateSubscriptionQuantity(id, planId string, quantity int64) error {
	client.UpdateSubscriptionQuantityParams = append(client.UpdateSubscriptionQuantityParams, MockUpdateSubscriptionQuantity{
		id,
		planId,
		quantity,
	})
	return client.UpdateSubscriptionQuantityError
}

func (client *MockClient) UpdateSubscriptionPrice(id, planId, priceId string) error {
	client.UpdateSubscriptionPriceParams = append(client.UpdateSubscriptionPriceParams, MockUpdateSubscriptionPrice{
		id,
		planId,
		priceId,
	})
	return client.UpdateSubscriptionPriceError
}

func (client *MockClient) UpdateSubscriptionTaxID(id, planId, taxId string) error {
	client.UpdateSubscriptionTaxIDParams = append(client.UpdateSubscriptionTaxIDParams, MockUpdateSubscriptionTaxID{
		id,
		planId,
		taxId,
	})
	return client.UpdateSubscriptionTaxIDError
}

func (client *MockClient) MarkSubscriptionCancelled(id string) error {
	client.MarkSubscriptionCancelledIDs = append(client.MarkSubscriptionCancelledIDs, id)
	return client.MarkSubscriptionCancelledError
}

func (client *MockClient) ResumeSubscription(id string) error {
	client.ResumeSubscriptionIDs = append(client.ResumeSubscriptionIDs, id)
	return client.ResumeSubscriptionError
}

func (client *MockClient) CancelSubscription(id string) error {
	client.CancelSubscriptionIDs = append(client.CancelSubscriptionIDs, id)
	return client.CancelSubscriptionError
}

func (client *MockClient) AttachPaymentMethod(customer *stripe.Customer, paymentMethodId string) (*stripe.PaymentMethod, error) {
	client.AttachPaymentMethodParams = append(client.AttachPaymentMethodParams, MockAttachPaymentMethod{
		customer,
		paymentMethodId,
	})
	return client.AttachPaymentMethodResult, client.AttachPaymentMethodError
}

func (client *MockClient) UpdateDefaultPaymentMethod(customer *stripe.Customer, paymentMethodId string) (*stripe.Customer, error) {
	client.UpdateDefaultPaymentMethodParams = append(client.UpdateDefaultPaymentMethodParams, MockUpdateDefaultPaymentMethod{
		customer,
		paymentMethodId,
	})
	return client.UpdateDefaultPaymentMethodResult, client.UpdateDefaultPaymentMethodError
}

func (client *MockClient) GetInvoices(subId string, invoiceId string, limit int) ([]*stripe.Invoice, error) {
	client.GetInvoicesParams = append(client.GetInvoicesParams, MockGetInvoices{
		subId,
		invoiceId,
		limit,
	})
	return client.GetInvoicesResult, client.GetInvoicesError
}

func (client *MockClient) DetachPaymentMethod(id string) error {
	client.DetachPaymentMethodIDs = append(client.DetachPaymentMethodIDs, id)
	return client.DetachPaymentMethodError
}

func (client *MockClient) GetPaymentMethod(id string) (*stripe.PaymentMethod, error) {
	client.GetPaymentMethodIDs = append(client.GetPaymentMethodIDs, id)
	return client.GetPaymentMethodResult, client.GetPaymentMethodError
}

func (client *MockClient) GetPaymentIntent(id string) (*stripe.PaymentIntent, error) {
	client.GetPaymentIntentIDs = append(client.GetPaymentIntentIDs, id)
	return client.GetPaymentIntentResult, client.GetPaymentIntentError
}

func (client *MockClient) AddBalance(customerId string, amount int64, currency, description string) error {
	client.AddBalanceParams = append(client.AddBalanceParams, MockAddBalance{
		customerId,
		amount,
		currency,
		description,
	})
	return client.AddBalanceError
}
