package invitation

import (
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
)

func CleanupInvitations() {
	if err := model.DeleteInvitationByDefTimeBeforeOneMonth(nil); err != nil {
		logbuch.Fatal("Error while cleaning up invitations", logbuch.Fields{"err": err})
	}
}
