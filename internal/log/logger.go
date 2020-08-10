package log

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Logging interface in app.
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelInfo
	LogLevelDebug
)

//
// Stdout Logger
//

// Stdout logger
type StdoutLogger struct {
	level  LogLevel
	logger *log.Logger
}

// Create stdout logger.
func NewStdoutLogger(outputLoglevel LogLevel) *StdoutLogger {
	return &StdoutLogger{
		level:  outputLoglevel,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Output log.
func (l *StdoutLogger) outputLog(level LogLevel, v ...interface{}) {
	if l.level < level {
		// Ignore the log with lower priorities than the output level.
		return
	}

	length := len(v)
	if length == 0 {
		return
	}

	var label string
	if level == LogLevelError {
		label = "ERROR"
	} else if level == LogLevelInfo {
		label = "INFO"
	} else if level == LogLevelDebug {
		label = "DEBUG"
	}

	body := strings.Trim(fmt.Sprintf("%s", v...), "[]")
	l.logger.Printf("[%s] %s", label, body)
}

// Output debug log.
func (l *StdoutLogger) Debug(v ...interface{}) {
	l.outputLog(LogLevelDebug, v)
}

// Output info log.
func (l *StdoutLogger) Info(v ...interface{}) {
	l.outputLog(LogLevelInfo, v)
}

// Output error log.
func (l *StdoutLogger) Error(v ...interface{}) {
	l.outputLog(LogLevelError, v)
}
