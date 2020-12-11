package notify

import (
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Controller implements ControllerHandler.
type Controller struct {
	InputPort notifyevent.InputPort
}

var _ ControllerHandler = (*Controller)(nil)

// NewController creates new struct.
func NewController() *Controller {
	h := &Controller{
		InputPort: di.Get("notifyevent.InputPort").(notifyevent.InputPort),
		// TODO: presenter get from di
	}
	return h
}

// Handle notifies the event to user.
func (n Controller) Handle(ac appcontext.AppContext, input HandleInput) *Response {
	log.InfoWithContext(ac, "call inputport", input)

	data := notifyevent.InputData{
		UserID:  input.UserID,
		Message: input.Message,
	}

	// Receive error to send error state to SQS.
	err := n.InputPort.NotifyEvent(ac, data)

	log.InfoWithContext(ac, "return from inputport", err)

	resp := &Response{
		Error: err,
	}

	log.InfoWithContext(ac, "handler output", *resp)

	return resp
}
