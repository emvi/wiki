package usergroup

import (
	"emviwiki/shared/config"
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	config.Load()
	conn := testutil.ConnectBackend(true)
	defer conn.Disconnect()
	os.Exit(m.Run())
}
