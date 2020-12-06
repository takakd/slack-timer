// Package enqueuecontroller provides that events reached the time enqueue queue.
package enqueue

import (
	"context"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

// Concrete struct.
type EnqueueController struct {
	InputPort enqueueevent.InputPort
}

var _ Controller = (*EnqueueController)(nil)

func NewEnqueueController() *EnqueueController {
	h := &EnqueueController{
		InputPort: di.Get("enqueueevent.InputPort").(enqueueevent.InputPort),
	}
	return h
}

func (e EnqueueController) Handle(ctx context.Context, input HandleInput) {
	log.Info("handler called", input)

	// TODO: Getting time from Lambda context?
	data := enqueueevent.InputData{
		EventTime: time.Now().UTC(),
	}

	e.InputPort.EnqueueEvent(ctx, data)

	log.Info("handler done")
}
