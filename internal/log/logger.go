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

type Level int

const (
	LevelError Level = iota
	LevelInfo
	LevelDebug
)

//
// Stdout Logger
//

// Stdout logger
type StdoutLogger struct {
	level  Level
	logger *log.Logger
}

// Create stdout logger.
func NewStdoutLogger(outputLoglevel Level) *StdoutLogger {
	return &StdoutLogger{
		level:  outputLoglevel,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Output log.
func (l *StdoutLogger) outputLog(level Level, v ...interface{}) {
	if l.level < level {
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
	l.logger.Printf("[%s] %s", label, body)
}

// Output debug log.
func (l *StdoutLogger) Debug(v ...interface{}) {
	l.outputLog(LevelDebug, v)
}

// Output info log.
func (l *StdoutLogger) Info(v ...interface{}) {
	l.outputLog(LevelInfo, v)
}

// Output error log.
func (l *StdoutLogger) Error(v ...interface{}) {
	l.outputLog(LevelError, v)
}
