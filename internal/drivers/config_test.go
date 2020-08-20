package drivers

import (
	"os"
	"proteinreminder/internal/testutil"
	"testing"
)

func TestEnvConfig_Get(t *testing.T) {
	patterns := [][]string{
		// key, value, want
		{"NAME1", "value1", "value1"},
		{"NAME2", "", ""},
		{"", "value2", ""},
	}
	for _, pattern := range patterns {
		os.Setenv(pattern[0], pattern[1])
	}

	config := NewEnvConfig()
	for _, pattern := range patterns {
		got := config.Get(pattern[0])
		if got != pattern[2] {
			t.Errorf(testutil.MakeTestMessageWithGotWant(got, pattern[2]))
		}
	}
}
