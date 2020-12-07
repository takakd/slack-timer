package notify

import (
	"context"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Controller implements ControllerHandler.
type Controller struct {
	InputPort notifyevent.InputPort
}

var _ ControllerHandler = (*Controller)(nil)

// NewController create new struct.
func NewController() *Controller {
	h := &Controller{
		InputPort: di.Get("notifyevent.InputPort").(notifyevent.InputPort),
	}
	return h
}

// Handle notifies the event to user.
func (n Controller) Handle(ctx context.Context, input HandleInput) *Response {
	log.Info("handler input", input)

	data := notifyevent.InputData{
		UserID:  input.UserID,
		Message: input.Message,
	}
	err := n.InputPort.NotifyEvent(ctx, data)

	resp := &Response{
		Error: err,
	}

	log.Info("handler output", *resp)

	return resp
}
