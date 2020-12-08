package settime

import (
	"context"
	"net/http"
	"slacktimer/internal/app/util/log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewURLVerificationRequestHandlerFunctor(t *testing.T) {
	assert.NotPanics(t, func() {
		NewURLVerificationRequestHandlerFunctor()
	})
}

func TestURLVerificationRequestHandlerFunctor_Handle(t *testing.T) {
	cases := []struct {
		name      string
		challenge string
		resp      *Response
	}{
		{"empty challenge", "", &Response{
			Body: &ResponseErrorBody{
				Message: "invalid challenge",
				Detail:  "empty",
			},
			StatusCode: http.StatusInternalServerError,
		}},
		{"ok", "valid token", &Response{
			StatusCode: http.StatusOK,
			Body: URLVerificationResponseBody{
				"valid token",
			},
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caseData := EventCallbackData{
				Challenge: c.challenge,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ml := log.NewMockLogger(ctrl)
			ml.EXPECT().Info("URLVerification called", caseData.Challenge)
			if caseData.Challenge != "" {
				ml.EXPECT().Info("URLVerification outputs", *c.resp)
			}
			log.SetDefaultLogger(ml)

			h := NewURLVerificationRequestHandlerFunctor()
			got := h.Handle(context.TODO(), caseData)
			assert.Equal(t, c.resp, got)
		})
	}
}
