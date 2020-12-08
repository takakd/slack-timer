package driver

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slacktimer/internal/app/util/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCloudWatchLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		NewCloudWatchLogger()
	})
}

func TestCloudWatchLogger_SetLevel(t *testing.T) {
	l := NewCloudWatchLogger()
	l.SetLevel(log.LevelError)
	assert.Equal(t, log.LevelError, l.level)
}

func TestCloudWatchLogger_outputLog(t *testing.T) {
	cases := []struct {
		name  string
		value []interface{}
	}{
		{"ng:marshal", []interface{}{make(chan int)}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := gotTestLogOutput(log.LevelDebug, log.LevelDebug, c.value)
			want := "marshal error in logging\n"
			assert.Equal(t, want, got)
		})
	}
}

func gotTestLogOutput(levelSetting log.Level, level log.Level, msg interface{}) string {
	// Ref: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewCloudWatchLogger()

	logger.SetLevel(levelSetting)

	switch level {
	case log.LevelDebug:
		logger.Debug(msg)
	case log.LevelInfo:
		logger.Info(msg)
	case log.LevelError:
		logger.Error(msg)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	got := buf.String()

	return got
}

func TestCloudWatchLogger_Debug(t *testing.T) {
	cases := []struct {
		name         string
		levelSetting log.Level
		msg          string
	}{
		{"ok:debug", log.LevelDebug, "a b テスト"},
		{"ok:info", log.LevelInfo, ""},
		{"ok:error", log.LevelError, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := gotTestLogOutput(c.levelSetting, log.LevelDebug, c.msg)

			if c.msg == "" {
				assert.Empty(t, got)
			} else {
				want := fmt.Sprintf(`{"level":"DEBUG","msg":"%s"}`+"\n", c.msg)
				assert.Equal(t, want, got)
			}
		})
	}
}

func TestCloudWatchLogger_Info(t *testing.T) {
	cases := []struct {
		name         string
		levelSetting log.Level
		msg          string
	}{
		{"ok:debug", log.LevelDebug, "a b テスト"},
		{"ok:info", log.LevelInfo, "a b テスト"},
		{"ok:error", log.LevelError, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := gotTestLogOutput(c.levelSetting, log.LevelInfo, c.msg)

			if c.msg == "" {
				assert.Empty(t, got)
			} else {
				want := fmt.Sprintf(`{"level":"INFO","msg":"%s"}`+"\n", c.msg)
				assert.Equal(t, want, got)
			}
		})
	}
}

func TestCloudWatchLogger_Error(t *testing.T) {
	cases := []struct {
		name         string
		levelSetting log.Level
		msg          string
	}{
		{"ok:debug", log.LevelDebug, "a b テスト"},
		{"ok:info", log.LevelInfo, "a b テスト"},
		{"ok:error", log.LevelError, "a b テスト"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := gotTestLogOutput(c.levelSetting, log.LevelError, c.msg)

			if c.msg == "" {
				assert.Empty(t, got)
			} else {
				want := fmt.Sprintf(`{"level":"ERROR","msg":"%s"}`+"\n", c.msg)
				assert.Equal(t, want, got)
			}
		})
	}
}
