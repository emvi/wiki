package feed

import (
	"emviwiki/shared/testutil"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testutil.SetTestLogger()
	code := m.Run()
	os.Exit(code)
}
