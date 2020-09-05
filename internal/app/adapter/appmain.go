package adapter

import (
	"context"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"proteinreminder/internal/errorutil"
	"proteinreminder/internal/fileutil"
	"proteinreminder/internal/ioc"
	"proteinreminder/internal/server"
	"runtime"
)

func loadEnv() {
	logger := ioc.GetLogger()

	_, filename, _, _ := runtime.Caller(0)
	appDir := filepath.Dir(filename)
	envPath := filepath.Join(appDir, "../configs/.env")
	if fileutil.FileExists(envPath) {
		logger.Info(".env found. loaded it.")
		godotenv.Load(envPath)
	}
}

func Run(ctx context.Context) {
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
