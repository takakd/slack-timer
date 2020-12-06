package settime

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"slacktimer/internal/app/util/log"
	"testing"
)

func TestNewUrlVerificationRequestHandlerFunctor(t *testing.T) {
	assert.NotPanics(t, func() {
		NewUrlVerificationRequestHandlerFunctor()
	})
}

func TestUrlVerificationRequestHandlerFunctor_Handle(t *testing.T) {
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
			Body: UrlVerificationResponseBody{
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
			ml.EXPECT().Info(gomock.Eq(fmt.Sprintf("UrlVerificationRequestHandler.Handler challenge=%s", caseData.Challenge)))
			if caseData.Challenge != "" {
				ml.EXPECT().Info(gomock.Eq(fmt.Sprintf("UrlVerificationRequestHandler.Handler output=%v", *c.resp)))
			}
			log.SetDefaultLogger(ml)

			h := NewUrlVerificationRequestHandlerFunctor()
			got := h.Handle(context.TODO(), caseData)
			assert.Equal(t, c.resp, got)
		})
	}
}
