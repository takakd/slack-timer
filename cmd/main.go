package main

import (
	"context"
	"os"
	"proteinreminder/internal/app/adapter"
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

	server := adapter.NewWebServer()
	if server == nil {
		log.Error("failed to create server")
	}
	ctx := context.Background()

	err := server.Run(ctx)
	if err != nil {
		log.Error(err)
	}
}
