package notifycontroller

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i := notifyevent.NewMockInputPort(ctrl)
	d := di.NewMockDI(ctrl)
	d.EXPECT().Get("notifycontroller.InputPort").Return(i)

	di.SetDi(d)

	h := NewHandler().(*SqsEventHandler)
	assert.Equal(t, i, h.InputPort)
}

func TestSqsEventHandler_Handler(t *testing.T) {
	ctx := context.TODO()
	caseInput := HandlerInput{
		UserId:  "test user",
		Message: "test message",
	}
	caseError := errors.New("test error")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i := notifyevent.NewMockInputPort(ctrl)
	i.EXPECT().NotifyEvent(gomock.Eq(ctx), gomock.Eq(notifyevent.InputData{
		UserId:  caseInput.UserId,
		Message: caseInput.Message,
	})).Return(caseError)

	d := di.NewMockDI(ctrl)
	d.EXPECT().Get("notifycontroller.InputPort").Return(i)

	di.SetDi(d)

	h := NewHandler().(*SqsEventHandler)
	resp := h.Handler(ctx, caseInput)
	assert.Equal(t, &Response{
		Error: caseError,
	}, resp)
}
