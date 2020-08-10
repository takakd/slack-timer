package errorutil

import (
	"proteinreminder/internal/testutil"
	"testing"
)

func TestMakePanicMessage(t *testing.T) {
	// OK case
	s := MakePanicMessage("Hi")
	if s != "PANIC: Hi" {
		t.Errorf("got: %s, want: %s", s, "PANIC: Hi")
	}

	// Error case
	ok := testutil.IsTestCallPanic(func() {
		MakePanicMessage(nil)
	})
	if !ok {
		t.Errorf("failed.")
	}
}
