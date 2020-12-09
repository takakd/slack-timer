// Package enqueue provides that events reached the time enqueue queue.
package enqueue

import (
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

// Controller implements ControllerHandler.
type Controller struct {
	InputPort enqueueevent.InputPort
}

var _ ControllerHandler = (*Controller)(nil)

// NewController create new struct.
func NewController() *Controller {
	h := &Controller{
		InputPort: di.Get("enqueueevent.InputPort").(enqueueevent.InputPort),
	}
	return h
}

// Handle enqueues events reached the time.
func (e Controller) Handle(ac appcontext.AppContext, input HandleInput) {
	// TODO: Getting time from Lambda context?
	data := enqueueevent.InputData{
		EventTime: time.Now().UTC(),
	}

	log.InfoWithContext(ac, "call inputport", input)
	e.InputPort.EnqueueEvent(ac, data)
	log.InfoWithContext(ac, "return from inputport")
}
