// Package enqueuecontroller provides that events reached the time enqueue queue.
package enqueuecontroller

import (
	"context"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

type Handler interface {
	Handler(ctx context.Context, input HandlerInput)
}

type CloudWatchEventHandler struct {
	InputPort enqueueevent.InputPort
}

type HandlerInput struct {
}

func NewHandler() Handler {
	h := &CloudWatchEventHandler{
		InputPort: di.Get("enqueueevent.InputPort").(enqueueevent.InputPort),
	}
	return h
}

func (c CloudWatchEventHandler) Handler(ctx context.Context, input HandlerInput) {
	log.Info("handler called", input)

	data := enqueueevent.InputData{
		EventTime: time.Now().UTC(),
	}

	c.InputPort.EnqueueEvent(ctx, data)

	log.Info("handler done")
}
