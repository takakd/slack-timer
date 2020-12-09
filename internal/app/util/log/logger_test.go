package log

import (
	"reflect"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSetDefaultLogger(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ml := NewMockLogger(ctrl)

		SetDefaultLogger(ml)
		assert.Equal(t, reflect.TypeOf(ml), reflect.TypeOf(logger))
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

			ml := NewMockLogger(ctrl)
			ml.EXPECT().SetLevel(c.logLevel)

			SetDefaultLogger(ml)

			SetLevel(c.level)
		})
	}
}

func TestDebug(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().Debug("test")

	SetDefaultLogger(ml)
	Debug("test")
}

func TestInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().Info("test")

	SetDefaultLogger(ml)
	Info("test")
}

func TestError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().Error("test")

	SetDefaultLogger(ml)
	Error("test")
}

func TestDebugWithContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ac := appcontext.TODO()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().DebugWithContext(ac, "test")

	SetDefaultLogger(ml)
	DebugWithContext(ac, "test")
}

func TestInfoWithContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ac := appcontext.TODO()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().InfoWithContext(ac, "test")

	SetDefaultLogger(ml)
	InfoWithContext(ac, "test")
}

func TestErrorWithContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ac := appcontext.TODO()

	ml := NewMockLogger(ctrl)
	ml.EXPECT().ErrorWithContext(ac, "test")

	SetDefaultLogger(ml)
	ErrorWithContext(ac, "test")
}
