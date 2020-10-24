package slackcontroller

import (
	"context"
	"net/http"
)

// UrlVerificationRequestHandler represents url_verification event
// Ref. https://api.slack.com/events/url_verification
type UrlVerificationRequestHandler struct {
	Data *UrlVerificationEventCallbackData
}

// URL verification process just depends on Slack Event API, so no usecase and outputport.
func (ur *UrlVerificationRequestHandler) Handler(ctx context.Context, w http.ResponseWriter) {
	if ur.Data.Challenge == "" {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("invalid challenge"))
		return
	}

	// URL verification process just depends on Slack Event API, so no usecase and outputport.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ur.Data.Challenge))
	return
}
