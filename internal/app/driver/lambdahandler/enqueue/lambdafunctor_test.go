package enqueue

import (
	"context"
	"slacktimer/internal/app/adapter/enqueue"
	"slacktimer/internal/app/util/di"
	"testing"

	"slacktimer/internal/app/util/appcontext"

	"time"

	"fmt"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type AppContextMatcher struct {
	testValue appcontext.AppContext
}

func (m *AppContextMatcher) String() string {
	return fmt.Sprintf("%v", m.testValue)
}
func (m *AppContextMatcher) Matches(x interface{}) bool {
	another, _ := x.(appcontext.AppContext)
	matched := true
	matched = matched && m.testValue.RequestID() == another.RequestID()
	return matched
}

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
			ac, _ := appcontext.NewLambdaAppContext(ctx, time.Now())

			caseInput := LambdaInput{}

			mc := enqueue.NewMockControllerHandler(ctrl)
			matcher := &AppContextMatcher{
				testValue: ac,
			}
			mc.EXPECT().Handle(matcher, gomock.Any())

			md := di.NewMockDI(ctrl)
			md.EXPECT().Get("enqueue.ControllerHandler").Return(mc)
			di.SetDi(md)

			h := NewLambdaFunctor()
			h.Handle(ctx, caseInput)
		})
	})
}
