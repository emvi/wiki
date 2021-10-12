package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
	"github.com/stripe/stripe-go/v71"
)

func WebhookEvent(event stripe.Event) error {
	defer func() {
		if err := recover(); err != nil {
			logbuch.Error("Error processing event", logbuch.Fields{"event_id": event.ID, "type": event.Type})
		}
	}()

	switch event.Type {
	case "invoice.paid":
		return invoicePaid(event)
	case "invoice.payment_action_required":
		return invoicePaymentRequiresAction(event)
	case "customer.subscription.deleted":
		return subscriptionDeleted(event)
	}

	return nil
}

func invoicePaid(event stripe.Event) error {
	logbuch.Info("Invoice paid", logbuch.Fields{"event_id": event.ID, "type": event.Type})
	subId := event.GetObjectValue("subscription")
	orga := model.GetOrganizationByStripeSubscriptionID(subId)

	if orga == nil {
		logbuch.Error("Organization for paid invoice not found", logbuch.Fields{"event_id": event.ID, "type": event.Type, "subscription_id": subId})
		return errs.OrganizationNotFound
	}

	return UpgradeOrganization(orga)
}

func invoicePaymentRequiresAction(event stripe.Event) error {
	logbuch.Info("Invoice payment action required", logbuch.Fields{"event_id": event.ID, "type": event.Type})
	customerId := event.GetObjectValue("customer")
	pmId := event.GetObjectValue("payment_intent")
	orga := model.GetOrganizationByStripeCustomerID(customerId)

	if orga == nil {
		logbuch.Error("Organization for invoice payment action required not found", logbuch.Fields{"event_id": event.ID, "type": event.Type, "customer_id": customerId})
		return errs.OrganizationNotFound
	}

	return PaymentActionRequired(orga, pmId)
}

func subscriptionDeleted(event stripe.Event) error {
	logbuch.Info("Subscription deleted", logbuch.Fields{"event_id": event.ID, "type": event.Type})
	subId := event.GetObjectValue("id")
	err := Downgrade(subId)

	// subscription not founds means the organization has been deleted
	if err == errs.SubscriptionNotFound {
		return nil
	}

	return err
}
