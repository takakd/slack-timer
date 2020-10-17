// Deprecated
package config

import (
	"os"
	"path/filepath"
	"proteinreminder/internal/pkg/testutil"
	"runtime"
	"testing"
)

func TestNewEnvConfig(t *testing.T) {
	called := testutil.DoesTestCallPanic(func() {
		NewEnvConfig()
	})
	if called {
		t.Error("must not be called")
	}

	called = testutil.DoesTestCallPanic(func() {
		NewEnvConfig("")
	})
	if !called {
		t.Error("must be called")
	}
}

func TestEnvConfig_Get(t *testing.T) {
	cases := []struct {
		name string
		env  string
		key  string
		want string
	}{
		{"OK: env1", ".env.test", "NAME1", "value1"},
		{"OK: env2", ".env.test", "NAME2", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, filePath, _, _ := runtime.Caller(0)
			// e.g. $(pwd)/testdata/.env.test
			envPath := filepath.Join(filepath.Dir(filePath), "testdata/", c.env)
			config := NewEnvConfig(envPath)

			got := config.Get(c.key, "")
			if c.want != got {
				t.Error(testutil.MakeTestMessageWithGotWant(got, c.want))
			}
		})
	}

	t.Run("OK: not effect os.env", func(t *testing.T) {
		want := "not define in .env"
		os.Setenv("NAME3", want)

		_, filePath, _, _ := runtime.Caller(0)
		envPath := filepath.Join(filepath.Dir(filePath), "testdata/.env.test")
		config := NewEnvConfig(envPath)

		got := config.Get("NAME3", "")
		if want != got {
			t.Error(testutil.MakeTestMessageWithGotWant(got, want))
		}
	})

	t.Run("OK: overwrite", func(t *testing.T) {
		_, filePath, _, _ := runtime.Caller(0)
		envPath := filepath.Join(filepath.Dir(filePath), "testdata/.env.test")
		config := NewEnvConfig(envPath)

		want := "changed"
		os.Setenv("NAME1", want)

		got := config.Get("NAME1", "")
		if want != got {
			t.Error(testutil.MakeTestMessageWithGotWant(got, want))
		}
	})
}
