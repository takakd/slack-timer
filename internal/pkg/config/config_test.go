package config

import (
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"path/filepath"
	"proteinreminder/internal/pkg/fileutil"
	"proteinreminder/internal/pkg/testutil"
	"runtime"
	"testing"
)

func TestEnvPath(t *testing.T) {
	t.Run("OK: no env", func(t *testing.T) {
		if envPath() != "" {
			t.Error("must be empty")
		}
	})

	t.Run("OK: exist env", func(t *testing.T) {
		appDir, _ := fileutil.GetAppDir()
		path := filepath.Join(appDir, ".env")
		err := ioutil.WriteFile(path, []byte(""), 0644)
		if err != nil {
			t.Error("failed to write test file.")
		}

		t.Log(path)

		if envPath() == "" {
			t.Error("must not be empty")
		}
	})
}

func TestGetConfig(t *testing.T) {
	_, filePath, _, _ := runtime.Caller(0)
	// e.g. $(pwd)/testdata/.env.test
	envPath := filepath.Join(filepath.Dir(filePath), "testdata/.env.test")

	cases := []struct {
		name   string
		params string
	}{
		{name: "OK: no params", params: ""},
		{name: "OK: params", params: envPath},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			called := testutil.DoesTestCallPanic(func() {
				GetConfig("", c.params)
			})
			if called {
				t.Error("must not be called")
			}
		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockConfig(ctrl)
	config = m

	key := "key"
	want := "value"
	m.EXPECT().Get(key).Return(want)

	got := Get(key)
	if got != want {
		t.Error(testutil.MakeTestMessageWithGotWant(got, want))
	}
}
