package settime

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/log"
)

type UrlVerificationRequestHandler interface {
	Handle(ctx context.Context, data EventCallbackData) *Response
}

// UrlVerificationRequestHandlerFunctor represents url_verification event
// Ref. https://api.slack.com/events/url_verification
type UrlVerificationRequestHandlerFunctor struct {
}

type UrlVerificationResponseBody struct {
	Challenge string `json:"challenge"`
}

var _ UrlVerificationRequestHandler = (*UrlVerificationRequestHandlerFunctor)(nil)

func NewUrlVerificationRequestHandlerFunctor() *UrlVerificationRequestHandlerFunctor {
	return &UrlVerificationRequestHandlerFunctor{}
}

func (ur UrlVerificationRequestHandlerFunctor) Handle(ctx context.Context, data EventCallbackData) *Response {

	log.Info(fmt.Sprintf("UrlVerificationRequestHandler.Handler challenge=%s", data.Challenge))

	if data.Challenge == "" {
		return newErrorHandlerResponse("invalid challenge", "empty")
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	resp := &Response{
		StatusCode: http.StatusOK,
		Body: UrlVerificationResponseBody{
			data.Challenge,
		},
	}

	log.Info(fmt.Sprintf("UrlVerificationRequestHandler.Handler output=%v", *resp))

	return resp
}
