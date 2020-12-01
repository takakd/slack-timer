package slackcontroller

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"testing"
)

func TestSetRequestHandler_validateTs(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		ts    string
		min   int
		valid bool
	}{
		{"ok", "set 10", "1606830655", 10, true},
		{"ng", "set 10", "", 10, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := SetRequestHandler{
				messageEvent: &MessageEvent{
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}
			bag := r.validate()
			_, exists := bag.GetError("timestamp")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, r.remindIntervalInMin)
			}
		})
	}
}

func TestSetRequestHandler_validateSet(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		ts    string
		min   int
		valid bool
	}{
		{"ok", "set 10", "1606830655", 10, true},
		{"ok", "set 1", "1606830655", 1, true},
		{"ng", "set -1", "1606830655", 0, false},
		{"ng", "set", "1606830655", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := SetRequestHandler{
				messageEvent: &MessageEvent{
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}
			bag := r.validate()
			_, exists := bag.GetError("interval")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, r.remindIntervalInMin)
			}
		})
	}
}

func TestSetRequestHandler_Handler(t *testing.T) {
	cases := []struct {
		name string
		text string
		ts   string
		resp *HandlerResponse
	}{
		{"timestamp validate error", "set 10", "", &HandlerResponse{
			StatusCode: http.StatusInternalServerError,
			Body: &HandlerResponseErrorBody{
				Message: "invalid format",
				Detail:  "invalid parameters",
			},
		}},
		{"set command validate error", "", "1606830655", &HandlerResponse{
			StatusCode: http.StatusInternalServerError,
			Body: &HandlerResponseErrorBody{
				Message: "invalid format",
				Detail:  "invalid parameters",
			},
		}},
		{"ok", "set 10", "1606830655", &HandlerResponse{
			StatusCode: http.StatusOK,
			Body:       "success",
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caseData := &EventCallbackData{
				MessageEvent: MessageEvent{
					Type:    "message",
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}

			ctx := context.TODO()

			var mu *updatetimerevent.MockUsecase
			if c.text != "" {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mu = updatetimerevent.NewMockUsecase(ctrl)

				if c.ts != "" {
					mu.EXPECT().SaveIntervalMin(gomock.Eq(ctx), gomock.Eq(caseData.MessageEvent.User), gomock.Any(), gomock.Eq(10), gomock.Any()).DoAndReturn(func(_, _, _, _, outputPort interface{}) {
						output := outputPort.(*SetRequestOutputPort)
						output.Resp = c.resp
					})
				}
			}

			h := SetRequestHandler{
				messageEvent: &caseData.MessageEvent,
				usecase:      mu,
			}
			got := h.Handler(ctx)
			assert.Equal(t, c.resp, got)
		})
	}
}

func TestSetRequestOutputPort_Output(t *testing.T) {
	cases := []struct {
		name string
		err  error
		msg  string
	}{
		{name: "ng:find", err: updatetimerevent.ErrFind, msg: "failed to find event"},
		{name: "ng:create", err: updatetimerevent.ErrCreate, msg: "failed to create event"},
		{name: "ng:save", err: updatetimerevent.ErrSave, msg: "failed to save event"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			caseData := &updatetimerevent.OutputData{
				Result: c.err,
			}
			wantResp := makeErrorHandlerResponse("failed to save event", ErrSaveEvent)

			outputPort := &SetRequestOutputPort{}
			outputPort.Output(caseData)

			assert.Equal(t, wantResp, outputPort.Resp)
		})
	}

	t.Run("ok", func(t *testing.T) {
		caseData := &updatetimerevent.OutputData{
			Result: nil,
		}
		wantResp := &HandlerResponse{
			StatusCode: http.StatusOK,
			Body:       "success",
		}

		outputPort := &SetRequestOutputPort{}
		outputPort.Output(caseData)

		assert.Equal(t, wantResp, outputPort.Resp)
	})
}
