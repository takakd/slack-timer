// Package slackcontroller provides the slack Event API callback handler.
// Ref: https://api.slack.com/events-api#the-events-api__receiving-events
package settime

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Command types entered by users.
const (
	CmdSet = "set"
)

type SetTimeController struct {
}

func NewSetTimeController() Controller {
	h := SetTimeController{}
	return h
}

func (c SetTimeController) Handle(ctx context.Context, input HandleInput) *Response {
	// Create request struct corresponding to input.

	// URL verification event
	if input.EventData.isVerificationEvent() {
		log.Info("url verification event")
		rh := di.Get("slackcontroller.urlverificationhandler").(UrlVerificationRequestHandler)
		return rh.Handle(ctx, input.EventData)
	}

	// Set interval minutes event
	if !input.EventData.MessageEvent.isSetEvent() {
		return makeErrorHandlerResponse("invalid event", fmt.Sprintf("type=%s", input.EventData.MessageEvent.Type))
	}

	rh := di.Get("slackcontroller.setcontroller").(SaveEventHandler)
	return rh.Handle(ctx, input.EventData)
}

type ResponseErrorBody struct {
	Message string
	Detail  string
}

func makeErrorHandlerResponse(message string, detail string) *Response {
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
