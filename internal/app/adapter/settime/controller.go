package settime

import (
	"encoding/json"
	"net/http"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Controller handles "set" and "urlverification" command.
type Controller struct {
}

// NewController creates new struct.
func NewController() *Controller {
	h := &Controller{}
	return h
}

var _ ControllerHandler = (*Controller)(nil)

// Handle calls handler according to "set" or "urlverification" command.
func (s Controller) Handle(ac appcontext.AppContext, input HandleInput) *Response {
	// URL verification event
	if input.EventData.isVerificationEvent() {
		log.InfoWithContext(ac, "URL verification event")
		rh := di.Get("settime.URLVerificationRequestHandler").(URLVerificationRequestHandler)
		return rh.Handle(ac, input.EventData)
	}

	var resp *Response

	if input.EventData.MessageEvent.isBotMessage() {
		log.InfoWithContext(ac, "ignore bot message")
		resp = &Response{
			StatusCode: http.StatusOK,
			Body:       "bot message",
		}
		return resp
	}

	if input.EventData.MessageEvent.isSetTimeEvent() {
		rh := di.Get("settime.SaveEventHandler").(SaveEventHandler)
		resp = rh.Handle(ac, input.EventData)

	} else if input.EventData.MessageEvent.isOnEvent() {
		rh := di.Get("settime.OnEventHandler").(OnEventHandler)
		resp = rh.Handle(ac, input.EventData)

	} else if input.EventData.MessageEvent.isOffEvent() {
		rh := di.Get("settime.OffEventHandler").(OffEventHandler)
		resp = rh.Handle(ac, input.EventData)

	} else {
		resp = newErrorHandlerResponse(ac, "invalid event", input.EventData)
	}
	return resp
}

// ResponseErrorBody is used if response status is error.
type ResponseErrorBody struct {
	Message string
	Detail  string
}

func newErrorHandlerResponse(ac appcontext.AppContext, message string, detail interface{}) *Response {
	body := &ResponseErrorBody{
		Message: message,
	}
	if detail != nil {
		if detailJSON, err := json.Marshal(detail); err != nil {
			log.ErrorWithContext(ac, "marshal error", err.Error())
		} else {
			body.Detail = string(detailJSON)
		}
	}
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}
