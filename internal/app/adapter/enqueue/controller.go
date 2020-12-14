// Package enqueue provides that events reached the time enqueue queue.
package enqueue

import (
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Controller implements ControllerHandler.
type Controller struct {
	inputPort enqueueevent.InputPort
}

var _ ControllerHandler = (*Controller)(nil)

// NewController creates new struct.
func NewController() *Controller {
	h := &Controller{
		inputPort: di.Get("enqueueevent.InputPort").(enqueueevent.InputPort),
	}
	return h
}

// Handle enqueues events reached the time.
func (e Controller) Handle(ac appcontext.AppContext, input HandleInput) {
	data := enqueueevent.InputData{
		EventTime: ac.HandlerCalledTime(),
	}

	log.InfoWithContext(ac, "call inputport", input)
	e.inputPort.EnqueueEvent(ac, data)
	log.InfoWithContext(ac, "return from inputport")
}
