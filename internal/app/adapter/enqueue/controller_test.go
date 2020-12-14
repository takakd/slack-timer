package enqueue

import (
	"slacktimer/internal/app/usecase/enqueueevent"
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

	mi := enqueueevent.NewMockInputPort(ctrl)
	md := di.NewMockDI(ctrl)
	md.EXPECT().Get("enqueueevent.InputPort").Return(mi)

	di.SetDi(md)

	h := NewController()
	assert.Equal(t, mi, h.inputPort)
}

func TestController_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ac := appcontext.TODO()
	caseInput := HandleInput{}

	ml := log.NewMockLogger(ctrl)
	gomock.InOrder(
		ml.EXPECT().InfoWithContext(ac, gomock.Any(), gomock.Any()),
		ml.EXPECT().InfoWithContext(ac, gomock.Any()),
	)
	log.SetDefaultLogger(ml)

	mi := enqueueevent.NewMockInputPort(ctrl)
	mi.EXPECT().EnqueueEvent(ac, gomock.Any())

	md := di.NewMockDI(ctrl)
	md.EXPECT().Get("enqueueevent.InputPort").Return(mi)
	di.SetDi(md)

	assert.NotPanics(t, func() {
		h := NewController()
		h.Handle(ac, caseInput)
	})
}
