// Package enqueuecontroller provides that events reached the time enqueue queue.
package enqueuecontroller

import (
	"context"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/appinit"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
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

// Lambda callback
// Ref: https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html
func LambdaHandleEvent(ctx context.Context, input LambdaInput) error {
	appinit.AppInit()

	log.Info("handler input", input)

	h := NewEventHandler()
	resp := h.Handler(ctx)

	log.Info("handler output", resp)

	return resp.Error
}
