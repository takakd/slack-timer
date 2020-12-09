package driver

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slacktimer/internal/app/util/log"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
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
			got := gotTestLogOutput(nil, log.LevelDebug, log.LevelDebug, c.value)
			want := "marshal error in logging\n"
			assert.Equal(t, want, got)
		})
	}
}

func gotTestLogOutput(ac appcontext.AppContext, levelSetting log.Level, level log.Level, msg interface{}) string {
	// Ref: https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := NewCloudWatchLogger()

	logger.SetLevel(levelSetting)

	if ac == nil {
		switch level {
		case log.LevelDebug:
			logger.Debug(msg)
		case log.LevelInfo:
			logger.Info(msg)
		case log.LevelError:
			logger.Error(msg)
		}
	} else {
		switch level {
		case log.LevelDebug:
			logger.DebugWithContext(ac, msg)
		case log.LevelInfo:
			logger.InfoWithContext(ac, msg)
		case log.LevelError:
			logger.ErrorWithContext(ac, msg)
		}
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
		withContext  bool
	}{
		{"ok:debug", log.LevelDebug, "a b テスト", false},
		{"ok:info", log.LevelInfo, "", false},
		{"ok:error", log.LevelError, "", false},
		{"ok:debug with ctx", log.LevelDebug, "a b テスト", true},
		{"ok:info with ctx", log.LevelInfo, "", true},
		{"ok:error with ctx", log.LevelError, "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.withContext {
				lc := &lambdacontext.LambdaContext{
					AwsRequestID: "test ID",
				}
				ctx := lambdacontext.NewContext(context.TODO(), lc)
				ac, _ := appcontext.FromContext(ctx)

				got := gotTestLogOutput(ac, c.levelSetting, log.LevelDebug, c.msg)

				if c.msg == "" {
					assert.Empty(t, got)
				} else {
					want := fmt.Sprintf(`{"AwsRequestID":"test ID","level":"DEBUG","msg":"%s"}`+"\n", c.msg)
					assert.Equal(t, want, got)
				}
			} else {
				got := gotTestLogOutput(nil, c.levelSetting, log.LevelDebug, c.msg)

				if c.msg == "" {
					assert.Empty(t, got)
				} else {
					want := fmt.Sprintf(`{"level":"DEBUG","msg":"%s"}`+"\n", c.msg)
					assert.Equal(t, want, got)
				}
			}
		})
	}
}

func TestCloudWatchLogger_Info(t *testing.T) {
	cases := []struct {
		name         string
		levelSetting log.Level
		msg          string
		withContext  bool
	}{
		{"ok:debug", log.LevelDebug, "a b テスト", false},
		{"ok:info", log.LevelInfo, "a b テスト", false},
		{"ok:error", log.LevelError, "", false},
		{"ok:debug with context", log.LevelDebug, "a b テスト", true},
		{"ok:info with context", log.LevelInfo, "a b テスト", true},
		{"ok:error with context", log.LevelError, "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.withContext {
				lc := &lambdacontext.LambdaContext{
					AwsRequestID: "test ID",
				}
				ctx := lambdacontext.NewContext(context.TODO(), lc)
				ac, _ := appcontext.FromContext(ctx)

				got := gotTestLogOutput(ac, c.levelSetting, log.LevelInfo, c.msg)

				if c.msg == "" {
					assert.Empty(t, got)
				} else {
					want := fmt.Sprintf(`{"AwsRequestID":"test ID","level":"INFO","msg":"%s"}`+"\n", c.msg)
					assert.Equal(t, want, got)
				}
			} else {
				got := gotTestLogOutput(nil, c.levelSetting, log.LevelInfo, c.msg)

				if c.msg == "" {
					assert.Empty(t, got)
				} else {
					want := fmt.Sprintf(`{"level":"INFO","msg":"%s"}`+"\n", c.msg)
					assert.Equal(t, want, got)
				}
			}
		})
	}
}

func TestCloudWatchLogger_Error(t *testing.T) {
	cases := []struct {
		name         string
		levelSetting log.Level
		msg          string
		withContext  bool
	}{
		{"ok:debug", log.LevelDebug, "a b テスト", false},
		{"ok:info", log.LevelInfo, "a b テスト", false},
		{"ok:error", log.LevelError, "a b テスト", false},
		{"ok:debug with context", log.LevelDebug, "a b テスト", true},
		{"ok:info with context", log.LevelInfo, "a b テスト", true},
		{"ok:error with context", log.LevelError, "a b テスト", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.withContext {
				lc := &lambdacontext.LambdaContext{
					AwsRequestID: "test ID",
				}
				ctx := lambdacontext.NewContext(context.TODO(), lc)
				ac, _ := appcontext.FromContext(ctx)

				got := gotTestLogOutput(ac, c.levelSetting, log.LevelError, c.msg)

				if c.msg == "" {
					assert.Empty(t, got)
				} else {
					want := fmt.Sprintf(`{"AwsRequestID":"test ID","level":"ERROR","msg":"%s"}`+"\n", c.msg)
					assert.Equal(t, want, got)
				}
			} else {
				got := gotTestLogOutput(nil, c.levelSetting, log.LevelError, c.msg)

				if c.msg == "" {
					assert.Empty(t, got)
				} else {
					want := fmt.Sprintf(`{"level":"ERROR","msg":"%s"}`+"\n", c.msg)
					assert.Equal(t, want, got)
				}
			}
		})
	}
}
