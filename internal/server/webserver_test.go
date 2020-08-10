package server

import (
	"proteinreminder/internal/testutil"
	"testing"
)

func TestNewServer(t *testing.T) {
	called := testutil.IsTestCallPanic(func() {
		NewServer()
	})
	if called {
		t.Errorf("panic")
	}
}

func TestInit(t *testing.T) {
	server := NewServer()
	err := server.Init()
	if err != nil {
		t.Errorf("failed.")
	}
}

// TODO: TestRun
