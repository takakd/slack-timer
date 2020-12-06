package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/util/di"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewEnqueueLambdaHandler(t *testing.T) {
	assert.NotPanics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := enqueue.NewMockControllerHandler(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("enqueue.Controller").Return(mc)
		di.SetDi(md)

		NewEnqueueLambdaHandler()
	})
}

func TestEnqueueLambdaHandler_Handle(t *testing.T) {
	t.Run("ok:notify", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()

			caseInput := LambdaInput{}

			mc := enqueue.NewMockControllerHandler(ctrl)
			mc.EXPECT().Handle(gomock.Eq(ctx), gomock.Any())

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("enqueue.Controller").Return(mc)
			di.SetDi(md)

			h := NewEnqueueLambdaHandler()
			h.Handle(ctx, caseInput)
		})
	})
}
