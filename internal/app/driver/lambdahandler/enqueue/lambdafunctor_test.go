package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewLambdaFunctor(t *testing.T) {
	assert.NotPanics(t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mc := enqueue.NewMockControllerHandler(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("enqueue.ControllerHandler").Return(mc)
		di.SetDi(md)

		NewLambdaFunctor()
	})
}

func TestLambdaFunctor_Handle(t *testing.T) {
	t.Run("ok:notify", func(t *testing.T) {
		assert.NotPanics(t, func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			lc := &lambdacontext.LambdaContext{}
			ctx := lambdacontext.NewContext(context.TODO(), lc)
			ac, _ := appcontext.FromContext(ctx)

			caseInput := LambdaInput{}

			mc := enqueue.NewMockControllerHandler(ctrl)
			mc.EXPECT().Handle(ac, gomock.Any())

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("enqueue.ControllerHandler").Return(mc)
			di.SetDi(md)

			h := NewLambdaFunctor()
			h.Handle(ctx, caseInput)
		})
	})
}
