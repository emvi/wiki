package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestChangePlan(t *testing.T) {
	testutil.CleanBackendDb(t)
	config.Get().Stripe.MonthlyPriceID = "monthly-plan"
	config.Get().Stripe.YearlyPriceID = "yearly-plan"
	mock := payment.NewMockClient()
	mock.GetSubscriptionResult = &stripe.Subscription{
		ID: "sub-id",
		Items: &stripe.SubscriptionItemList{
			Data: []*stripe.SubscriptionItem{
				{ID: "item-id"},
			},
		},
	}
	client = mock
	orga, user := testutil.CreateOrgaAndUser(t)
	nonAdmin := testutil.CreateUser(t, orga, 321, "non-admin@user.com")

	if err := ChangePlan(orga, nonAdmin.ID, ""); err != errs.BillingIntervalInvalid {
		t.Fatalf("Billing interval must be invalid, but was: %v", err)
	}

	if err := ChangePlan(orga, nonAdmin.ID, monthly); err != errs.PermissionDenied {
		t.Fatalf("Access must be denied, but was: %v", err)
	}

	if err := ChangePlan(orga, user.ID, monthly); err != errs.SubscriptionNotFound {
		t.Fatalf("Subscription must not be found, but was: %v", err)
	}

	orga.StripeSubscriptionID.SetValid("sub-id")
	orga.SubscriptionPlan.SetValid(monthly)

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	if err := ChangePlan(orga, user.ID, monthly); err != nil {
		t.Fatalf("Plan must have been changed, but was: %v", err)
	}

	if err := ChangePlan(orga, user.ID, yearly); err != nil {
		t.Fatalf("Plan must have been changed, but was: %v", err)
	}

	orga = model.GetOrganizationById(orga.ID)

	if orga.SubscriptionPlan.String != yearly {
		t.Fatalf("Plan must have been changed for organization, but was: %v", orga.SubscriptionPlan)
	}

	if len(mock.UpdateSubscriptionPriceParams) != 1 ||
		mock.UpdateSubscriptionPriceParams[0].ID != "sub-id" ||
		mock.UpdateSubscriptionPriceParams[0].PlanID != "item-id" ||
		mock.UpdateSubscriptionPriceParams[0].PriceID != "yearly-plan" {
		t.Fatalf("Update subscription request not as expected: %v", mock.UpdateSubscriptionQuantityParams)
	}
}
