package log

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSetDefaultLogger(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		l := GetLogger("")
		SetDefaultLogger(l)
		assert.Equal(t, reflect.TypeOf(l), reflect.TypeOf(logger))
	})
}

func TestDebug(t *testing.T) {
	cases := []struct {
		name  string
		level string
		msg   string
	}{
		{"ok:debug", "debug", "a b テスト"},
		{"ok:info", "info", ""},
		{"ok:error", "error", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SetLevel(c.level)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockLogger(ctrl)
			logger = m

			if c.msg == "" {
				m.EXPECT().Print().MaxTimes(0)
			} else {
				m.EXPECT().Print(gomock.Eq("[DEBUG] " + c.msg))
			}

			Debug(c.msg)
		})
	}
}

func TestInfo(t *testing.T) {
	cases := []struct {
		name  string
		level string
		msg   string
	}{
		{"ok:debug", "debug", "a b テスト"},
		{"ok:info", "info", "a b テスト"},
		{"ok:error", "error", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SetLevel(c.level)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockLogger(ctrl)
			logger = m

			if c.msg == "" {
				m.EXPECT().Print().MaxTimes(0)
			} else {
				m.EXPECT().Print(gomock.Eq("[INFO] " + c.msg))
			}

			Info(c.msg)
		})
	}
}

func TestError(t *testing.T) {
	cases := []struct {
		name  string
		level string
		msg   string
	}{
		{"ok:debug", "debug", "a b テスト"},
		{"ok:info", "info", "a b テスト"},
		{"ok:error", "error", "a b テスト"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			SetLevel(c.level)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockLogger(ctrl)
			logger = m

			m.EXPECT().Print(gomock.Eq("[ERROR] " + c.msg))

			Error(c.msg)
		})
	}
}
