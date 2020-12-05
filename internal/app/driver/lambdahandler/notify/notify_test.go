package notify

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/adapter/notifycontroller"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestSqsMessage_HandlerInput(t *testing.T) {
	m := &SqsMessage{
		Body: "test user",
	}
	h := m.HandlerInput()
	assert.Equal(t, m.Body, h.UserId)
	assert.Equal(t, "test user", h.UserId)
}

func TestLambdaHandler(t *testing.T) {
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
		caseResponse := &notifycontroller.Response{
			Error: nil,
		}

		mi := notifycontroller.NewMockHandler(ctrl)
		mi.EXPECT().Handler(gomock.Eq(ctx), gomock.Any()).Return(caseResponse)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.Handler").Return(mi)
		di.SetDi(md)

		h := NewNotifyLambdaHandler()
		err := h.LambdaHandler(ctx, caseInput)
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

		caseResponse := &notifycontroller.Response{
			Error: errors.New("test error"),
		}

		mi := notifycontroller.NewMockHandler(ctrl)
		mi.EXPECT().Handler(gomock.Eq(ctx), gomock.Any()).Return(caseResponse)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("notify.Handler").Return(mi)
		di.SetDi(md)

		h := NewNotifyLambdaHandler()
		err := h.LambdaHandler(ctx, caseInput)
		assert.Error(t, errors.New("error happend count=1"), err)
	})
}
