package enqueuecontroller

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/util/di"
	"testing"
)

func TestNewEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := di.NewMockDI(ctrl)

	caseUseCase := &enqueueevent.Interactor{}
	m.EXPECT().Get("enqueuecontroller.EnqueueNotification").Return(caseUseCase)
	di.SetDi(m)

	h := NewEventHandler()
	assert.Equal(t, &CloudWatchEventHandler{caseUseCase}, h)
}

func TestLambdaHandleEvent(t *testing.T) {
	ctx := context.TODO()
	caseInput := LambdaInput{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	caseResponse := &HandlerResponse{
		Error: errors.New("test error"),
	}
	u := enqueueevent.NewMockUsecase(ctrl)
	u.EXPECT().EnqueueEvent(gomock.Eq(ctx), gomock.Any()).Return(caseResponse.Error)

	m := di.NewMockDI(ctrl)
	m.EXPECT().Get("enqueuecontroller.EnqueueNotification").Return(u)
	di.SetDi(m)

	os.Setenv("APP_ENV", "ignore set DI")
	err := LambdaHandleEvent(ctx, caseInput)
	assert.Equal(t, caseResponse.Error, err)
}
