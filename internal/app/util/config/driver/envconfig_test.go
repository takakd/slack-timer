// Deprecated
package driver

import (
	"os"
	"path/filepath"
	"runtime"
	"slacktimer/internal/pkg/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnvConfig(t *testing.T) {
	called := helper.DoesTestCallPanic(func() {
		NewEnvConfig()
	})
	assert.False(t, called)

	called = helper.DoesTestCallPanic(func() {
		NewEnvConfig("")
	})
	assert.True(t, called)
}

func TestEnvConfig_Get(t *testing.T) {
	cases := []struct {
		name string
		env  string
		key  string
		want string
	}{
		{"ok:env1", ".env.test", "NAME1", "value1"},
		{"ok:env2", ".env.test", "NAME2", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, filePath, _, _ := runtime.Caller(0)
			// e.g. $(pwd)/testdata/.env.test
			envPath := filepath.Join(filepath.Dir(filePath), "../testdata/", c.env)
			config := NewEnvConfig(envPath)

			got := config.Get(c.key, "")
			assert.Equal(t, c.want, got)
		})
	}

	t.Run("ok:not effect os.env", func(t *testing.T) {
		want := "not define in .env"
		os.Setenv("NAME3", want)

		_, filePath, _, _ := runtime.Caller(0)
		envPath := filepath.Join(filepath.Dir(filePath), "../testdata/.env.test")
		config := NewEnvConfig(envPath)

		got := config.Get("NAME3", "")
		assert.Equal(t, want, got)
	})

	t.Run("ok:overwrite", func(t *testing.T) {
		_, filePath, _, _ := runtime.Caller(0)
		envPath := filepath.Join(filepath.Dir(filePath), "../testdata/.env.test")
		config := NewEnvConfig(envPath)

		want := "changed"
		os.Setenv("NAME1", want)

		got := config.Get("NAME1", "")
		assert.Equal(t, want, got)
	})
}
