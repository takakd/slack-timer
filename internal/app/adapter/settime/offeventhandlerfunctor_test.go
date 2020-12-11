package settime

import (
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"slacktimer/internal/app/usecase/timeroffevent"
	"slacktimer/internal/app/util/log"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewOffEventHandlerFunctor(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mi := timeroffevent.NewMockInputPort(ctrl)
		mp := NewOffEventOutputReceivePresenter()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("timeroffevent.InputPort").Return(mi)
		md.EXPECT().Get("settime.OffEventOutputReceivePresenter").Return(mp)
		di.SetDi(md)

		assert.NotPanics(t, func() {
			NewOffEventHandlerFunctor()
		})
	})

	t.Run("ng", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("timeroffevent.InputPort").Return(nil)
		di.SetDi(md)

		assert.Panics(t, func() {
			NewOffEventHandlerFunctor()
		})
	})
}

func TestOffEventHandlerFunctor_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	caseData := EventCallbackData{
		MessageEvent: MessageEvent{
			Type: "message",
			User: "test",
			Text: "off",
		},
	}

	ac := appcontext.TODO()
	caseInput := timeroffevent.InputData{
		UserID: caseData.MessageEvent.User,
	}

	mp := NewOffEventOutputReceivePresenter()

	ml := log.NewMockLogger(ctrl)
	ml.EXPECT().InfoWithContext(ac, "timer on", "call inputport", map[string]interface {
	}{
		"user": caseData.MessageEvent.User,
		"text": caseData.MessageEvent.Text,
	})
	ml.EXPECT().InfoWithContext(ac, "return from inputport", mp.Resp)
	log.SetDefaultLogger(ml)

	mu := timeroffevent.NewMockInputPort(ctrl)
	mu.EXPECT().SetEventOff(ac, caseInput, mp)

	md := di.NewMockDI(ctrl)
	md.EXPECT().Get("timeroffevent.InputPort").Return(mu)
	md.EXPECT().Get("settime.OffEventOutputReceivePresenter").Return(mp)
	di.SetDi(md)

	h := NewOffEventHandlerFunctor()
	h.Handle(ac, caseData)
}
