package notify

import (
	"context"
	"errors"
	"slacktimer/internal/app/adapter/notify"
	"slacktimer/internal/app/util/di"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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

		ctx := context.TODO()

		caseInput := LambdaInput{
			Records: []SqsMessage{
				{
					Body: "test user",
				},
			},
		}
		caseResponse := &notify.Response{
			Error: nil,
		}

		mi := notify.NewMockControllerHandler(ctrl)
		mi.EXPECT().Handle(gomock.Eq(ctx), gomock.Any()).Return(caseResponse)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.ControllerHandler").Return(mi)
		di.SetDi(md)

		h := NewLambdaFunctor()
		err := h.Handle(ctx, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ng:notify", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()

		caseInput := LambdaInput{
			Records: []SqsMessage{
				{
					Body: "test_user",
				},
			},
		}

		caseResponse := &notify.Response{
			Error: errors.New("test error"),
		}

		mi := notify.NewMockControllerHandler(ctrl)
		mi.EXPECT().Handle(gomock.Eq(ctx), gomock.Any()).Return(caseResponse)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.ControllerHandler").Return(mi)
		di.SetDi(md)

		h := NewLambdaFunctor()
		err := h.Handle(ctx, caseInput)
		assert.Error(t, errors.New("error happend count=1"), err)
	})
}
