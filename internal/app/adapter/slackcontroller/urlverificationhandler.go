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
func (ur *UrlVerificationRequestHandler) Handler(ctx context.Context) EventCallbackResponse {
	if ur.Data.Challenge == "" {
		return *makeErrorCallbackResponse("invalid challenge", nil)
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	return EventCallbackResponse{
		Message:    "success",
		StatusCode: http.StatusOK,
	}
}
