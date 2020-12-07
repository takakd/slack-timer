// Package driver provides implementation of logger.
package driver

import (
	"fmt"
	"log"
	"os"
	log2 "slacktimer/internal/app/util/log"
	"strings"
)

// CloudWatchLogger implements log.Logger with CloudwatchLogs.
type CloudWatchLogger struct {
	logger *log.Logger
	level  log2.Level
}

var _ log2.Logger = (*CloudWatchLogger)(nil)

// NewCloudWatchLogger create new struct.
func NewCloudWatchLogger() *CloudWatchLogger {
	return &CloudWatchLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// SetLevel sets logging level.
func (l *CloudWatchLogger) SetLevel(level log2.Level) {
	l.level = level
}

func (l CloudWatchLogger) outputLog(level log2.Level, v ...interface{}) {
	if l.level < level {
		// Ignore the log with lower priorities than the output level.
		return
	}

	length := len(v)
	if length == 0 {
		return
	}

	var label string
	if level == log2.LevelError {
		label = "ERROR"
	} else if level == log2.LevelInfo {
		label = "INFO"
	} else if level == log2.LevelDebug {
		label = "DEBUG"
	}

	body := strings.Trim(fmt.Sprintf("%s", v...), "[]")
	l.logger.Print(fmt.Sprintf("[%s] %s", label, body))
}

// Debug implements Logger.Debug.
func (l CloudWatchLogger) Debug(v ...interface{}) {
	l.outputLog(log2.LevelDebug, v)
}

// Info implements Logger.Info.
func (l CloudWatchLogger) Info(v ...interface{}) {
	l.outputLog(log2.LevelInfo, v)
}

// Error implements Logger.Error.
func (l CloudWatchLogger) Error(v ...interface{}) {
	l.outputLog(log2.LevelError, v)
}
