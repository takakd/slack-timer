// +build heroku

package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// flag.Parse() here if TestMain uses flags
	os.Exit(m.Run())
}
