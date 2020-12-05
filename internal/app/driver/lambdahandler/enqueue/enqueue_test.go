package enqueue

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/adapter/enqueuecontroller"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestLambdaInput_HandlerInput(t *testing.T) {
	caseInput := LambdaInput{}
	assert.Equal(t, enqueuecontroller.HandlerInput{}, caseInput.HandlerInput())
}

func TestLambdaHandler(t *testing.T) {
	t.Run("ok:notify", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()

			caseInput := LambdaInput{}

			mc := enqueuecontroller.NewMockHandler(ctrl)
			mc.EXPECT().Handler(gomock.Eq(ctx), gomock.Any())

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("enqueue.Handler").Return(mc)
			di.SetDi(md)

			h := NewEnqueueLambdaHandler()
			h.LambdaHandler(ctx, caseInput)
		})
	})
}
