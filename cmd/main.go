package main

import (
	"context"
	"fmt"
	"os"
	"proteinreminder/internal/app/adapter/webserver"
	"proteinreminder/internal/app/driver/di"
	"proteinreminder/internal/app/driver/di/container"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/errorutil"
	"proteinreminder/internal/pkg/log"
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

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error(errorutil.MakePanicMessage(r))
			os.Exit(1)
		}
		log.Info("exit server")
	}()

	ctx := context.Background()
	log.SetLevel(config.Get("LOG_LEVEL", "debug"))

	setDi()

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
