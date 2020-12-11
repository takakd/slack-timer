package settime

import (
	"slacktimer/internal/app/usecase/timeronevent"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// OnEventHandler handles "set" command.
type OnEventHandler interface {
	Handle(ac appcontext.AppContext, data EventCallbackData) *Response
}

// OnEventHandlerFunctor handle "Set" command.
type OnEventHandlerFunctor struct {
	inputPort timeronevent.InputPort
	presenter *OnEventOutputReceivePresenter
}

var _ OnEventHandler = (*OnEventHandlerFunctor)(nil)

// NewOnEventHandlerFunctor creates new struct.
func NewOnEventHandlerFunctor() *OnEventHandlerFunctor {
	return &OnEventHandlerFunctor{
		inputPort: di.Get("timeronevent.InputPort").(timeronevent.InputPort),
		presenter: di.Get("settime.OnEventOutputReceivePresenter").(*OnEventOutputReceivePresenter),
	}
}

// Handle saves event sent by user.
func (se OnEventHandlerFunctor) Handle(ac appcontext.AppContext, data EventCallbackData) *Response {

	log.InfoWithContext(ac, "timer on", "call inputport", map[string]interface{}{
		"user": data.MessageEvent.User,
		"text": data.MessageEvent.Text,
	})

	input := timeronevent.InputData{
		UserID: data.MessageEvent.User,
	}
	se.inputPort.SetEventOn(ac, input, se.presenter)

	log.InfoWithContext(ac, "return from inputport", se.presenter.Resp)

	return &se.presenter.Resp
}
