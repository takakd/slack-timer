package log

import (
	"github.com/golang/mock/gomock"
	"proteinreminder/internal/pkg/testutil"
	"reflect"
	"testing"
)

func TestSetDefaultLogger(t *testing.T) {
	l := GetLogger("")
	SetDefaultLogger("")
	if reflect.TypeOf(l) != reflect.TypeOf(logger) {
		t.Error(testutil.MakeTestMessageWithGotWant(reflect.TypeOf(logger), reflect.TypeOf(l)))
	}
}

func TestDebug(t *testing.T) {
	cases := []struct {
		name  string
		level Level
		msg   string
	}{
		{"OK: debug", LevelDebug, "a b テスト"},
		{"OK: info", LevelInfo, ""},
		{"OK: error", LevelError, ""},
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
		level Level
		msg   string
	}{
		{"OK: debug", LevelDebug, "a b テスト"},
		{"OK: info", LevelInfo, "a b テスト"},
		{"OK: error", LevelError, ""},
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
		level Level
		msg   string
	}{
		{"OK: debug", LevelDebug, "a b テスト"},
		{"OK: info", LevelInfo, "a b テスト"},
		{"OK: error", LevelError, "a b テスト"},
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
