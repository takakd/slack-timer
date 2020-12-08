package settime

import (
	"context"
	"fmt"
	"net/http"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	assert.NotPanics(t, func() {
		NewController()
	})
}

func TestController_Handle(t *testing.T) {
	t.Run("ok:verification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				Type:      "url_verification",
				Challenge: "challenge",
			},
		}

		wantResp := &Response{}

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().Info("URL verification event")
		log.SetDefaultLogger(ml)

		mu := NewMockURLVerificationRequestHandler(ctrl)
		mu.EXPECT().Handle(gomock.Eq(ctx), gomock.Eq(caseInput.EventData)).Return(wantResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.URLVerificationRequestHandler")).Return(mu)
		di.SetDi(md)

		h := NewController()
		got := h.Handle(ctx, caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("not support event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				Type: "not support",
				MessageEvent: MessageEvent{
					Type: "not support",
				},
			},
		}

		wantResp := newErrorHandlerResponse("invalid event", fmt.Sprintf("type=%s", caseInput.EventData.Type))

		h := NewController()
		got := h.Handle(ctx, caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				MessageEvent: MessageEvent{
					Type: "message",
					Text: "set 10",
				},
			},
		}

		wantResp := &Response{
			StatusCode: http.StatusOK,
			Body:       "test",
		}

		mu := NewMockSaveEventHandler(ctrl)
		mu.EXPECT().Handle(gomock.Eq(ctx), gomock.Eq(caseInput.EventData)).Return(wantResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("settime.SaveEventHandler")).Return(mu)
		di.SetDi(md)

		h := NewController()
		got := h.Handle(ctx, caseInput)
		assert.Equal(t, wantResp, got)

	})
}

func TestMakeErrorHandleResponse(t *testing.T) {
	cases := []struct {
		name    string
		message string
		detail  string
	}{
		{"no error", "test", ""},
		{"error", "test", "test err"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			wantBody := &ResponseErrorBody{
				Message: c.message,
			}
			if c.detail != "" {
				wantBody.Detail = c.detail
			}
			want := &Response{
				StatusCode: http.StatusInternalServerError,
				Body:       wantBody,
			}

			got := newErrorHandlerResponse(c.message, c.detail)
			assert.Equal(t, want, got)
		})
	}
}
