package settime

import (
	"slacktimer/internal/app/usecase/timeroffevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// OffEventHandler handles "set" command.
type OffEventHandler interface {
	Handle(ac appcontext.AppContext, data EventCallbackData) *Response
}

// OffEventHandlerFunctor handle "Set" command.
type OffEventHandlerFunctor struct {
	inputPort timeroffevent.InputPort
	presenter *OffEventOutputReceivePresenter
}

var _ OffEventHandler = (*OffEventHandlerFunctor)(nil)

// NewOffEventHandlerFunctor creates new struct.
func NewOffEventHandlerFunctor() *OffEventHandlerFunctor {
	return &OffEventHandlerFunctor{
		inputPort: di.Get("timeroffevent.InputPort").(timeroffevent.InputPort),
		presenter: di.Get("settime.OffEventOutputReceivePresenter").(*OffEventOutputReceivePresenter),
	}
}

// Handle saves event sent by user.
func (se OffEventHandlerFunctor) Handle(ac appcontext.AppContext, data EventCallbackData) *Response {

	log.InfoWithContext(ac, "timer on", "call inputport", map[string]interface{}{
		"user": data.MessageEvent.User,
		"text": data.MessageEvent.Text,
	})

	input := timeroffevent.InputData{
		UserID: data.MessageEvent.User,
	}
	se.inputPort.SetEventOff(ac, input, se.presenter)

	log.InfoWithContext(ac, "return from inputport", se.presenter.Resp)

	return &se.presenter.Resp
}
