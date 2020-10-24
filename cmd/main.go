package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slacktimer/internal/app/adapter/webserver"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/driver/di/container"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/config/driver"
	"slacktimer/internal/pkg/errorutil"
	"slacktimer/internal/pkg/fileutil"
	"slacktimer/internal/pkg/log"
)

func setDi() {
	env := config.Get("APP_ENV", "development")

	log.Info(fmt.Sprintf("set di env=%s", env))

	if env == "production" {
		di.SetDi(&container.Production{})
	} else if env == "development" {
		di.SetDi(&container.Development{})
	} else if env == "test" {
		di.SetDi(&container.Test{})
	}
}

func setConfig() {
	configType := os.Getenv("APP_CONFIG_TYPE")
	if configType == "" {
		configType = "env"
	}

	log.Info(fmt.Sprintf("set config type=%s", configType))

	if configType == "env" {
		// Get .env path
		appDir, err := fileutil.GetAppDir()
		if err != nil {
			panic(errorutil.MakePanicMessage("need app directory path."))
		}
		names := make([]string, 0)
		path := filepath.Join(appDir, ".env")
		if fileutil.FileExists(path) {
			names = append(names, path)
		}
		config.SetConfig(driver.NewEnvConfig(names...))
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error(errorutil.MakePanicMessage(r))
			os.Exit(1)
		}
		log.Info("exit server")
	}()

	setConfig()

	setDi()

	ctx := context.Background()
	log.SetLevel(config.Get("LOG_LEVEL", "debug"))

	server := webserver.NewWebServer(ctx)
	if server == nil {
		log.Error("failed to create server")
	}

	// Start web server.
	err := server.Run()
	if err != nil {
		log.Error(err)
	}
}
