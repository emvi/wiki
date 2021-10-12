package notification

import (
	"emviwiki/shared/config"
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	config.Load()
	config.Get().Template.MailTemplateDir = "../../template/mail/*"
	LoadConfig()
	conn := testutil.ConnectBackend(false)
	defer conn.Disconnect()
	code := m.Run()
	testutil.CheckOpenConnectionsNull(conn)
	os.Exit(code)
}
