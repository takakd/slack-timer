package slackcontroller

import (
	"context"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUrlVerificationHandler_Handler(t *testing.T) {
	cases := []struct {
		name      string
		challenge string
		resp      EventCallbackResponse
	}{
		{"empty challenge", "", EventCallbackResponse{
			Message:    "invalid challenge",
			StatusCode: http.StatusInternalServerError,
		}},
		{"ok", "valid token", EventCallbackResponse{
			Message:    "success",
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
