package driver

import (
	"fmt"
	"log"
	"os"
	log2 "slacktimer/internal/app/util/log"
	"strings"
)

// Stdout logger
type CloudWatchLogger struct {
	logger *log.Logger
	level  log2.Level
}

func NewCloudWatchLogger() *CloudWatchLogger {
	return &CloudWatchLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Set logging level.
func (l *CloudWatchLogger) SetLevel(level log2.Level) {
	l.level = level
}

// Output log.
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

// Output debug log.
func (l CloudWatchLogger) Debug(v ...interface{}) {
	l.outputLog(log2.LevelDebug, v)
}

// Output info log.
func (l CloudWatchLogger) Info(v ...interface{}) {
	l.outputLog(log2.LevelInfo, v)
}

// Output error log.
func (l CloudWatchLogger) Error(v ...interface{}) {
	l.outputLog(log2.LevelError, v)
}
