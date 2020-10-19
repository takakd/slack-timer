package main

import (
	"context"
	"os"
	"proteinreminder/internal/app/adapter/webserver"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/errorutil"
	"proteinreminder/internal/pkg/log"
)

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
