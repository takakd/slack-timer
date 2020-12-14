package settime

import (
	"net/http"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/log"
	"time"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSaveEventHandlerFunctor_parseTs(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		ts    string
		min   int
		valid bool
	}{
		{"ok", "set 10 Hi!", "1606830655.000003", 10, true},
		{"ng", "set 10 Hi!", "1606830655", 10, true},
		{"ng", "set 10 Hi!", "", 10, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mi := updatetimerevent.NewMockInputPort(ctrl)
			mp := NewSaveEventOutputReceivePresenter()

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("updatetimerevent.InputPort").Return(mi)
			md.EXPECT().Get("settime.SaveEventOutputReceivePresenter").Return(mp)
			di.SetDi(md)

			caseData := EventCallbackData{
				MessageEvent: MessageEvent{
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}

			ct := NewSaveEventHandlerFunctor()

			bag := ct.parse(caseData)
			_, exists := bag.GetError("timestamp")
			assert.Equal(t, c.valid, !exists)
			if c.valid {
				assert.Equal(t, c.min, ct.remindIntervalInMin)
			}
		})
	}
}

func TestSaveEventHandlerFunctor_parseSet(t *testing.T) {
	cases := []struct {
		name  string
		text  string
		ts    string
		min   int
		valid bool
	}{
		{"ok", "set 10 Hi!", "1606830655", 10, true},
		{"ok", "set 1 Hi!", "1606830655", 1, true},
		{"ng", "set -1 Hi!", "1606830655", 0, false},
		{"ng", "set 10", "1606830655", 0, false},
		{"ng", "set", "1606830655", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mi := updatetimerevent.NewMockInputPort(ctrl)
			mp := NewSaveEventOutputReceivePresenter()

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("updatetimerevent.InputPort").Return(mi)
			md.EXPECT().Get("settime.SaveEventOutputReceivePresenter").Return(mp)
			di.SetDi(md)

			caseData := EventCallbackData{
				MessageEvent: MessageEvent{
					User:    "test",
					Text:    c.text,
					EventTs: c.ts,
				},
			}

			ct := NewSaveEventHandlerFunctor()

			bag := ct.parse(caseData)
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
		name    string
		text    string
		message string
		ts      string
		resp    Response
	}{
		{"timestamp validate error", "set 10 Hi!", "Hi!", "", Response{
			StatusCode: http.StatusInternalServerError,
			Body: &ResponseErrorBody{
				Message: "invalid parameter",
				Detail:  `"invalid format"`,
			},
		}},
		{"set command validate error", "", "", "1606830655", Response{
			StatusCode: http.StatusInternalServerError,
			Body: &ResponseErrorBody{
				Message: "invalid parameter",
				Detail:  `"invalid format"`,
			},
		}},
		{"ok", "set 10 Hi!", "Hi!", "1606830655.000010", Response{
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
					Ts:      c.ts,
				},
			}

			ac := appcontext.TODO()

			mu := updatetimerevent.NewMockInputPort(ctrl)
			mp := NewSaveEventOutputReceivePresenter()

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("updatetimerevent.InputPort").Return(mu)
			md.EXPECT().Get("settime.SaveEventOutputReceivePresenter").Return(mp)
			di.SetDi(md)

			h := NewSaveEventHandlerFunctor()

			if c.text != "" && c.ts != "" {
				h.parse(caseData)

				input := updatetimerevent.SaveEventInputData{
					UserID:      caseData.MessageEvent.User,
					CurrentTime: h.notificationTime,
					Minutes:     h.remindIntervalInMin,
					Text:        c.message,
				}
				mu.EXPECT().SaveIntervalMin(ac, input, gomock.Any()).DoAndReturn(func(_, _, outputPort interface{}) {
					output := outputPort.(*SaveEventOutputReceivePresenter)
					output.Resp = c.resp
				})
			}

			if c.resp.StatusCode == http.StatusOK {
				ml := log.NewMockLogger(ctrl)

				ts, _ := caseData.MessageEvent.eventUnixTimeStamp()
				ml.EXPECT().InfoWithContext(ac, "call inputport", map[string]interface{}{
					"user":              caseData.MessageEvent.User,
					"interval":          10,
					"text":              c.message,
					"notification time": time.Unix(ts, 0).UTC(),
					"ts":                caseData.MessageEvent.Ts,
				})

				ml.EXPECT().InfoWithContext(ac, "return from inputport", c.resp)
				log.SetDefaultLogger(ml)
			}

			got := h.Handle(ac, caseData)
			assert.Equal(t, &c.resp, got)
		})
	}
}
