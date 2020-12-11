package settime

import (
	"net/http"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/log"
)

// URLVerificationRequestHandler handles "urlverification" command.
type URLVerificationRequestHandler interface {
	Handle(ac appcontext.AppContext, data EventCallbackData) *Response
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

// NewURLVerificationRequestHandlerFunctor creates new struct.
func NewURLVerificationRequestHandlerFunctor() *URLVerificationRequestHandlerFunctor {
	return &URLVerificationRequestHandlerFunctor{}
}

// Handle response according to the Slack URL verification specification.
func (ur URLVerificationRequestHandlerFunctor) Handle(ac appcontext.AppContext, data EventCallbackData) *Response {

	log.InfoWithContext(ac, "URLVerification called", data.Challenge)

	if data.Challenge == "" {
		return newErrorHandlerResponse(ac, "invalid challenge", "empty")
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	resp := &Response{
		StatusCode: http.StatusOK,
		Body: URLVerificationResponseBody{
			data.Challenge,
		},
	}

	log.InfoWithContext(ac, "URLVerification outputs", *resp)

	return resp
}
