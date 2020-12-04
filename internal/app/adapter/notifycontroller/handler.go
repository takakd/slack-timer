package notifycontroller

import (
	"context"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type Handler interface {
	Handler(ctx context.Context, input HandlerInput) *Response
}

type Response struct {
	Error error
}

type SqsEventHandler struct {
	InputPort notifyevent.InputPort
}

type HandlerInput struct {
	UserId  string
	Message string
}

func NewHandler() Handler {
	h := &SqsEventHandler{
		InputPort: di.Get("notifycontroller.InputPort").(notifyevent.InputPort),
	}
	return h
}

func (s *SqsEventHandler) Handler(ctx context.Context, input HandlerInput) *Response {
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
