package log

import (
	"fmt"
	"strings"
)

type Level int

const (
	LevelError Level = iota
	LevelInfo
	LevelDebug
)

// Output logs interface.
type Logger interface {
	Print(v ...interface{})
}

// Use this interface for logging.
var logger Logger

// Output log above this level.
var logLevel Level

func init() {
	logger = GetLogger("")
	logLevel = LevelInfo
}

// Get logger implementation.
func GetLogger(name string) Logger {
	var l Logger = nil
	if name == "" {
		l = NewStdoutLogger()
	}
	return l
}

// Set default logger which is called log.Info, log.Error...
func SetDefaultLogger(l Logger) {
	logger = l
}

// Set logging level.
func SetLevel(level string) {
	switch level {
	case "error":
		logLevel = LevelError
	case "info":
		logLevel = LevelInfo
	case "debug":
		logLevel = LevelDebug
	}
}

// Output log.
func outputLog(level Level, v ...interface{}) {
	if logLevel < level {
		// Ignore the log with lower priorities than the output level.
		return
	}

	length := len(v)
	if length == 0 {
		return
	}

	var label string
	if level == LevelError {
		label = "ERROR"
	} else if level == LevelInfo {
		label = "INFO"
	} else if level == LevelDebug {
		label = "DEBUG"
	}

	body := strings.Trim(fmt.Sprintf("%s", v...), "[]")
	logger.Print(fmt.Sprintf("[%s] %s", label, body))
}

// Output debug log.
func Debug(v ...interface{}) {
	outputLog(LevelDebug, v)
}

// Output info log.
func Info(v ...interface{}) {
	outputLog(LevelInfo, v)
}

// Output error log.
func Error(v ...interface{}) {
	outputLog(LevelError, v)
}
