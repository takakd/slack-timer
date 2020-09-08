// +build server

package main

func main() {
	logger := ioc.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error(r)
			os.Exit(1)
		}
		logger.Info("exit server.")
	}()

	loadEnv()

	server := server.NewServer()
	err := server.Init()
	if err != nil {
		panic(errorutil.MakePanicMessage(err))
	}

	err = server.Run()
	if err != nil {
		logger.Error(err)
	}
}
