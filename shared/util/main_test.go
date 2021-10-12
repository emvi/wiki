package util

import (
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	conn := testutil.ConnectBackend(false)
	defer conn.Disconnect()
	os.Exit(m.Run())
}
