package article

import (
	"emviwiki/backend/content"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/config"
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	os.RemoveAll("bucket")
	os.Setenv("EMVI_WIKI_TEMPLATE_DIR", "../../template/backend/*")
	os.Setenv("EMVI_WIKI_MAIL_TEMPLATE_DIR", "../../template/mail/*")
	config.Load()
	LoadConfig()
	content.LoadConfig()
	InitTemplates()
	mailtpl.InitTemplates()
	conn := testutil.ConnectBackend(true)
	defer conn.Disconnect()
	code := m.Run()
	testutil.CheckOpenConnectionsNull(conn)
	os.Exit(code)
}
