package jwt

import (
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
	LoadRSAKeys()
	os.Exit(m.Run())
}
