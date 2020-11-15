package slackcontroller

import (
	"context"
	"net/http"
)

// UrlVerificationRequestHandler represents url_verification event
// Ref. https://api.slack.com/events/url_verification
type UrlVerificationRequestHandler struct {
	Data *EventCallbackData
}

// URL verification process just depends on Slack Event API, so no usecase and outputport.
func (ur *UrlVerificationRequestHandler) Handler(ctx context.Context) *HandlerResponse {
	if ur.Data.Challenge == "" {
		return makeErrorHandlerResponse("invalid challenge", nil)
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	return &HandlerResponse{
		StatusCode: http.StatusOK,
		Body:       ur.Data.Challenge,
	}
}
