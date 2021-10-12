package articlelist

import (
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	conn := testutil.ConnectBackend(true)
	defer conn.Disconnect()
	code := m.Run()
	testutil.CheckOpenConnectionsNull(conn)
	os.Exit(code)
}
