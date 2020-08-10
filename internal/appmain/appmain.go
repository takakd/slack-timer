package appmain

import (
	"os"
	"proteinreminder/internal/errorutil"
	"proteinreminder/internal/ioc"
	"proteinreminder/internal/server"
)

// Run webserver.
func Run() {
	defer func() {
		logger := ioc.GetLogger()
		if r := recover(); r != nil {
			logger.Error(r)
			os.Exit(1)
		}
		logger.Info("exit server.")
	}()

	server := server.NewServer()
	err := server.Init()
	if err != nil {
		panic(errorutil.MakePanicMessage(err))
	}

	server.Run()
}
