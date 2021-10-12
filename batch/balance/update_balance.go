package balance

import (
	"emviwiki/shared/db"
	"emviwiki/shared/model"
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	"sync"
)

const (
	updateBalanceConsumers = 10
	currency               = "USD" // we always use US dollar on invoices and for prices...
)

func UpdateBalance() {
	sendChan := updateBalanceProducer()
	var wg sync.WaitGroup
	wg.Add(updateBalanceConsumers)

	for i := 0; i < updateBalanceConsumers; i++ {
		go updateBalanceConsumer(sendChan, &wg)
	}

	wg.Wait()
}

func updateBalanceProducer() <-chan *model.Organization {
	organizationRows, err := model.FindOrganizationWithSubscriptionCycleReached()

	if err != nil {
		logbuch.Fatal("Error reading organization to update balance", logbuch.Fields{"err": err})
	}

	sendChan := make(chan *model.Organization)

	go func() {
		for organizationRows.Next() {
			var orga model.Organization

			if err := organizationRows.StructScan(&orga); err != nil {
				logbuch.Fatal("Error scanning organization", logbuch.Fields{"err": err})
				continue
			}

			sendChan <- &orga
		}

		db.CloseRows(organizationRows)
		close(sendChan)
	}()

	return sendChan
}

func updateBalanceConsumer(sendChan <-chan *model.Organization, wg *sync.WaitGroup) {
	for orga := range sendChan {
		applyBalance(orga)
	}

	wg.Done()
}

func applyBalance(orga *model.Organization) {
	if !orga.SubscriptionCycle.Valid {
		logbuch.Error("Invalid subscription cycle for organization", logbuch.Fields{"orga_id": orga.ID})
		return
	}

	if !orga.StripeCustomerID.Valid {
		logbuch.Error("Invalid subscription customer for organization", logbuch.Fields{"orga_id": orga.ID})
		return
	}

	price, err := getPrice(orga)

	if err != nil {
		logbuch.Error("Price for subscription not found", logbuch.Fields{"err": err, "orga_id": orga.ID})
		return
	}

	billable := int64(model.CountOrganizationMemberByOrganizationIdAndActiveAndNotReadOnly(orga.ID))
	active := int64(model.CountOrganizationMemberByOrganizationIdAndLastSeenAfter(orga.ID, orga.SubscriptionCycle.Time))
	inactive := billable - active
	amount := inactive * price * -1 // negative because we this is a debit from our perspective

	if amount >= 0 {
		logbuch.Info("Dropping balance modification, as it is <= 0", logbuch.Fields{"orga_id": orga.ID, "amount": amount})
		return
	}

	logbuch.Info("Applying calculated amount to balance", logbuch.Fields{"orga_id": orga.ID, "amount": amount})
	description := fmt.Sprintf("Organization ID: %d", orga.ID)

	if err := client.AddBalance(orga.StripeCustomerID.String, amount, currency, description); err != nil {
		logbuch.Error("Error applying amount to customer balance", logbuch.Fields{"err": err, "orga_id": orga.ID, "amount": amount})
		return
	}

	orga.SubscriptionCycle.SetValid(orga.SubscriptionCycle.Time.AddDate(0, 1, 0))

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization to update subscription cycle", logbuch.Fields{"err": err, "orga_id": orga.ID})
	}
}

func getPrice(orga *model.Organization) (int64, error) {
	if !orga.StripeSubscriptionID.Valid {
		return 0, errors.New("stripe subscription invalid")
	}

	sub, err := client.GetSubscription(orga.StripeSubscriptionID.String)

	if err != nil {
		return 0, errors.New("subscription not found")
	}

	return sub.Items.Data[0].Price.UnitAmount, nil
}
