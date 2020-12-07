package notify

import (
	"context"
	"errors"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i := notifyevent.NewMockInputPort(ctrl)
	d := di.NewMockDI(ctrl)
	d.EXPECT().Get(gomock.Eq("notifyevent.InputPort")).Return(i)

	di.SetDi(d)

	h := NewController()
	assert.Equal(t, i, h.InputPort)
}

func TestController_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	caseInput := HandleInput{
		UserID:  "test user",
		Message: "test message",
	}
	caseError := errors.New("test error")

	ml := log.NewMockLogger(ctrl)
	gomock.InOrder(
		ml.EXPECT().Info(gomock.Any(), gomock.Any()),
		ml.EXPECT().Info(gomock.Any()),
	)
	log.SetDefaultLogger(ml)

	i := notifyevent.NewMockInputPort(ctrl)
	i.EXPECT().NotifyEvent(gomock.Eq(ctx), gomock.Eq(notifyevent.InputData{
		UserID:  caseInput.UserID,
		Message: caseInput.Message,
	})).Return(caseError)

	d := di.NewMockDI(ctrl)
	d.EXPECT().Get("notifyevent.InputPort").Return(i)

	di.SetDi(d)

	h := NewController()
	resp := h.Handle(ctx, caseInput)
	assert.Equal(t, &Response{
		Error: caseError,
	}, resp)
}
