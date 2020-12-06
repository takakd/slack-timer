package settime

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

type SetTimeController struct {
}

func NewSetTimeController() Controller {
	h := SetTimeController{}
	return h
}

func (c SetTimeController) Handle(ctx context.Context, input HandleInput) *Response {
	// URL verification event
	if input.EventData.isVerificationEvent() {
		log.Info("url verification event")
		rh := di.Get("settime.UrlVerificationRequestHandler").(UrlVerificationRequestHandler)
		return rh.Handle(ctx, input.EventData)
	}

	// Set interval minutes event
	if !input.EventData.MessageEvent.isSetTimeEvent() {
		return newErrorHandlerResponse("invalid event", fmt.Sprintf("type=%s", input.EventData.MessageEvent.Type))
	}

	rh := di.Get("settime.SaveEventHandler").(SaveEventHandler)
	return rh.Handle(ctx, input.EventData)
}

type ResponseErrorBody struct {
	Message string
	Detail  string
}

func newErrorHandlerResponse(message string, detail string) *Response {
	body := &ResponseErrorBody{
		Message: message,
	}
	if detail != "" {
		body.Detail = detail
	}
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}
