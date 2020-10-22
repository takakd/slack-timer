package slackcontroller

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUrlVerificationHandler_Handler(t *testing.T) {
	cases := []struct {
		name      string
		challenge string
		status    int
		body      string
	}{
		{"empty challenge", "", http.StatusInternalServerError, "invalid challenge"},
		{"ok", "valid token", http.StatusOK, "valid token"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := UrlVerificationEventCallbackData{
				Challenge: c.challenge,
			}

			ctx := context.TODO()
			w := httptest.NewRecorder()

			h := UrlVerificationRequestHandler{
				Data: &data,
			}
			h.Handler(ctx, w)

			assert.Equal(t, w.Code, c.status)
			assert.Equal(t, w.Body.String(), c.body)
		})
	}
}
