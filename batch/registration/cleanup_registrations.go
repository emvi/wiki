package registration

import (
	"emviwiki/auth/model"
	"github.com/emvi/logbuch"
)

func CleanupRegistrations() {
	if err := model.DeleteUserByInactiveAndRegistrationStepAndDefTime(nil); err != nil {
		logbuch.Fatal("Error while cleaning up registrations", logbuch.Fields{"err": err})
	}
}
