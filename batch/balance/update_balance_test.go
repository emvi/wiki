package balance

import (
	"emviwiki/shared/model"
	"emviwiki/shared/payment"
	"emviwiki/shared/testutil"
	"fmt"
	"github.com/stripe/stripe-go/v71"
	"testing"
	"time"
)

func TestUpdateBalance(t *testing.T) {
	testutil.CleanBackendDb(t)
	mock := payment.NewMockClient()
	mock.GetSubscriptionResult = &stripe.Subscription{
		ID: "sub-id",
		Items: &stripe.SubscriptionItemList{
			Data: []*stripe.SubscriptionItem{
				{
					Price: &stripe.Price{
						ID:         "price-id",
						UnitAmount: 500, // $5
					},
				},
			},
		},
	}
	client = mock
	expert, _ := testutil.CreateOrgaAndUser(t)
	expert.StripeCustomerID.SetValid("cust-id")
	expert.StripeSubscriptionID.SetValid("sub-id")
	expert.SubscriptionCycle.SetValid(time.Date(2020, 5, 25, 0, 0, 0, 0, time.Local))

	if err := model.SaveOrganization(nil, expert); err != nil {
		t.Fatal(err)
	}

	active1 := testutil.CreateUser(t, expert, 111, "user@test.com")
	active2 := testutil.CreateUser(t, expert, 222, "user@test.com")
	inactive := testutil.CreateUser(t, expert, 333, "user@test.com")
	ro := testutil.CreateUser(t, expert, 444, "ro@test.com")
	disabled := testutil.CreateUser(t, expert, 555, "disabled@test.com")
	active1.OrganizationMember.LastSeen = time.Date(2020, 6, 14, 0, 0, 0, 0, time.Local)
	active2.OrganizationMember.LastSeen = time.Date(2020, 6, 9, 0, 0, 0, 0, time.Local)
	inactive.OrganizationMember.LastSeen = time.Date(2020, 5, 24, 0, 0, 0, 0, time.Local)
	ro.OrganizationMember.ReadOnly = true
	ro.OrganizationMember.LastSeen = time.Date(2020, 5, 26, 0, 0, 0, 0, time.Local)
	disabled.OrganizationMember.Active = false
	disabled.OrganizationMember.LastSeen = time.Date(2020, 5, 26, 0, 0, 0, 0, time.Local)

	if err := model.SaveOrganizationMember(nil, active1.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, active2.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, inactive.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, ro.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	if err := model.SaveOrganizationMember(nil, disabled.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	UpdateBalance()

	if len(mock.AddBalanceParams) != 1 {
		t.Fatal("Balance must have been applied")
	}

	if mock.AddBalanceParams[0].CustomerID != "cust-id" ||
		mock.AddBalanceParams[0].Amount != -1000 ||
		mock.AddBalanceParams[0].Currency != currency ||
		mock.AddBalanceParams[0].Description != fmt.Sprintf("Organization ID: %d", expert.ID) {
		t.Fatalf("API call not as expected: %v", mock.AddBalanceParams[0])
	}

	expert = model.GetOrganizationById(expert.ID)

	if expert.SubscriptionCycle.Time.Year() != 2020 ||
		expert.SubscriptionCycle.Time.Month() != 6 ||
		expert.SubscriptionCycle.Time.Day() != 25 {
		t.Fatalf("Subscription cylce must have been updated, but was: %v", expert.SubscriptionCycle)
	}
}
