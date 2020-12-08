package settime

import (
	"context"
	"net/http"
	"slacktimer/internal/app/util/log"
)

// URLVerificationRequestHandler handles "urlverification" command.
type URLVerificationRequestHandler interface {
	Handle(ctx context.Context, data EventCallbackData) *Response
}

// URLVerificationRequestHandlerFunctor represents url_verification event.
// Ref. https://api.slack.com/events/url_verification
type URLVerificationRequestHandlerFunctor struct {
}

// URLVerificationResponseBody represents the url verification event payload in EventCallbackData.
type URLVerificationResponseBody struct {
	Challenge string `json:"challenge"`
}

var _ URLVerificationRequestHandler = (*URLVerificationRequestHandlerFunctor)(nil)

// NewURLVerificationRequestHandlerFunctor create new struct.
func NewURLVerificationRequestHandlerFunctor() *URLVerificationRequestHandlerFunctor {
	return &URLVerificationRequestHandlerFunctor{}
}

// Handle response according to the Slack URL verification specification.
func (ur URLVerificationRequestHandlerFunctor) Handle(ctx context.Context, data EventCallbackData) *Response {

	log.Info("URLVerification called", data.Challenge)

	if data.Challenge == "" {
		return newErrorHandlerResponse("invalid challenge", "empty")
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	resp := &Response{
		StatusCode: http.StatusOK,
		Body: URLVerificationResponseBody{
			data.Challenge,
		},
	}

	log.Info("URLVerification outputs", *resp)

	return resp
}
