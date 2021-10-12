package pages

import (
	"emviwiki/shared/config"
	"emviwiki/shared/recaptcha"
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	os.Setenv("EMVI_WIKI_RECAPTCHA_CLIENT_SECRET", "test")
	os.Setenv("EMVI_WIKI_TEMPLATE_DIR", "../../template/auth/*")
	os.Setenv("EMVI_WIKI_MAIL_TEMPLATE_DIR", "../../template/mail/*")
	config.Load()
	recaptcha.LoadConfig()
	InitTemplates()
	conn := testutil.ConnectAuth(true)
	defer conn.Disconnect()
	os.Exit(m.Run())
}
