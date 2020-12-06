package settime

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestSaveEventHandlerFunctor_validateTs(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		ts    string
		min   int
		valid bool
	}{
		{"ok", "set 10", "1606830655.000003", 10, true},
		{"ng", "set 10", "1606830655", 10, true},
		{"ng", "set 10", "", 10, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mi := updatetimerevent.NewMockInputPort(ctrl)
			md := di.NewMockDI(ctrl)
			md.EXPECT().Get(gomock.Eq("updatetimerevent.InputPort")).Return(mi)
			di.SetDi(md)

			caseData := EventCallbackData{
				MessageEvent: MessageEvent{
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}

			ct := NewSaveEventHandlerFunctor()

			bag := ct.validate(caseData)
			_, exists := bag.GetError("timestamp")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, ct.remindIntervalInMin)
			}
		})
	}
}

func TestSaveEventHandlerFunctor_validateSet(t *testing.T) {
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
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mi := updatetimerevent.NewMockInputPort(ctrl)
			md := di.NewMockDI(ctrl)
			md.EXPECT().Get(gomock.Eq("updatetimerevent.InputPort")).Return(mi)
			di.SetDi(md)

			caseData := EventCallbackData{
				MessageEvent: MessageEvent{
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}

			ct := NewSaveEventHandlerFunctor()

			bag := ct.validate(caseData)
			_, exists := bag.GetError("interval")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, ct.remindIntervalInMin)
			}
		})
	}
}

func TestSaveEventHandlerFunctor_Handle(t *testing.T) {
	cases := []struct {
		name string
		text string
		ts   string
		resp Response
	}{
		{"timestamp validate error", "set 10", "", Response{
			StatusCode: http.StatusInternalServerError,
			Body: &ResponseErrorBody{
				Message: "invalid parameter",
				Detail:  "invalid format",
			},
		}},
		{"set command validate error", "", "1606830655", Response{
			StatusCode: http.StatusInternalServerError,
			Body: &ResponseErrorBody{
				Message: "invalid parameter",
				Detail:  "invalid format",
			},
		}},
		{"ok", "set 10", "1606830655", Response{
			StatusCode: http.StatusOK,
			Body:       "success",
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			caseData := EventCallbackData{
				MessageEvent: MessageEvent{
					Type:    "message",
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}

			ctx := context.TODO()

			mu := updatetimerevent.NewMockInputPort(ctrl)
			if c.text != "" && c.ts != "" {
				mu.EXPECT().SaveIntervalMin(gomock.Eq(ctx), gomock.Eq(caseData.MessageEvent.User), gomock.Any(), gomock.Eq(10), gomock.Any()).DoAndReturn(func(_, _, _, _, outputPort interface{}) {

					output := outputPort.(*SaveEventOutputReceivePresenter)
					output.Resp = c.resp
				})
			}

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("updatetimerevent.InputPort").Return(mu)
			di.SetDi(md)

			h := NewSaveEventHandlerFunctor()
			got := h.Handle(ctx, caseData)
			assert.Equal(t, &c.resp, got)
		})
	}
}
