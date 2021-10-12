package api

import (
	"emviwiki/backend/billing"
	"emviwiki/backend/context"
	"emviwiki/shared/config"
	"emviwiki/shared/rest"
	"github.com/emvi/logbuch"
	"github.com/stripe/stripe-go/v71/webhook"
	"io/ioutil"
	"net/http"
)

func CreateSubscriptionHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	var order billing.Order

	if err := rest.DecodeJSON(r, &order); err != nil {
		return []error{err}
	}

	clientSecret, err := billing.Subscribe(ctx.Organization, ctx.UserId, order)

	if err != nil {
		return err
	}

	rest.WriteResponse(w, struct {
		ClientSecret string `json:"client_secret"`
	}{
		clientSecret,
	})
	logbuch.Debug("Subscribe request done", logbuch.Fields{"client_secret": clientSecret})
	return nil
}

func CancelSubscriptionHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	if err := billing.CancelSubscription(ctx.Organization, ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func ResumeSubscriptionHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	if err := billing.ResumeSubscription(ctx.Organization, ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func GetInvoicesHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	startInvoiceId := rest.GetParam(r, "start_invoice_id")
	invoices, err := billing.GetInvoices(ctx.Organization, ctx.UserId, startInvoiceId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, invoices)
	return nil
}

func GetSubscriptionHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	customer, pm, err := billing.GetSubscription(ctx.Organization, ctx.UserId)

	if err != nil {
		return []error{err}
	}

	rest.WriteResponse(w, struct {
		Customer      *billing.Customer      `json:"customer"`
		PaymentMethod *billing.PaymentMethod `json:"payment_method"`
	}{
		customer,
		pm,
	})
	return nil
}

func UpdateCustomerHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	var order billing.Order

	if err := rest.DecodeJSON(r, &order); err != nil {
		return []error{err}
	}

	if err := billing.UpdateCustomer(ctx.Organization, ctx.UserId, order); err != nil {
		return err
	}

	return nil
}

func UpdatePaymentMethodHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		PaymentMethodID string `json:"payment_method_id"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := billing.UpdatePaymentMethod(ctx.Organization, ctx.UserId, req.PaymentMethodID); err != nil {
		return []error{err}
	}

	return nil
}

func UpdatePlanHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	req := struct {
		Interval string `json:"interval"`
	}{}

	if err := rest.DecodeJSON(r, &req); err != nil {
		return []error{err}
	}

	if err := billing.ChangePlan(ctx.Organization, ctx.UserId, req.Interval); err != nil {
		return []error{err}
	}

	return nil
}

func RemovePaymentIntentClientSecretHandler(ctx context.EmviContext, w http.ResponseWriter, r *http.Request) []error {
	if err := billing.RemovePaymentIntentClientSecret(ctx.Organization, ctx.UserId); err != nil {
		return []error{err}
	}

	return nil
}

func StripeWebhookHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logbuch.Error("Error reading body from stripe webhook", logbuch.Fields{"err": err})
		rest.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), config.Get().Stripe.WebhookKey)

	if err != nil {
		logbuch.Error("Error constructing webhook event", logbuch.Fields{"err": err})
		rest.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if err := billing.WebhookEvent(event); err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err)
	}
}
