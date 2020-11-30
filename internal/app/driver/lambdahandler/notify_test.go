package lambdahandler

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"slacktimer/internal/app/adapter/notifycontroller"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/usecase/notifyevent"
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

func TestNotifyLambdaHandler(t *testing.T) {
	t.Run("ok:notify", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := LambdaInput{
			Records: []SqsMessage{
				{
					Body: "test user",
				},
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseResponse := &notifycontroller.Response{
			Error: nil,
		}

		i := notifyevent.NewMockInputPort(ctrl)
		i.EXPECT().NotifyEvent(gomock.Eq(ctx), gomock.Any()).Return(caseResponse.Error)

		m := di.NewMockDI(ctrl)
		m.EXPECT().Get("notifycontroller.InputPort").Return(i)
		di.SetDi(m)

		os.Setenv("APP_ENV", "ignore set DI")
		err := NotifyLambdaHandler(ctx, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ok:notify", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := LambdaInput{
			Records: []SqsMessage{
				{
					Body: "test_user",
				},
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseResponse := &notifycontroller.Response{
			Error: errors.New("test error"),
		}

		i := notifyevent.NewMockInputPort(ctrl)
		i.EXPECT().NotifyEvent(gomock.Eq(ctx), gomock.Any()).Return(caseResponse.Error)

		m := di.NewMockDI(ctrl)
		m.EXPECT().Get("notifycontroller.InputPort").Return(i)
		di.SetDi(m)

		os.Setenv("APP_ENV", "ignore set DI")
		err := NotifyLambdaHandler(ctx, caseInput)
		assert.Error(t, errors.New("error happend count=1"), err)
	})
}
