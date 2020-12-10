// Package driver provides implementation of logger.
package driver

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"slacktimer/internal/app/util/appcontext"
	log2 "slacktimer/internal/app/util/log"
)

// CloudWatchLogger implements log.Logger with CloudWatchLogs.
type CloudWatchLogger struct {
	logger *log.Logger
	level  log2.Level
}

var _ log2.Logger = (*CloudWatchLogger)(nil)

// NewCloudWatchLogger create new struct.
func NewCloudWatchLogger() *CloudWatchLogger {
	return &CloudWatchLogger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// SetLevel sets logging level.
func (l *CloudWatchLogger) SetLevel(level log2.Level) {
	l.level = level
}

func (l CloudWatchLogger) outputLog(ac appcontext.AppContext, level log2.Level, v []interface{}) {
	if l.level < level {
		// Ignore the log with lower priorities than the output level.
		return
	}

	length := len(v)
	if length == 0 {
		return
	}

	var label string
	switch level {
	case log2.LevelError:
		label = "ERROR"
	case log2.LevelInfo:
		label = "INFO"
	case log2.LevelDebug:
		label = "DEBUG"
	}

	data := map[string]interface{}{
		"level": label,
		"msg":   v,
	}

	// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-context.html
	if ac != nil {
		data["AwsRequestID"] = ac.RequestID()
	}

	if length == 1 {
		data["msg"] = v[0]
	}

	var msg string
	if body, err := json.Marshal(data); err == nil {
		var buf bytes.Buffer
		if json.Compact(&buf, body); err == nil {
			msg = buf.String()
		} else {
			msg = "marshal error in logging"
		}
	} else {
		msg = "marshal error in logging"
	}

	l.logger.Print(msg)
}

// Debug implements Logger.Debug.
func (l CloudWatchLogger) Debug(v ...interface{}) {
	l.outputLog(nil, log2.LevelDebug, v)
}

// Info implements Logger.Info.
func (l CloudWatchLogger) Info(v ...interface{}) {
	l.outputLog(nil, log2.LevelInfo, v)
}

// Error implements Logger.Error.
func (l CloudWatchLogger) Error(v ...interface{}) {
	l.outputLog(nil, log2.LevelError, v)
}

// DebugWithContext implements Logger.DebugWithContext.
func (l CloudWatchLogger) DebugWithContext(ac appcontext.AppContext, v ...interface{}) {
	l.outputLog(ac, log2.LevelDebug, v)
}

// InfoWithContext implements Logger.InfoWithContext.
func (l CloudWatchLogger) InfoWithContext(ac appcontext.AppContext, v ...interface{}) {
	l.outputLog(ac, log2.LevelInfo, v)
}

// ErrorWithContext implements Logger.ErrorWithContext.
func (l CloudWatchLogger) ErrorWithContext(ac appcontext.AppContext, v ...interface{}) {
	l.outputLog(ac, log2.LevelError, v)
}
