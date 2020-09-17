package log

import (
	"log"
	"os"
)

// Stdout logger
type StdoutLogger struct {
	logger *log.Logger
}

func NewStdoutLogger() *StdoutLogger {
	return &StdoutLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *StdoutLogger) Print(v ...interface{}) {
	l.logger.Print(v...)
}
