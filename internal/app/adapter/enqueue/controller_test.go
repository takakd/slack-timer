package enqueue

import (
	"context"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mi := enqueueevent.NewMockInputPort(ctrl)
	md := di.NewMockDI(ctrl)
	md.EXPECT().Get(gomock.Eq("enqueueevent.InputPort")).Return(mi)

	di.SetDi(md)

	h := NewController()
	assert.Equal(t, mi, h.InputPort)
}

func TestController_Handle(t *testing.T) {
	ctx := context.TODO()
	caseInput := HandleInput{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ml := log.NewMockLogger(ctrl)
	gomock.InOrder(
		ml.EXPECT().Info(gomock.Any(), gomock.Any()),
		ml.EXPECT().Info(gomock.Any()),
	)
	log.SetDefaultLogger(ml)

	i := enqueueevent.NewMockInputPort(ctrl)
	i.EXPECT().EnqueueEvent(gomock.Eq(ctx), gomock.Any())

	d := di.NewMockDI(ctrl)
	d.EXPECT().Get("enqueueevent.InputPort").Return(i)
	di.SetDi(d)

	assert.NotPanics(t, func() {
		h := NewController()
		h.Handle(ctx, caseInput)
	})
}
