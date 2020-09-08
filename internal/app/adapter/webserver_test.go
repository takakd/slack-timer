package adapter

import (
	"testing"
	"proteinreminder/internal/pkg/testutil"
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
