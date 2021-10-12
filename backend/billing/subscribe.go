package billing

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/perm"
	"emviwiki/shared/config"
	"emviwiki/shared/constants"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/stripe/stripe-go/v71"
	"html/template"
	"strings"
)

const (
	monthly             = "monthly"
	yearly              = "yearly"
	subscriptionSubject = "subscription"
)

var subscriptionMailI18n = i18n.Translation{
	"en": {
		"title":    "Your subscription at Emvi",
		"text-1":   "Your subscription for your organization",
		"text-2":   "has been placed. Your organization will be upgraded as soon as we have processed the first payment. You will be notified then and receive the bill in a separate mail.",
		"greeting": "Thank you for subscribing to Emvi!",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Dein Abonnement bei Emvi",
		"text-1":   "Das Abonnement für deine Organisation",
		"text-2":   "wurde platziert. Sobald wir die erste Zahlung verarbeitet haben, wird deine Organisation heraufgestuft. Du wirst dann darüber benachrichtigt und erhälst die Rechnung in einer separaten Mail.",
		"greeting": "Vielen Dank für dein Abonnement bei Emvi!",
		"goodbye":  "Dein Emvi Team",
	},
}

// Order is an order for an Expert organization.
type Order struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Country         string `json:"country"`
	AddressLine1    string `json:"address_line_1"`
	AddressLine2    string `json:"address_line_2"`
	PostalCode      string `json:"postal_code"`
	City            string `json:"city"`
	Phone           string `json:"phone"`
	TaxNumber       string `json:"tax_number"`
	Interval        string `json:"interval"`
	PaymentMethodId string `json:"payment_method_id"`
}

func (order *Order) validate() []error {
	order.Email = strings.TrimSpace(order.Email)
	order.Name = strings.TrimSpace(order.Name)
	order.Country = strings.ToUpper(strings.TrimSpace(order.Country))
	order.AddressLine1 = strings.TrimSpace(order.AddressLine1)
	order.AddressLine2 = strings.TrimSpace(order.AddressLine2)
	order.PostalCode = strings.TrimSpace(order.PostalCode)
	order.City = strings.TrimSpace(order.City)
	order.Phone = strings.TrimSpace(order.Phone)
	order.TaxNumber = strings.TrimSpace(order.TaxNumber)
	order.Interval = strings.TrimSpace(strings.ToLower(order.Interval))
	e := make([]error, 0)

	if !mail.EmailValid(order.Email) {
		e = append(e, errs.EmailInvalid)
	}

	if order.Name == "" {
		e = append(e, errs.NameTooShort)
	}

	if order.Country == "" {
		e = append(e, errs.CountryInvalid)
	}

	if order.AddressLine1 == "" {
		e = append(e, errs.AddressLineTooShort)
	}

	if order.PostalCode == "" {
		e = append(e, errs.PostalCodeTooShort)
	}

	if order.City == "" {
		e = append(e, errs.CityTooShort)
	}

	if order.Interval != monthly && order.Interval != yearly {
		e = append(e, errs.BillingIntervalInvalid)
	}

	if len(e) == 0 {
		return nil
	}

	return e
}

// Subscribe upgrades the organization and stores the customer and subscription ID.
// The subscription will only be created if the organization hasn't been upgraded and no subscription exists yet.
func Subscribe(orga *model.Organization, userId hide.ID, order Order) (string, []error) {
	if err := order.validate(); err != nil {
		return "", err
	}

	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return "", []error{err}
	}

	if orga.Expert && orga.StripeSubscriptionID.Valid {
		logbuch.Error("Running subscription found while subscribing", logbuch.Fields{"orga_id": orga.ID})
		return "", []error{errs.ActiveSubscriptionFound}
	}

	if !orga.Expert && orga.StripeSubscriptionID.Valid {
		if err := ResetSubscription(orga); err != nil {
			return "", []error{err}
		}
	}

	paid, clientSecret, err := createSubscription(orga, order)

	if err != nil {
		return "", []error{err}
	}

	if paid {
		if err := UpgradeOrganization(orga); err != nil {
			return "", []error{err}
		}
	}

	go func() {
		admins, err := getAdmins(orga)

		if err != nil {
			logbuch.Error("No administrators found while subscribing", logbuch.Fields{"err": err, "orga_id": orga.ID})
		}

		for _, admin := range admins {
			sendSubscriptionMail(orga, &admin)
		}
	}()

	return clientSecret, nil
}

func createSubscription(orga *model.Organization, order Order) (bool, string, error) {
	member := model.CountOrganizationMemberByOrganizationIdAndActiveAndNotReadOnly(orga.ID)

	if member == 0 {
		logbuch.Error("Error selecting active non-read-only members", logbuch.Fields{"orga_id": orga.ID})
		return false, "", errs.ProcessingPayment
	}

	customer, err := createUpdateCustomer(orga, order)

	if err != nil {
		return false, "", err
	}

	pm, err := client.AttachPaymentMethod(customer, order.PaymentMethodId)

	if err != nil {
		logbuch.Error("Error attaching payment method", logbuch.Fields{"orga_id": orga.ID})
		return false, "", errs.ProcessingPayment
	}

	if _, err := client.UpdateDefaultPaymentMethod(customer, order.PaymentMethodId); err != nil {
		logbuch.Error("Error updating default payment method", logbuch.Fields{"orga_id": orga.ID})
		return false, "", errs.ProcessingPayment
	}

	taxId := getTaxId(order.Country, order.TaxNumber)
	logbuch.Debug("Using tax ID for subscription", logbuch.Fields{"orga_id": orga.ID, "customer_id": customer.ID, "country": order.Country, "tax_number": order.TaxNumber, "tax_id": taxId})
	sub, err := client.CreateSubscription(customer, getStripePriceID(order.Interval), int64(member), taxId)

	if err != nil {
		logbuch.Error("Error creating subscription", logbuch.Fields{"orga_id": orga.ID})
		return false, "", errs.ProcessingPayment
	}

	orga.StripeCustomerID.SetValid(customer.ID)
	orga.StripeSubscriptionID.SetValid(sub.ID)
	orga.StripePaymentMethodID.SetValid(pm.ID)
	orga.SubscriptionPlan.SetValid(order.Interval)

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while creating subscription", logbuch.Fields{"orga_id": orga.ID, "customer_id": customer.ID, "subscription_id": sub.ID})
		return false, "", errs.ProcessingPayment
	}

	paid := sub.LatestInvoice != nil &&
		sub.LatestInvoice.PaymentIntent != nil &&
		sub.LatestInvoice.PaymentIntent.Status == stripe.PaymentIntentStatusSucceeded
	clientSecret := ""

	if !paid && sub.LatestInvoice != nil &&
		sub.LatestInvoice.PaymentIntent != nil &&
		sub.LatestInvoice.PaymentIntent.Status == stripe.PaymentIntentStatusRequiresAction {
		clientSecret = sub.LatestInvoice.PaymentIntent.ClientSecret
	}

	logbuch.Info("New subscription", logbuch.Fields{
		"orga_id":                              orga.ID,
		"customer_id":                          customer.ID,
		"subscription_id":                      sub.ID,
		"paid":                                 paid,
		"payment_authentication_client_secret": clientSecret,
	})
	return paid, clientSecret, nil
}

func createUpdateCustomer(orga *model.Organization, order Order) (*stripe.Customer, error) {
	var customer *stripe.Customer

	if !orga.StripeCustomerID.Valid {
		var err error
		customer, err = client.CreateCustomer(order.Email,
			order.Name,
			order.Country,
			order.AddressLine1,
			order.AddressLine2,
			order.PostalCode,
			order.City,
			order.Phone,
			order.TaxNumber,
			getTaxIdType(order.Country),
			getTaxExempt(order.Country),
			orga.ID)

		if err != nil {
			logbuch.Error("Error creating new customer", logbuch.Fields{"err": err, "orga_id": orga.ID})
			return nil, errs.ProcessingPayment
		}
	} else {
		var err error
		customer, err = client.GetCustomer(orga.StripeCustomerID.String)

		if err != nil {
			logbuch.Error("Error getting existing customer", logbuch.Fields{"err": err, "orga_id": orga.ID, "customer_id": orga.StripeCustomerID.String})
			return nil, errs.ProcessingPayment
		}

		_, err = client.UpdateCustomer(customer.ID,
			order.Email,
			order.Name,
			order.Country,
			order.AddressLine1,
			order.AddressLine2,
			order.PostalCode,
			order.City,
			order.Phone,
			order.TaxNumber,
			getTaxIdType(order.Country),
			getTaxExempt(order.Country))

		if err != nil {
			logbuch.Error("Error updating existing customer", logbuch.Fields{"err": err, "orga_id": orga.ID, "customer_id": customer.ID})
			return nil, errs.ProcessingPayment
		}

		if err := client.DetachPaymentMethod(customer.InvoiceSettings.DefaultPaymentMethod.ID); err != nil {
			logbuch.Error("Error detaching payment method from customer", logbuch.Fields{"err": err, "orga_id": orga.ID, "customer_id": customer.ID})
			return nil, errs.ProcessingPayment
		}
	}

	return customer, nil
}

func getStripePriceID(interval string) string {
	if interval == yearly {
		return config.Get().Stripe.YearlyPriceID
	}

	return config.Get().Stripe.MonthlyPriceID
}

func getAdmins(orga *model.Organization) ([]model.User, error) {
	adminGroup := model.GetUserGroupByOrganizationIdAndNameTx(nil, orga.ID, constants.GroupAdminName)

	if adminGroup == nil {
		logbuch.Error("Admin group not found", logbuch.Fields{"orga_id": orga.ID})
		return nil, errs.GroupNotFound
	}

	return model.FindUserGroupMemberUserByOrganizationIdAndUserGroupId(orga.ID, adminGroup.ID), nil
}

func sendSubscriptionMail(orga *model.Organization, user *model.User) {
	tpl := mailtpl.Cache.Get()
	var buffer bytes.Buffer
	langCode := util.DetermineSystemSupportedLangCode(orga.ID, user.ID)
	data := struct {
		Organization *model.Organization
		EndVars      map[string]template.HTML
		Vars         map[string]template.HTML
	}{
		orga,
		i18n.GetMailEndI18n(langCode),
		i18n.GetVars(langCode, subscriptionMailI18n),
	}

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.SubscriptionMailTemplate, data); err != nil {
		logbuch.Error("Error executing subscription mail template", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}

	subject := i18n.GetMailTitle(langCode)[subscriptionSubject]

	if err := mailProvider(subject, buffer.String(), user.Email); err != nil {
		logbuch.Error("Error sending subscription mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": user.ID})
	}
}
