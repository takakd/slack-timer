package config

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"proteinreminder/internal/pkg/fileutil"
	"testing"
)

func TestEnvPath(t *testing.T) {
	t.Run("ok:no env", func(t *testing.T) {
		if envPath() != "" {
			t.Error("must be empty")
		}
	})

	t.Run("ok:exist env", func(t *testing.T) {
		appDir, _ := fileutil.GetAppDir()
		path := filepath.Join(appDir, ".env")
		err := ioutil.WriteFile(path, []byte(""), 0644)
		assert.NoError(t, err)
		assert.NotEmpty(t, envPath())
	})
}

func TestGetConfig(t *testing.T) {
	cases := []struct {
		name         string
		key          string
		defaultValue string
		value        string
		want         string
	}{
		{name: "ok:no params", key: "test1", defaultValue: "other1", value: "value1", want: "value1"},
		{name: "ok:params", key: "test2", defaultValue: "value2", value: "value2", want: "value2"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := NewMockConfig(ctrl)
			m.EXPECT().Get(gomock.Eq(c.key), gomock.Eq(c.defaultValue)).Return(c.value)

			SetConfig(m)

			assert.Equal(t, c.want, Get(c.key, c.defaultValue))
		})
	}
}

// Deprecated
//func TestGet(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	m := NewMockConfig(ctrl)
//	config = m
//
//	key := "key"
//	want := "value"
//	m.EXPECT().Get(key).Return(want)
//
//	got := Get(key)
//	if got != want {
//		t.Error(testutil.MakeTestMessageWithGotWant(got, want))
//	}
//}
