package newsletter

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
	config.Get().Migrate.Dir = "../../dashboard/schema"
	dashboardconn := testutil.ConnectDashboard(true)
	defer dashboardconn.Disconnect()
	code := m.Run()
	testutil.CheckOpenConnectionsNull(conn)
	testutil.CheckOpenConnectionsNull(dashboardconn)
	os.Exit(code)
}
