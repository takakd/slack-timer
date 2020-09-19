package main

import (
	"context"
	"os"
	"proteinreminder/internal/app/adapter"
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

	server := adapter.NewWebServer(ctx, config.GetConfig(""))
	if server == nil {
		log.Error("failed to create server")
	}

	err := server.Run()
	if err != nil {
		log.Error(err)
	}
}
