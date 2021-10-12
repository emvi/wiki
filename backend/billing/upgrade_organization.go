package billing

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
	"time"
)

// UpgradeOrganization upgrades the given organization to expert.
// CAREFUL! This must stay the one thing this function does! It must not have any side effects!
func UpgradeOrganization(orga *model.Organization) error {
	logbuch.Debug("Upgrading organization", logbuch.Fields{"orga_id": orga.ID})
	member := model.CountOrganizationMemberByOrganizationIdAndActiveAndNotReadOnly(orga.ID)
	orga.Expert = true
	orga.MaxStorageGB = int64(constants.StorageGBPerUser * member)

	if !orga.SubscriptionCycle.Valid {
		orga.SubscriptionCycle.SetValid(yesterday())
	}

	if err := model.SaveOrganization(nil, orga); err != nil {
		logbuch.Error("Error saving organization while upgrading", logbuch.Fields{"err": err, "orga_id": orga.ID})
		return errs.Saving
	}

	return nil
}

func yesterday() time.Time {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return today.AddDate(0, 0, -1)
}
