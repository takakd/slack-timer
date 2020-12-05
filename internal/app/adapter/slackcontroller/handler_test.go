package slackcontroller

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"slacktimer/internal/app/util/di"
	"strconv"
	"testing"
)

func TestEventCallbackData_isVerificationEvent(t *testing.T) {
	cases := []struct {
		name         string
		dataType     string
		verification bool
	}{
		{"ok", "url_verification", true},
		{"ng", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &EventCallbackData{
				Type: c.dataType,
			}
			assert.Equal(t, d.isVerificationEvent(), c.verification)
		})
	}
}

func TestMessageEvent_isSetEvent(t *testing.T) {
	cases := []struct {
		name      string
		eventType string
		text      string
		isSet     bool
	}{
		{"ok", "message", "set 10", true},
		{"ng:wrong type", "wrong", "set 10", false},
		{"ng:wrong body", "messazge", "set a", false},
		{"ng:empty body", "messazge", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				Type: c.eventType,
				Text: c.text,
			}
			got := d.isSetEvent()
			assert.Equal(t, c.isSet, got)
		})
	}
}

func TestMessageEvent_eventUnixTimeStamp(t *testing.T) {
	cases := []struct {
		name    string
		ts      string
		tsNano  string
		success bool
	}{
		{"ok", "1607165903", "000010", true},
		{"ok:empty nano", "1607165903", "", true},
		{"ng:invalid format", "abc", "", false},
		{"ng:empty", "", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := &MessageEvent{
				EventTs: fmt.Sprintf("%s.%s", c.ts, c.tsNano),
			}
			got, err := d.eventUnixTimeStamp()

			if c.success {
				assert.NoError(t, err)
				want, err := strconv.ParseInt(c.ts, 10, 64)
				require.NoError(t, err)
				assert.Equal(t, want, got)

			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	assert.NotPanics(t, func() {
		NewHandler()
	})
}

func TestSlackEventHandler_Handler(t *testing.T) {
	t.Run("ok:verification", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := HandlerInput{
			EventData: EventCallbackData{
				Type:      "url_verification",
				Challenge: "challenge",
			},
		}

		wantResp := &Response{}

		mu := NewMockUrlVerificationRequestHandler(ctrl)
		mu.EXPECT().Handler(gomock.Eq(ctx), gomock.Eq(caseInput.EventData)).Return(wantResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("slackcontroller.urlverificationhandler")).Return(mu)
		di.SetDi(md)

		h := NewHandler()
		got := h.Handler(ctx, caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("not support event", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := HandlerInput{
			EventData: EventCallbackData{
				Type: "not support",
				MessageEvent: MessageEvent{
					Type: "not support",
				},
			},
		}

		wantResp := makeErrorHandlerResponse("invalid event", fmt.Sprintf("type=%s", caseInput.EventData.Type))

		h := NewHandler()
		got := h.Handler(ctx, caseInput)
		assert.Equal(t, wantResp, got)
	})

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := HandlerInput{
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

		mu := NewMockSetRequestHandler(ctrl)
		mu.EXPECT().Handler(gomock.Eq(ctx), gomock.Eq(caseInput.EventData)).Return(wantResp)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get(gomock.Eq("slackcontroller.setcontroller")).Return(mu)
		di.SetDi(md)

		h := NewHandler()
		got := h.Handler(ctx, caseInput)
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

			got := makeErrorHandlerResponse(c.message, c.detail)
			assert.Equal(t, want, got)
		})
	}
}
