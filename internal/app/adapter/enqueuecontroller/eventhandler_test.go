package enqueuecontroller

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/appinit"
	"testing"
)

func TestCloudWatchEventHandler_Handler(t *testing.T) {
	appinit.AppInit()

	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := enqueueevent.NewMockUsecase(ctrl)
	m.EXPECT().EnqueueEvent(gomock.Eq(ctx), gomock.Any())

	h := CloudWatchEventHandler{
		usecase: m,
	}
	resp := h.Handler(ctx)
	// TODO: modify along codes.
	assert.Equal(t, resp.Error, nil)
}
