package slackcontroller

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/log"
)

type UrlVerificationRequestHandler interface {
	Handler(ctx context.Context, data EventCallbackData) *Response
}

// UrlVerificationRequestHandler represents url_verification event
// Ref. https://api.slack.com/events/url_verification
type UrlVerificationController struct {
}

type UrlVerificationResponseBody struct {
	Challenge string `json:"challenge"`
}

func NewUrlVerificationController() UrlVerificationRequestHandler {
	return &UrlVerificationController{}
}

// URL verification process just depends on Slack Event API, so no usecase and outputport.
func (ur *UrlVerificationController) Handler(ctx context.Context, data EventCallbackData) *Response {

	log.Info(fmt.Sprintf("UrlVerificationRequestHandler.Handler challenge=%s", data.Challenge))

	if data.Challenge == "" {
		return makeErrorHandlerResponse("invalid challenge", "empty")
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
