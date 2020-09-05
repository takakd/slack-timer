package adapter

import (
	"proteinreminder/internal/testutil"
	"testing"
)

func TestRun(t *testing.T) {
	ok := testutil.IsTestCallPanic(func() {
		// TODO: correspond to wait server process.
		// Run()
	})
	if ok {
		t.Errorf(testutil.MakeTestMessageWithGotWant(ok, false))
	}
}
