package settime

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Controller handles "set" and "urlverification" command.
type Controller struct {
}

// NewController create new struct.
func NewController() *Controller {
	h := &Controller{}
	return h
}

var _ ControllerHandler = (*Controller)(nil)

// Handle calls handler according to "set" or "urlverification" command.
func (s Controller) Handle(ctx context.Context, input HandleInput) *Response {
	// URL verification event
	if input.EventData.isVerificationEvent() {
		log.Info("url verification event")
		rh := di.Get("settime.URLVerificationRequestHandler").(URLVerificationRequestHandler)
		return rh.Handle(ctx, input.EventData)
	}

	// Set interval minutes event
	if !input.EventData.MessageEvent.isSetTimeEvent() {
		return newErrorHandlerResponse("invalid event", fmt.Sprintf("type=%s", input.EventData.MessageEvent.Type))
	}

	rh := di.Get("settime.SaveEventHandler").(SaveEventHandler)
	return rh.Handle(ctx, input.EventData)
}

// ResponseErrorBody is used if response status is error.
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
