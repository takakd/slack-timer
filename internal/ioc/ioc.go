package ioc

import "proteinreminder/internal/log"

var logger log.Logger

func GetLogger() log.Logger {
	if logger == nil {
		logger = log.NewStdoutLogger(log.LogLevelDebug)
	}
	return logger
}
