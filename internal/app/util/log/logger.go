// Package log provides logging feature.
package log

import "context"

// Level represents the type of logging level.
type Level int

const (
	// LevelError outputs Error logs.
	LevelError Level = iota
	// LevelInfo outputs Error and Info logs.
	LevelInfo
	// LevelDebug outputs all level logs.
	LevelDebug
)

// Logger defines logging methods.
type Logger interface {
	SetLevel(level Level)
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
	DebugWithContext(ctx context.Context, v ...interface{})
	InfoWithContext(ctx context.Context, v ...interface{})
	ErrorWithContext(ctx context.Context, v ...interface{})
}

// Use this interface for logging.
var logger Logger

// Output log above this level.
var logLevel Level

// SetDefaultLogger sets default logger which is called log.Info, log.Error...
func SetDefaultLogger(l Logger) {
	logger = l
}

// SetLevel sets logging level.
func SetLevel(level string) {
	if logger == nil {
		return
	}

	switch level {
	case "error":
		logLevel = LevelError
	case "info":
		logLevel = LevelInfo
	case "debug":
		logLevel = LevelDebug
	}
	logger.SetLevel(logLevel)
}

// Debug outputs debug log.
func Debug(v ...interface{}) {
	defer func() {
		// don't panic
	}()

	if logger != nil {
		logger.Debug(v...)
	}
}

// Info outputs info log.
func Info(v ...interface{}) {
	defer func() {
		// don't panic
	}()

	if logger != nil {
		logger.Info(v...)
	}
}

// Error outputs info log.
func Error(v ...interface{}) {
	defer func() {
		// don't panic
	}()

	if logger != nil {
		logger.Error(v...)
	}
}

// DebugWithContext outputs debug log with context information.
func DebugWithContext(ctx context.Context, v ...interface{}) {
	defer func() {
		// don't panic
	}()

	if logger != nil {
		logger.DebugWithContext(ctx, v...)
	}
}

// InfoWithContext outputs info log with context information.
func InfoWithContext(ctx context.Context, v ...interface{}) {
	defer func() {
		// don't panic
	}()

	if logger != nil {
		logger.InfoWithContext(ctx, v...)
	}
}

// ErrorWithContext outputs info log with context information.
func ErrorWithContext(ctx context.Context, v ...interface{}) {
	defer func() {
		// don't panic
	}()

	if logger != nil {
		logger.ErrorWithContext(ctx, v...)
	}
}
