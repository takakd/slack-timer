package adapter

import (
	"proteinreminder/internal/pkg/testutil"
	"testing"
)

func TestNewServer(t *testing.T) {
	called := testutil.DoesTestCallPanic(func() {
		NewWebServer()
	})
	if called {
		t.Errorf("panic")
	}
}

// TODO: TestRun
