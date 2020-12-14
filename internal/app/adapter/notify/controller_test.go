package notify

import (
	"errors"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mi := notifyevent.NewMockInputPort(ctrl)
	md := di.NewMockDI(ctrl)
	md.EXPECT().Get("notifyevent.InputPort").Return(mi)

	di.SetDi(md)

	h := NewController()
	assert.Equal(t, mi, h.inputPort)
}

func TestController_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ac := appcontext.TODO()
	caseInput := HandleInput{
		UserID:  "test user",
		Message: "test message",
	}
	caseError := errors.New("test error")

	ml := log.NewMockLogger(ctrl)
	ml.EXPECT().InfoWithContext(ac, "call inputport", gomock.Any())
	ml.EXPECT().InfoWithContext(ac, "return from inputport", gomock.Any())
	ml.EXPECT().InfoWithContext(ac, "handler output", gomock.Any())
	log.SetDefaultLogger(ml)

	mi := notifyevent.NewMockInputPort(ctrl)
	mi.EXPECT().NotifyEvent(ac, notifyevent.InputData{
		UserID:  caseInput.UserID,
		Message: caseInput.Message,
	}).Return(caseError)

	md := di.NewMockDI(ctrl)
	md.EXPECT().Get("notifyevent.InputPort").Return(mi)

	di.SetDi(md)

	h := NewController()
	resp := h.Handle(ac, caseInput)
	assert.Equal(t, &Response{
		Error: caseError,
	}, resp)
}
