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

	mi := enqueueevent.NewMockInputPort(ctrl)
	md := di.NewMockDI(ctrl)
	md.EXPECT().Get(gomock.Eq("enqueueevent.InputPort")).Return(mi)

	di.SetDi(md)

	h := NewHandler().(*CloudWatchEventHandler)
	assert.Equal(t, mi, h.InputPort)
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

	assert.NotPanics(t, func() {
		h := NewHandler()
		h.Handler(ctx, caseInput)
	})
}
