package feed

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/jmoiron/sqlx"
)

// ToggleNotificationRead toggles a notification read/unread or all read depending on the id parameter.
// If set, the noification status will be toggled, else all notifications will be marked as read.
func ToggleNotificationRead(tx *sqlx.Tx, organization *model.Organization, userId, id hide.ID) error {
	if id != 0 {
		// toggle single notification as read/unread
		access := model.GetFeedAccessByOrganizationIdAndUserIdAndFeedIdAndNotificationTx(tx, organization.ID, userId, id, true)

		if access == nil {
			return errs.FeedAccessNotFound
		}

		access.Read = !access.Read

		if err := model.SaveFeedAccess(tx, access); err != nil {
			return errs.Saving
		}
	} else {
		// mark all notifications as read
		if err := model.UpdateFeedAccessNotificationByOrganizationIdAndUserId(tx, organization.ID, userId, true); err != nil {
			return errs.Saving
		}
	}

	return nil
}
