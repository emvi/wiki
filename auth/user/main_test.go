package user

import (
	"emviwiki/auth/jwt"
	"emviwiki/shared/config"
	"emviwiki/shared/recaptcha"
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	os.Setenv("EMVI_WIKI_RECAPTCHA_CLIENT_SECRET", "test")
	os.Setenv("EMVI_WIKI_TOKEN_PUBLIC_KEY", "../secrets/token.public")
	os.Setenv("EMVI_WIKI_TOKEN_PRIVATE_KEY", "../secrets/token.private")
	os.Setenv("EMVI_WIKI_MAIL_TEMPLATE_DIR", "../../template/mail/*")
	config.Load()
	jwt.LoadRSAKeys()
	InitTemplates()
	recaptcha.LoadConfig()
	conn := testutil.ConnectAuth(true)
	defer conn.Disconnect()
	code := m.Run()
	testutil.CheckOpenConnectionsNull(conn)
	os.Exit(code)
}
