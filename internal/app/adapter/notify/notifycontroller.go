package notify

import (
	"context"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type NotifyController struct {
	InputPort notifyevent.InputPort
}

func NewNotifyController() Controller {
	h := &NotifyController{
		InputPort: di.Get("notifycontroller.InputPort").(notifyevent.InputPort),
	}
	return h
}

func (s NotifyController) Handle(ctx context.Context, input HandleInput) *Response {
	log.Info("handler input", input)

	data := notifyevent.InputData{
		UserId:  input.UserId,
		Message: input.Message,
	}

	err := s.InputPort.NotifyEvent(ctx, data)

	resp := &Response{
		Error: err,
	}

	log.Info("handler output", *resp)

	return resp
}
