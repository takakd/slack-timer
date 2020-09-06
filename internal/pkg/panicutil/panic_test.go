package panicutil

import (
	"proteinreminder/internal/pkg/testutil"
	"testing"
)

func TestMakePanicMessage(t *testing.T) {
	// ok case
	s := MakePanicMessage("Hi")
	if s != "PANIC: Hi" {
		t.Errorf("got: %s, want: %s", s, "PANIC: Hi")
	}

	// error case
	ok := testutil.DoesTestCallPanic(func() {
		MakePanicMessage(nil)
	})
	if !ok {
		t.Errorf("failed.")
	}
}
