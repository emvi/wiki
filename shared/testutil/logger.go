package testutil

import (
	"github.com/emvi/logbuch"
)

func SetTestLogger() {
	logbuch.SetFormatter(logbuch.NewFieldFormatter("", "\t"))
}
