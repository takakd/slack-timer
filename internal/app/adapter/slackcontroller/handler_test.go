package slackcontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestMakeErrorHandleResponse(t *testing.T) {
	cases := []struct {
		name    string
		message string
		err     error
	}{
		{"no error", "test", nil},
		{"error", "test", errors.New("test err")},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			wantBody := &HandlerResponseErrorBody{
				Message: c.message,
			}
			if c.err != nil {
				wantBody.Detail = c.err.Error()
			}
			want := &HandlerResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       wantBody,
			}

			got := makeErrorHandlerResponse(c.message, c.err)
			assert.Equal(t, want, got)
		})
	}
}

func TestHandler(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		caseJson, err := json.Marshal(EventCallbackData{})
		require.NoError(t, err)

		caseInput := LambdaInput{
			Body: string(caseJson),
		}
		want := makeErrorHandlerResponse("parameter error", ErrInvalidParameters)
		got, err := LambdaHandleRequest(context.TODO(), caseInput)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ok", func(t *testing.T) {
		caseData := &EventCallbackData{
			Type:      "url_verification",
			Challenge: "test challenge",
		}
		caseJson, err := json.Marshal(caseData)
		require.NoError(t, err)

		caseInput := LambdaInput{
			Body: string(caseJson),
		}
		ctx := context.TODO()

		body, _ := json.Marshal(UrlVerificationResponseBody{
			caseData.Challenge,
		})

		want := LambdaOutput{
			StatusCode: http.StatusOK,
			Body:       string(body),
		}
		got, err := LambdaHandleRequest(ctx, caseInput)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
