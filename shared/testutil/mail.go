package testutil

import (
	"github.com/emvi/logbuch"
)

func TestMailSender(subject, msgHTML, from string, to ...string) error {
	logbuch.Debug("Sending test mail", logbuch.Fields{"subject": subject, "msg": msgHTML, "from": from, "to": to})
	return nil
}
