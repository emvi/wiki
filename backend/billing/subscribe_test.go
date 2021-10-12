package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/config"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"strings"
	"testing"
)

func TestSubscribeActiveSubscription(t *testing.T) {
	testutil.CleanBackendDb(t)
	mailProvider = mailMock
	mock := payment.NewMockClient()
	mock.GetCustomerResult = &stripe.Customer{}
	mock.GetSubscriptionResult = &stripe.Subscription{
		Plan: &stripe.Plan{
			Interval: stripe.PlanIntervalMonth,
		},
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.StripeSubscriptionID.SetValid("sub-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	order := Order{
		Email:           user.Email,
		Name:            "Sherlock Holmes",
		Country:         "gb",
		AddressLine1:    "221b Baker Street",
		PostalCode:      "NW1 6XE",
		City:            "London",
		Phone:           "01+44+207 224 3688",
		TaxNumber:       "GB123456789",
		Interval:        yearly,
		PaymentMethodId: "pm-id",
	}

	if _, err := Subscribe(orga, user.ID, order); len(err) != 1 || err[0] != errs.ActiveSubscriptionFound {
		t.Fatalf("Organization must be expert already, but was: %v", err)
	}
}

func TestSubscribeExistingCustomer(t *testing.T) {
	config.Get().Stripe.YearlyPriceID = "yearly"
	config.Get().Stripe.MonthlyPriceID = "monthly"
	config.Get().Stripe.TaxIDDE = "tax-id"
	testutil.CleanBackendDb(t)
	mailProvider = mailMock
	mock := payment.NewMockClient()
	mock.GetCustomerResult = &stripe.Customer{
		ID: "cust-id",
		InvoiceSettings: &stripe.CustomerInvoiceSettings{
			DefaultPaymentMethod: &stripe.PaymentMethod{
				ID: "pm-id",
			},
		},
	}
	mock.GetSubscriptionResult = &stripe.Subscription{
		Plan: &stripe.Plan{
			Interval: stripe.PlanIntervalMonth,
		},
	}
	mock.CreateSubscriptionResult = &stripe.Subscription{
		ID: "sub-id",
		LatestInvoice: &stripe.Invoice{
			PaymentIntent: &stripe.PaymentIntent{
				Status: stripe.PaymentIntentStatusSucceeded,
			},
		},
	}
	mock.AttachPaymentMethodResult = &stripe.PaymentMethod{
		ID: "pm-id",
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false
	orga.MaxStorageGB = constants.DefaultMaxStorageGb
	orga.StripeCustomerID.SetValid("cust-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	testutil.CreateUser(t, orga, 111, "user1@test.com")
	testutil.CreateUser(t, orga, 222, "user2@test.com")
	user3 := testutil.CreateUser(t, orga, 333, "user3@test.com") // read only
	user4 := testutil.CreateUser(t, orga, 444, "user4@test.com") // inactive
	user3.OrganizationMember.ReadOnly = true
	user4.OrganizationMember.Active = false

	if err := model.SaveOrganizationMember(nil, user3.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, user4.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	order := Order{
		Email:           user.Email,
		Name:            "Sherlock Holmes",
		Country:         "IT",
		AddressLine1:    "221b Baker Street",
		PostalCode:      "NW1 6XE",
		City:            "London",
		Phone:           "01+44+207 224 3688",
		TaxNumber:       "",
		Interval:        yearly,
		PaymentMethodId: "pm-id",
	}

	if _, err := Subscribe(orga, user.ID, order); err != nil {
		t.Fatalf("Subsciption must have been started, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.StripeCustomerID.String != "cust-id" ||
		orga.StripeSubscriptionID.String != "sub-id" ||
		orga.StripePaymentMethodID.String != "pm-id" {
		t.Fatalf("Organization must have been updated, but was: %v %v %v", orga.StripeCustomerID.String, orga.StripeSubscriptionID.String, orga.StripePaymentMethodID.String)
	}

	if !orga.Expert || orga.MaxStorageGB != constants.StorageGBPerUser*3 {
		t.Fatalf("Orga must have been upgraded immediately, but was: %v %v", orga.Expert, orga.MaxStorageGB)
	}

	if len(mock.AttachPaymentMethodParams) != 1 ||
		mock.AttachPaymentMethodParams[0].Customer.ID != "cust-id" ||
		mock.AttachPaymentMethodParams[0].PaymentMethodID != "pm-id" {
		t.Fatalf("Payment method must have been attached, but was: %v", mock.AttachPaymentMethodParams)
	}

	if len(mock.UpdateDefaultPaymentMethodParams) != 1 ||
		mock.UpdateDefaultPaymentMethodParams[0].Customer.ID != "cust-id" ||
		mock.UpdateDefaultPaymentMethodParams[0].PaymentMethodID != "pm-id" {
		t.Fatalf("Default payment method must have been updated, but was: %v", mock.UpdateDefaultPaymentMethodParams)
	}

	if len(mock.CreateSubscriptionParams) != 1 ||
		mock.CreateSubscriptionParams[0].Customer.ID != "cust-id" ||
		mock.CreateSubscriptionParams[0].Quantity != 3 ||
		mock.CreateSubscriptionParams[0].PlanID != "yearly" ||
		mock.CreateSubscriptionParams[0].TaxID != "tax-id" {
		t.Fatalf("Subscription must have been created, but was: %v", mock.CreateSubscriptionParams)
	}

	if len(mock.GetCustomerIDs) != 1 || mock.GetCustomerIDs[0] != "cust-id" {
		t.Fatalf("Existing customer must have been read, but was: %v", mock.GetCustomerIDs)
	}

	if len(mock.UpdateCustomerParams) != 1 ||
		mock.UpdateCustomerParams[0].ID != "cust-id" ||
		mock.UpdateCustomerParams[0].Name != "Sherlock Holmes" ||
		mock.UpdateCustomerParams[0].Country != "IT" ||
		mock.UpdateCustomerParams[0].AddressLine1 != "221b Baker Street" ||
		mock.UpdateCustomerParams[0].PostalCode != "NW1 6XE" ||
		mock.UpdateCustomerParams[0].City != "London" ||
		mock.UpdateCustomerParams[0].Phone != "01+44+207 224 3688" ||
		mock.UpdateCustomerParams[0].TaxNumber != "" ||
		mock.UpdateCustomerParams[0].TaxNumberType != string(stripe.TaxIDTypeEUVAT) ||
		mock.UpdateCustomerParams[0].Email != user.Email {
		t.Fatalf("Updated customer not as expected: %v", mock.UpdateCustomerParams)
	}

	if len(mock.DetachPaymentMethodIDs) != 1 {
		t.Fatalf("Detach payment method not as expected: %v", mock.DetachPaymentMethodIDs)
	}
}

func TestSubscribeNewCustomer(t *testing.T) {
	config.Get().Stripe.YearlyPriceID = "yearly"
	config.Get().Stripe.MonthlyPriceID = "monthly"
	config.Get().Stripe.TaxIDDE = "tax-id"
	testutil.CleanBackendDb(t)
	mailProvider = mailMock
	mock := payment.NewMockClient()
	mock.CreateCustomerResult = &stripe.Customer{ID: "cust-id"}
	mock.GetSubscriptionResult = &stripe.Subscription{
		Plan: &stripe.Plan{
			Interval: stripe.PlanIntervalMonth,
		},
	}
	mock.CreateSubscriptionResult = &stripe.Subscription{
		ID: "sub-id",
		LatestInvoice: &stripe.Invoice{
			PaymentIntent: &stripe.PaymentIntent{
				Status: stripe.PaymentIntentStatusSucceeded,
			},
		},
	}
	mock.AttachPaymentMethodResult = &stripe.PaymentMethod{
		ID: "pm-id",
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	orga.Expert = false
	orga.MaxStorageGB = constants.DefaultMaxStorageGb

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	testutil.CreateUser(t, orga, 111, "user1@test.com")
	testutil.CreateUser(t, orga, 222, "user2@test.com")
	user3 := testutil.CreateUser(t, orga, 333, "user3@test.com") // read only
	user4 := testutil.CreateUser(t, orga, 444, "user4@test.com") // inactive
	user3.OrganizationMember.ReadOnly = true
	user4.OrganizationMember.Active = false

	if err := model.SaveOrganizationMember(nil, user3.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, user4.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	order := Order{
		Email:           user.Email,
		Name:            "Sherlock Holmes",
		Country:         "GB",
		AddressLine1:    "221b Baker Street",
		PostalCode:      "NW1 6XE",
		City:            "London",
		Phone:           "01+44+207 224 3688",
		TaxNumber:       "GB123456789",
		Interval:        monthly,
		PaymentMethodId: "pm-id",
	}

	if _, err := Subscribe(orga, user.ID, order); err != nil {
		t.Fatalf("Subsciption must have been started, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.StripeCustomerID.String != "cust-id" ||
		orga.StripeSubscriptionID.String != "sub-id" {
		t.Fatalf("Organization must have been updated, but was: %v %v", orga.StripeCustomerID.String, orga.StripeSubscriptionID.String)
	}

	if !orga.Expert ||
		orga.MaxStorageGB != constants.StorageGBPerUser*3 ||
		orga.SubscriptionPlan.String != monthly {
		t.Fatalf("Orga must have been upgraded immediately, but was: %v %v %v", orga.Expert, orga.MaxStorageGB, orga.SubscriptionPlan)
	}

	if len(mock.AttachPaymentMethodParams) != 1 ||
		mock.AttachPaymentMethodParams[0].Customer.ID != "cust-id" ||
		mock.AttachPaymentMethodParams[0].PaymentMethodID != "pm-id" {
		t.Fatalf("Payment method must have been attached, but was: %v", mock.AttachPaymentMethodParams)
	}

	if len(mock.UpdateDefaultPaymentMethodParams) != 1 ||
		mock.UpdateDefaultPaymentMethodParams[0].Customer.ID != "cust-id" ||
		mock.UpdateDefaultPaymentMethodParams[0].PaymentMethodID != "pm-id" {
		t.Fatalf("Default payment method must have been updated, but was: %v", mock.UpdateDefaultPaymentMethodParams)
	}

	if len(mock.CreateSubscriptionParams) != 1 ||
		mock.CreateSubscriptionParams[0].Customer.ID != "cust-id" ||
		mock.CreateSubscriptionParams[0].Quantity != 3 ||
		mock.CreateSubscriptionParams[0].PlanID != "monthly" ||
		mock.CreateSubscriptionParams[0].TaxID != "" {
		t.Fatalf("Subscription must have been created, but was: %v", mock.CreateSubscriptionParams)
	}

	if len(mock.CreateCustomerParams) != 1 ||
		mock.CreateCustomerParams[0].Email != user.Email ||
		mock.CreateCustomerParams[0].Name != "Sherlock Holmes" ||
		mock.CreateCustomerParams[0].Country != "GB" ||
		mock.CreateCustomerParams[0].AddressLine1 != "221b Baker Street" ||
		mock.CreateCustomerParams[0].PostalCode != "NW1 6XE" ||
		mock.CreateCustomerParams[0].City != "London" ||
		mock.CreateCustomerParams[0].Phone != "01+44+207 224 3688" ||
		mock.CreateCustomerParams[0].TaxNumber != "GB123456789" ||
		mock.CreateCustomerParams[0].TaxNumberType != string(stripe.TaxIDTypeUnknown) ||
		mock.CreateCustomerParams[0].OrganizationID != orga.ID {
		t.Fatalf("Customer must have been created, but was: %v", mock.CreateSubscriptionParams)
	}

	if len(mock.UpdateCustomerParams) != 0 {
		t.Fatal("Customer must not have been updated")
	}

	if len(mock.DetachPaymentMethodIDs) != 0 {
		t.Fatal("Payment method must not have been detached")
	}
}

func TestValidateSubscriptionOrder(t *testing.T) {
	input := []Order{
		{"", "User Name", "gb", "address", "", "1234", "city", "", "", "yearly", "id"},
		{"user@mail.com", "", "gb", "address", "", "1234", "city", "", "", "yearly", "id"},
		{"user@mail.com", "User Name", "", "address", "", "1234", "city", "", "", "yearly", "id"},
		{"user@mail.com", "User Name", "gb", "", "", "1234", "city", "", "", "yearly", "id"},
		{"user@mail.com", "User Name", "gb", "address", "", "", "city", "", "", "yearly", "id"},
		{"user@mail.com", "User Name", "gb", "address", "", "1234", "", "", "", "yearly", "id"},
		{"user@mail.com", "User Name", "gb", "address", "", "1234", "city", "", "", "invalid", "id"},
	}
	expected := []error{
		errs.EmailInvalid,
		errs.NameTooShort,
		errs.CountryInvalid,
		errs.AddressLineTooShort,
		errs.PostalCodeTooShort,
		errs.CityTooShort,
		errs.BillingIntervalInvalid,
	}

	for i, in := range input {
		if err := in.validate(); err[0] != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], err)
		}
	}

	order := Order{"user@mail.com", "User Name", "gb", "address", "", "1234", "city", "", "", "yearly", "id"}

	if err := order.validate(); err != nil {
		t.Fatalf("Order must be valid, but was: %v", err)
	}

	order = Order{"user@mail.com", "User Name", "gb", "address", "", "1234", "city", "", "", "monthly", "id"}

	if err := order.validate(); err != nil {
		t.Fatalf("Order must be valid, but was: %v", err)
	}
}

func TestSendSubscriptionMail(t *testing.T) {
	testutil.CleanBackendDb(t)
	var subject, body, to string
	mailProvider = func(mailSubject string, mailBody string, mailFrom string, mailTo ...string) error {
		subject = mailSubject
		body = mailBody
		to = mailFrom
		return nil
	}
	orga, user := testutil.CreateOrgaAndUser(t)
	sendSubscriptionMail(orga, user)

	if subject != "Your subscription at Emvi" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if to != user.Email {
		t.Fatalf("Receiver not as expected: %v", to)
	}

	t.Log(body)

	if !strings.Contains(body, orga.Name) ||
		!strings.Contains(body, string(subscriptionMailI18n["en"]["title"])) ||
		!strings.Contains(body, string(subscriptionMailI18n["en"]["text-1"])) ||
		!strings.Contains(body, string(subscriptionMailI18n["en"]["text-2"])) ||
		!strings.Contains(body, string(subscriptionMailI18n["en"]["greeting"])) ||
		!strings.Contains(body, string(subscriptionMailI18n["en"]["goodbye"])) {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func mailMock(subject string, body string, from string, to ...string) error {
	return nil
}
