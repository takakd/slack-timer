package log

type Level int

// Log level
const (
	LevelError Level = iota
	LevelInfo
	LevelDebug
)

// Output logs interface
type Logger interface {
	SetLevel(level Level)
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{})
}

// Use this interface for logging.
var logger Logger

// Output log above this level.
var logLevel Level

// Set default logger which is called log.Info, log.Error...
func SetDefaultLogger(l Logger) {
	logger = l
}

// Set logging level.
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

// Output debug log.
func Debug(v ...interface{}) {
	if logger != nil {
		logger.Debug(v...)
	}
}

// Output info log.
func Info(v ...interface{}) {
	if logger != nil {
		logger.Info(v...)
	}
}

// Output error log.
func Error(v ...interface{}) {
	if logger != nil {
		logger.Error(v...)
	}
}
