package slackcontroller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUrlVerificationHandler_Handler(t *testing.T) {
	cases := []struct {
		name      string
		challenge string
		resp      *HandlerResponse
	}{
		{"empty challenge", "", &HandlerResponse{
			Body: &HandlerResponseErrorBody{
				Message: "invalid challenge",
			},
			StatusCode: http.StatusInternalServerError,
		}},
		{"ok", "valid token", &HandlerResponse{
			Body:       "valid token",
			StatusCode: http.StatusOK,
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caseData := &EventCallbackData{
				Challenge: c.challenge,
			}

			h := UrlVerificationRequestHandler{
				Data: caseData,
			}

			got := h.Handler(context.TODO())
			assert.Equal(t, c.resp, got)
		})
	}
}
