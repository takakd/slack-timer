package log

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSetDefaultLogger(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockLogger(ctrl)

		SetDefaultLogger(m)
		assert.Equal(t, reflect.TypeOf(m), reflect.TypeOf(logger))
	})
}

func TestSetLevel(t *testing.T) {
	cases := []struct {
		name     string
		level    string
		logLevel Level
	}{
		{"error", "error", LevelError},
		{"info", "info", LevelInfo},
		{"debug", "debug", LevelDebug},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockLogger(ctrl)
			m.EXPECT().SetLevel(gomock.Eq(c.logLevel))

			SetDefaultLogger(m)

			SetLevel(c.level)
		})
	}
}

func TestDebug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockLogger(ctrl)
	m.EXPECT().Debug(gomock.Eq("test"))

	SetDefaultLogger(m)
	Debug("test")
}

func TestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockLogger(ctrl)
	m.EXPECT().Info(gomock.Eq("test"))

	SetDefaultLogger(m)
	Info("test")
}

func TestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockLogger(ctrl)
	m.EXPECT().Error(gomock.Eq("test"))

	SetDefaultLogger(m)
	Error("test")
}
