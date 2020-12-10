package notify

import (
	"context"
	"errors"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/util/di"
	"testing"

	"encoding/json"
	"slacktimer/internal/app/driver/queue"

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

		mc := notify.NewMockControllerHandler(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.ControllerHandler").Return(mc)
		di.SetDi(md)

		NewLambdaFunctor()
	})
}

func TestLambdaFunctor_Handle(t *testing.T) {
	t.Run("ok:notify", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := &lambdacontext.LambdaContext{}
		ctx := lambdacontext.NewContext(context.TODO(), lc)
		ac, _ := appcontext.NewLambdaAppContext(ctx, time.Now())

		caseBody, err := json.Marshal(queue.SqsMessageBody{
			UserID: "test user",
			Text:   "test text",
		})
		caseInput := LambdaInput{
			Records: []SqsMessage{
				{
					Body: string(caseBody),
				},
			},
		}
		caseResponse := &notify.Response{
			Error: nil,
		}

		mi := notify.NewMockControllerHandler(ctrl)
		matcher := &AppContextMatcher{
			testValue: ac,
		}
		mi.EXPECT().Handle(matcher, gomock.Any()).Return(caseResponse)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.ControllerHandler").Return(mi)
		di.SetDi(md)

		h := NewLambdaFunctor()
		err = h.Handle(ctx, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ng:notify", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		lc := &lambdacontext.LambdaContext{}
		ctx := lambdacontext.NewContext(context.TODO(), lc)
		ac, _ := appcontext.NewLambdaAppContext(ctx, time.Now())

		caseBody, err := json.Marshal(queue.SqsMessageBody{
			UserID: "test user",
			Text:   "test text",
		})
		caseInput := LambdaInput{
			Records: []SqsMessage{
				{
					Body: string(caseBody),
				},
			},
		}

		caseResponse := &notify.Response{
			Error: errors.New("test error"),
		}

		mi := notify.NewMockControllerHandler(ctrl)
		matcher := &AppContextMatcher{
			testValue: ac,
		}
		mi.EXPECT().Handle(matcher, gomock.Any()).Return(caseResponse)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.ControllerHandler").Return(mi)
		di.SetDi(md)

		h := NewLambdaFunctor()
		err = h.Handle(ctx, caseInput)
		assert.Error(t, errors.New("error happend count=1"), err)
	})
}
