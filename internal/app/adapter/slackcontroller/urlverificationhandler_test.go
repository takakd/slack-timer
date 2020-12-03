package slackcontroller

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"slacktimer/internal/app/util/log"
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
				Detail:  "empty",
			},
			StatusCode: http.StatusInternalServerError,
		}},
		{"ok", "valid token", &HandlerResponse{
			StatusCode: http.StatusOK,
			Body: UrlVerificationResponseBody{
				"valid token",
			},
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caseData := &EventCallbackData{
				Challenge: c.challenge,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			l := log.NewMockLogger(ctrl)
			l.EXPECT().Info(gomock.Eq(fmt.Sprintf("UrlVerificationRequestHandler.Handler challenge=%s", caseData.Challenge)))
			if caseData.Challenge != "" {
				l.EXPECT().Info(gomock.Eq(fmt.Sprintf("UrlVerificationRequestHandler.Handler output=%v", c.resp)))
			}
			log.SetDefaultLogger(l)

			h := UrlVerificationRequestHandler{
				Data: caseData,
			}
			got := h.Handler(context.TODO())
			assert.Equal(t, c.resp, got)
		})
	}
}
