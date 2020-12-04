package enqueuecontroller

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i := enqueueevent.NewMockInputPort(ctrl)
	d := di.NewMockDI(ctrl)
	d.EXPECT().Get(gomock.Eq("enqueueevent.InputPort")).Return(i)

	di.SetDi(d)

	h := NewHandler().(*CloudWatchEventHandler)
	assert.Equal(t, i, h.InputPort)
}

func TestCloudWatchEventHandler_Handler(t *testing.T) {
	ctx := context.TODO()
	caseInput := HandlerInput{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i := enqueueevent.NewMockInputPort(ctrl)
	i.EXPECT().EnqueueEvent(gomock.Eq(ctx), gomock.Any())

	d := di.NewMockDI(ctrl)
	d.EXPECT().Get("enqueueevent.InputPort").Return(i)
	di.SetDi(d)

	h := NewHandler().(*CloudWatchEventHandler)
	resp := h.Handler(ctx, caseInput)
	assert.Equal(t, &Response{}, resp)

	//os.Setenv("APP_ENV", "ignore set DI")
	//err := LambdaHandleEvent(ctx, caseInput)
	//assert.Equal(t, caseResponse.Error, err)
}
