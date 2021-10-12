package billing

import (
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"github.com/stripe/stripe-go/v71"
	"testing"
)

func TestUpdateSubscriptionNonExpert(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	client = mock
	orga, _ := testutil.CreateOrgaAndUser(t)

	if err := UpdateSubscription(orga); err != nil {
		t.Fatalf("Subscription must not have been updated, but was: %v", err)
	}

	if len(mock.GetSubscriptionIDs) != 0 || len(mock.UpdateSubscriptionQuantityParams) != 0 {
		t.Fatalf("Subscription must not have been updated, but was: %v %v", mock.GetSubscriptionIDs, mock.UpdateSubscriptionQuantityParams)
	}
}

func TestUpdateSubscription(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetSubscriptionResult = &stripe.Subscription{
		ID: "sub-id",
		Items: &stripe.SubscriptionItemList{
			Data: []*stripe.SubscriptionItem{
				{ID: "plan-id", Quantity: 1},
			},
		},
	}
	client = mock
	orga, _ := testutil.CreateOrgaAndUser(t)
	orga.StripeCustomerID.SetValid("cust-id")
	orga.StripeSubscriptionID.SetValid("sub-id")

	if err := model.SaveOrganization(nil, orga); err != nil {
		t.Fatal(err)
	}

	testutil.CreateUser(t, orga, 321, "second@user.com")

	if err := UpdateSubscription(orga); err != nil {
		t.Fatalf("Subscription must have been updated, but was: %v", err)
	}

	if len(mock.GetSubscriptionIDs) != 1 ||
		len(mock.UpdateSubscriptionQuantityParams) != 1 {
		t.Fatalf("Subscription must have been updated, but was: %v %v", mock.GetSubscriptionIDs, mock.UpdateSubscriptionQuantityParams)
	}

	if mock.UpdateSubscriptionQuantityParams[0].ID != "sub-id" ||
		mock.UpdateSubscriptionQuantityParams[0].PlanID != "plan-id" ||
		mock.UpdateSubscriptionQuantityParams[0].Quantity != 2 {
		t.Fatalf("Update not as expected: %v", mock.UpdateSubscriptionQuantityParams[0])
	}
}
