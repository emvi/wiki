package client

import (
	"emviwiki/auth/jwt"
	"emviwiki/shared/config"
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	os.Setenv("EMVI_WIKI_TOKEN_PUBLIC_KEY", "../secrets/token.public")
	os.Setenv("EMVI_WIKI_TOKEN_PRIVATE_KEY", "../secrets/token.private")
	config.Load()
	jwt.LoadRSAKeys()
	conn := testutil.ConnectAuth(true)
	defer conn.Disconnect()
	code := m.Run()
	testutil.CheckOpenConnectionsNull(conn)
	os.Exit(code)
}
