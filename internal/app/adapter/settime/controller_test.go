package settime

import (
	"net/http"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	assert.NotPanics(t, func() {
		NewController()
	})
}

func TestController_Handle(t *testing.T) {
	t.Run("ok:bot", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				MessageEvent: MessageEvent{
					Type:        "message",
					ChannelType: "im",
					Text:        "OK, start the notification.",
					BotID:       "ABC1234567",
				},
			},
		}

		wantResp := &Response{
			StatusCode: http.StatusOK,
			Body:       "bot message",
		}

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().InfoWithContext(ac, "ignore bot message")
		log.SetDefaultLogger(ml)

		h := NewController()
		got := h.Handle(ac, caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("ok:verification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				Type:      "url_verification",
				Challenge: "challenge",
			},
		}

		wantResp := &Response{}

		ml := log.NewMockLogger(ctrl)
		ml.EXPECT().InfoWithContext(ac, "URL verification event")
		log.SetDefaultLogger(ml)

		mu := NewMockURLVerificationRequestHandler(ctrl)
		mu.EXPECT().Handle(ac, caseInput.EventData).Return(wantResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("settime.URLVerificationRequestHandler").Return(mu)
		di.SetDi(md)

		h := NewController()
		got := h.Handle(ac, caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("not support event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				Type:         "not support",
				MessageEvent: MessageEvent{},
			},
		}

		wantResp := newErrorHandlerResponse(appcontext.TODO(), "invalid event", caseInput.EventData)

		h := NewController()
		got := h.Handle(appcontext.TODO(), caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("ok:set time", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseInput := HandleInput{
			EventData: EventCallbackData{
				MessageEvent: MessageEvent{
					Type:        "message",
					ChannelType: "im",
					Text:        "set 10 message",
				},
			},
		}

		wantResp := &Response{
			StatusCode: http.StatusOK,
			Body:       "test",
		}

		mu := NewMockSaveEventHandler(ctrl)
		mu.EXPECT().Handle(ac, caseInput.EventData).Return(wantResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("settime.SaveEventHandler").Return(mu)
		di.SetDi(md)

		h := NewController()
		got := h.Handle(ac, caseInput)
		assert.Equal(t, wantResp, got)

	})

	t.Run("ng:invalid command", func(t *testing.T) {
		cases := []struct {
			name string
			cmd  string
		}{
			{name: "set time", cmd: "set 10"},
			{name: "on", cmd: "onn"},
			{name: "off", cmd: "offf"},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				ac := appcontext.TODO()

				caseInput := HandleInput{
					EventData: EventCallbackData{
						MessageEvent: MessageEvent{
							Type:        "message",
							ChannelType: "im",
							Text:        c.cmd,
						},
					},
				}

				wantResp := newErrorHandlerResponse(ac, "invalid event", caseInput.EventData)

				h := NewController()
				got := h.Handle(ac, caseInput)
				assert.Equal(t, wantResp, got)
			})
		}
	})
}

func TestNewErrorHandleResponse(t *testing.T) {
	cases := []struct {
		name    string
		message string
		detail  interface{}
		want    string
	}{
		{"no detail", "test", nil, ""},
		{"error", "test", "test err", `"test err"`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			wantBody := &ResponseErrorBody{
				Message: c.message,
			}
			if c.detail != nil {
				wantBody.Detail = c.want
			}
			want := &Response{
				StatusCode: http.StatusInternalServerError,
				Body:       wantBody,
			}

			got := newErrorHandlerResponse(appcontext.TODO(), c.message, c.detail)
			assert.Equal(t, want, got)
		})
	}
}
