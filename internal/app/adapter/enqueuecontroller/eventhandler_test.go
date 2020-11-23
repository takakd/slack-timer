package enqueuecontroller

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/usecase/enqueueevent"
	"testing"
)

func TestCloudWatchEventHandler_Handler(t *testing.T) {
	ctx := context.TODO()
	caseError := errors.New("dummy")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := enqueueevent.NewMockUsecase(ctrl)
	m.EXPECT().EnqueueEvent(gomock.Eq(ctx), gomock.Any()).Return(caseError)

	h := CloudWatchEventHandler{
		usecase: m,
	}
	resp := h.Handler(ctx)
	assert.Equal(t, resp.Error, caseError)
}
