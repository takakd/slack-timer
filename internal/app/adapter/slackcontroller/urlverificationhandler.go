package slackcontroller

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/log"
)

// UrlVerificationRequestHandler represents url_verification event
// Ref. https://api.slack.com/events/url_verification
type UrlVerificationRequestHandler struct {
	Data *EventCallbackData
}

type UrlVerificationResponseBody struct {
	Challenge string `json:"challenge"`
}

// URL verification process just depends on Slack Event API, so no usecase and outputport.
func (ur *UrlVerificationRequestHandler) Handler(ctx context.Context) *HandlerResponse {

	log.Info(fmt.Sprintf("UrlVerificationRequestHandler.Handler challenge=%s", ur.Data.Challenge))

	if ur.Data.Challenge == "" {
		return makeErrorHandlerResponse("invalid challenge", "empty")
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	resp := &HandlerResponse{
		StatusCode: http.StatusOK,
		Body: UrlVerificationResponseBody{
			ur.Data.Challenge,
		},
	}

	log.Info(fmt.Sprintf("UrlVerificationRequestHandler.Handler output=%v", *resp))

	return resp
}
