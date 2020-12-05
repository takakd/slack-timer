// Package enqueuecontroller provides that events reached the time enqueue queue.
package enqueue

import (
	"context"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

type EnqueueController struct {
	InputPort enqueueevent.InputPort
}

func NewEnqueueController() Controller {
	h := &EnqueueController{
		InputPort: di.Get("enqueueevent.InputPort").(enqueueevent.InputPort),
	}
	return h
}

func (c EnqueueController) Handle(ctx context.Context, input HandleInput) {
	log.Info("handler called", input)

	data := enqueueevent.InputData{
		EventTime: time.Now().UTC(),
	}

	c.InputPort.EnqueueEvent(ctx, data)

	log.Info("handler done")
}
