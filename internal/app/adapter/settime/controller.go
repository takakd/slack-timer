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

// NewController create new struct.
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

	// Set interval minutes event
	if !input.EventData.MessageEvent.isSetTimeEvent() {
		return newErrorHandlerResponse(ac, "invalid event", input.EventData)
	}

	rh := di.Get("settime.SaveEventHandler").(SaveEventHandler)
	return rh.Handle(ac, input.EventData)
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
			log.ErrorWithContext(ac, "marshal error", err)
		} else {
			body.Detail = string(detailJSON)
		}
	}
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
	}
}
