package slackcontroller

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/usecase/updatetimerevent"
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

func TestNewRequestHandler(t *testing.T) {
	t.Run("ok:verify", func(t *testing.T) {
		caseData := &EventCallbackData{
			Type: "url_verification",
		}

		h, err := NewRequestHandler(caseData)
		assert.NoError(t, err)
		_, ok := h.(*UrlVerificationRequestHandler)
		assert.True(t, ok)
	})

	t.Run("not support event", func(t *testing.T) {
		caseData := &EventCallbackData{
			Type: "",
			MessageEvent: MessageEvent{
				Type: "test",
			},
		}
		caseErr := fmt.Errorf("invalid event type, type=%s", caseData.MessageEvent.Type)

		h, err := NewRequestHandler(caseData)
		assert.Nil(t, h)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseData := &EventCallbackData{
			MessageEvent: MessageEvent{
				Type: "message",
				Text: "set 10",
			},
		}
		caseUsecase := &updatetimerevent.Interactor{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := di.NewMockDI(ctrl)
		m.EXPECT().Get(gomock.Eq("UpdateTimerEvent")).Return(caseUsecase)
		di.SetDi(m)

		h, err := NewRequestHandler(caseData)
		assert.NoError(t, err)
		assert.Equal(t, &SetRequestHandler{
			messageEvent: &caseData.MessageEvent,
			usecase:      caseUsecase,
		}, h)
	})

	cases := []struct {
		name      string
		eventType string
		text      string
		err       error
	}{
		{"invalid format", "", "test", fmt.Errorf("invalid event type, type=%s", "")},
		{"invalid type", "message", "invalid 1", fmt.Errorf("invalid sub type, subtype=%s", "invalid")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caseData := &EventCallbackData{
				MessageEvent: MessageEvent{
					Type: c.eventType,
					Text: c.text,
				},
			}

			h, err := NewRequestHandler(caseData)
			assert.Nil(t, h)
			assert.Equal(t, c.err, err)
		})
	}
}

func TestMakeErrorCallbackResponseBody(t *testing.T) {
	cases := []struct {
		name    string
		message string
		err     string
	}{
		{"no error", "test", ""},
		{"error", "test", "test err"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			want := &EventCallbackResponse{
				Message:    c.message,
				StatusCode: http.StatusInternalServerError,
			}
			if c.err != "" {
				want.Detail = c.err
			}
			got := makeErrorCallbackResponse(c.message, errors.New(c.err))
			assert.Equal(t, want, got)
		})
	}
}

func TestHandler(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		caseData := EventCallbackData{
		}
		want := *makeErrorCallbackResponse("parameter error", ErrInvalidRequest)
		got, err := LambdaHandleRequest(context.TODO(), caseData)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ok", func(t *testing.T) {
		caseData := EventCallbackData{
			Type: "url_verification",
		}
		ctx := context.TODO()
		h := UrlVerificationRequestHandler{
			Data: &caseData,
		}

		want := h.Handler(ctx)
		got, err := LambdaHandleRequest(ctx, caseData)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
