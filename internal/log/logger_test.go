package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"proteinreminder/internal/testutil"
	"regexp"
	"testing"
)

const (
	DateTimePattern = "\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"
)

func TestNewStdoutLogger(t *testing.T) {
	testutil.IsTestCallPanic(func() {
		NewStdoutLogger(LogLevelInfo)
	})
}

func outputLogTest(t *testing.T, loggerLevel LogLevel, level LogLevel, wantPattern string, v ...interface{}) string {

	// Ref: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	l := NewStdoutLogger(loggerLevel)
	l.outputLog(level, v)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	got := buf.String()

	result := ""
	re := regexp.MustCompile(wantPattern)
	if !re.MatchString(got) {
		result = testutil.MakeTestMessageWithGotWant(got, "regexp:"+wantPattern)
	}
	return result
}

func TestStdoutLogger_outputLog(t *testing.T) {
	pattern := fmt.Sprintf("%s %s", DateTimePattern, regexp.QuoteMeta("[INFO] a b テスト"))
	result := outputLogTest(t, LogLevelInfo, LogLevelInfo, pattern, "a", "b", "テスト")
	if result != "" {
		t.Error(result)
	}
}

func TestStdoutLogger_Debug(t *testing.T) {
	dataList := []struct {
		loggerLevel LogLevel
		level       LogLevel
		want        string
	}{
		{LogLevelDebug, LogLevelDebug, "[DEBUG] a b テスト"},
		{LogLevelDebug, LogLevelInfo, "[INFO] a b テスト"},
		{LogLevelDebug, LogLevelError, "[ERROR] a b テスト"},
	}
	for _, v := range dataList {
		pattern := ""
		if v.want != "" {
			pattern = fmt.Sprintf("%s %s", DateTimePattern, regexp.QuoteMeta(v.want))
		}
		result := outputLogTest(t, v.loggerLevel, v.level, pattern, "a", "b", "テスト")
		if result != "" {
			t.Error(result)
		}
	}
}

func TestStdoutLogger_Info(t *testing.T) {
	dataList := []struct {
		loggerLevel LogLevel
		level       LogLevel
		want        string
	}{
		{LogLevelInfo, LogLevelDebug, ""},
		{LogLevelInfo, LogLevelInfo, "[INFO] a b テスト"},
		{LogLevelInfo, LogLevelError, "[ERROR] a b テスト"},
	}
	for _, v := range dataList {
		pattern := ""
		if v.want != "" {
			pattern = fmt.Sprintf("%s %s", DateTimePattern, regexp.QuoteMeta(v.want))
		}
		result := outputLogTest(t, v.loggerLevel, v.level, pattern, "a", "b", "テスト")
		if result != "" {
			t.Error(result)
		}
	}
}

func TestStdoutLogger_Error(t *testing.T) {
	dataList := []struct {
		loggerLevel LogLevel
		level       LogLevel
		want        string
	}{
		{LogLevelError, LogLevelDebug, ""},
		{LogLevelError, LogLevelInfo, ""},
		{LogLevelError, LogLevelError, "[ERROR] a b テスト"},
	}
	for _, v := range dataList {
		pattern := ""
		if v.want != "" {
			pattern = fmt.Sprintf("%s %s", DateTimePattern, regexp.QuoteMeta(v.want))
		}
		result := outputLogTest(t, v.loggerLevel, v.level, pattern, "a", "b", "テスト")
		if result != "" {
			t.Error(result)
		}
	}
}
