package driver

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
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

//func TestCloudWatchLogger_outputLog(t *testing.T) {
//	cases := []struct {
//		name string
//		want string
//	}{
//		{name: "ok", want: "test log."},
//		{name: "ng", want: ""},
//	}
//	for _, c := range cases {
//		t.Run(c.name, func(t *testing.T) {
//
//			pattern := ""
//			if c.want != "" {
//				dateTimePattern := "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"
//				pattern = fmt.Sprintf("%s %s", dateTimePattern, regexp.QuoteMeta(c.want))
//			}
//
//			// Ref: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
//			old := os.Stdout
//			r, w, _ := os.Pipe()
//			os.Stdout = w
//
//			logger := NewCloudWatchLogger()
//			logger.outputLog(c.want)
//
//			w.Close()
//			os.Stdout = old
//
//			var buf bytes.Buffer
//			io.Copy(&buf, r)
//			got := buf.String()
//
//			re := regexp.MustCompile(pattern)
//			assert.True(t, re.MatchString(got))
//		})
//	}
//}

func gotTestLogOutput(levelSetting log.Level, level log.Level, msg string) string {
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
				msg := fmt.Sprintf("[DEBUG] %s", c.msg)
				dateTimePattern := "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"
				want := fmt.Sprintf("%s %s", dateTimePattern, regexp.QuoteMeta(msg))
				re := regexp.MustCompile(want)

				assert.True(t, re.MatchString(got))
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
				msg := fmt.Sprintf("[INFO] %s", c.msg)
				dateTimePattern := "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"
				want := fmt.Sprintf("%s %s", dateTimePattern, regexp.QuoteMeta(msg))
				re := regexp.MustCompile(want)

				assert.True(t, re.MatchString(got))
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
				msg := fmt.Sprintf("[ERROR] %s", c.msg)
				dateTimePattern := "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"
				want := fmt.Sprintf("%s %s", dateTimePattern, regexp.QuoteMeta(msg))
				re := regexp.MustCompile(want)

				assert.True(t, re.MatchString(got))
			}
		})
	}
}
