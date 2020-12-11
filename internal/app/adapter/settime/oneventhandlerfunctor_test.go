package settime

import (
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"slacktimer/internal/app/usecase/timeronevent"
	"slacktimer/internal/app/util/log"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewOnEventHandlerFunctor(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mi := timeronevent.NewMockInputPort(ctrl)
		mp := NewOnEventOutputReceivePresenter()
		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("timeronevent.InputPort").Return(mi)
		md.EXPECT().Get("settime.OnEventOutputReceivePresenter").Return(mp)
		di.SetDi(md)

		assert.NotPanics(t, func() {
			NewOnEventHandlerFunctor()
		})
	})

	t.Run("ng", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("timeronevent.InputPort").Return(nil)
		di.SetDi(md)

		assert.Panics(t, func() {
			NewOnEventHandlerFunctor()
		})
	})
}

func TestOnEventHandlerFunctor_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	caseData := EventCallbackData{
		MessageEvent: MessageEvent{
			Type: "message",
			User: "test",
			Text: "on",
		},
	}

	ac := appcontext.TODO()
	caseInput := timeronevent.InputData{
		UserID: caseData.MessageEvent.User,
	}

	ml := log.NewMockLogger(ctrl)
	ml.EXPECT().InfoWithContext(ac, "timer on", "call inputport", map[string]interface {
	}{
		"user": caseData.MessageEvent.User,
		"text": caseData.MessageEvent.Text,
	})
	ml.EXPECT().InfoWithContext(ac, "return from inputport")
	log.SetDefaultLogger(ml)

	mp := NewOnEventOutputReceivePresenter()

	mu := timeronevent.NewMockInputPort(ctrl)
	mu.EXPECT().SetEventOn(ac, caseInput, mp)

	md := di.NewMockDI(ctrl)
	md.EXPECT().Get("timeronevent.InputPort").Return(mu)
	md.EXPECT().Get("settime.OnEventOutputReceivePresenter").Return(mp)
	di.SetDi(md)

	h := NewOnEventHandlerFunctor()
	h.Handle(ac, caseData)
}
