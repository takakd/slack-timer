// Package enqueuecontroller provides that events reached the time enqueue queue.
package enqueuecontroller

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/driver/di/container/dev"
	"slacktimer/internal/app/driver/di/container/prod"
	"slacktimer/internal/app/driver/di/container/test"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/config/driver"
	"slacktimer/internal/pkg/errorutil"
	"slacktimer/internal/pkg/fileutil"
	"slacktimer/internal/pkg/log"
)

// Lambda handler input data
type LambdaInput struct {
	Version    string   `json:"version"`
	Id         string   `json:"id"`
	DetailType string   `json:"detail-type"`
	Source     string   `json:"source"`
	Account    string   `json:"account"`
	Time       string   `json:"time"`
	Region     string   `json:"region"`
	Resources  []string `json:"resources"`
	Detail     struct {
		EventCategories  []string `json:"EventCategories"`
		SourceType       string   `json:"SourceType"`
		SourceArn        string   `json:"SourceArn"`
		Date             string   `json:"Date"`
		Message          string   `json:"Message"`
		SourceIdentifier string   `json:"SourceIdentifier"`
	} `json:"detail"`
}

type HandlerResponse struct {
	Error error
}

// Create handler corresponding to input.
func NewEventHandler() EventHandler {
	h := &CloudWatchEventHandler{
		usecase: di.Get("enqueuecontroller.EnqueueNotification").(enqueueevent.Usecase),
	}
	return h
}

// Provides handlers to event.
type EventHandler interface {
	Handler(ctx context.Context) *HandlerResponse
}

// Setup config.
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

// Setup DI container by env.
func setDi() {
	env := config.Get("APP_ENV", "dev")

	log.Info(fmt.Sprintf("set di env=%s", env))

	if env == "prod" {
		di.SetDi(&prod.Container{})
	} else if env == "dev" {
		di.SetDi(&dev.Container{})
	} else if env == "test" {
		di.SetDi(&test.Container{})
	}
}

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func LambdaHandleEvent(ctx context.Context, input LambdaInput) error {
	log.Debug(fmt.Sprintf("handler, input=%v", input))

	setConfig()
	setDi()

	h := NewEventHandler()
	resp := h.Handler(ctx)
	return resp.Error
}
