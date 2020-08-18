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
		NewStdoutLogger(LevelInfo)
	})
}

func outputLogTest(t *testing.T, loggerLevel Level, level Level, wantPattern string, v ...interface{}) string {

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
	result := outputLogTest(t, LevelInfo, LevelInfo, pattern, "a", "b", "テスト")
	if result != "" {
		t.Error(result)
	}
}

func TestStdoutLogger_Debug(t *testing.T) {
	dataList := []struct {
		loggerLevel Level
		level       Level
		want        string
	}{
		{LevelDebug, LevelDebug, "[DEBUG] a b テスト"},
		{LevelDebug, LevelInfo, "[INFO] a b テスト"},
		{LevelDebug, LevelError, "[ERROR] a b テスト"},
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
		loggerLevel Level
		level       Level
		want        string
	}{
		{LevelInfo, LevelDebug, ""},
		{LevelInfo, LevelInfo, "[INFO] a b テスト"},
		{LevelInfo, LevelError, "[ERROR] a b テスト"},
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
		loggerLevel Level
		level       Level
		want        string
	}{
		{LevelError, LevelDebug, ""},
		{LevelError, LevelInfo, ""},
		{LevelError, LevelError, "[ERROR] a b テスト"},
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
