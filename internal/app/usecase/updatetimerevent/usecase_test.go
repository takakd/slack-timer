package updatetimerevent

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/enterpriserule"
	"testing"
	"time"
)

func TestInteractor_saveTimerEventValue(t *testing.T) {

	t.Run("ok:create", func(t *testing.T) {
		ctx := context.TODO()
		userId := "abc"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		now := time.Now()
		event, _ := enterpriserule.NewTimerEvent(userId)
		event.NotificationTime = now.UTC()
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(_, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error) {

		})

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, userId, 0)
		assert.NoError(t, data.Result)
	})

	t.Run("ok:next notify", func(t *testing.T) {
		ctx := context.TODO()
		event, _ := enterpriserule.NewTimerEvent("abc")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(event.UserId)).
			Return(event, nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(event)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, event.UserId, 0)
		assert.NoError(t, data.Result)
	})

	t.Run("ok:interval", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc"
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		interval := 1
		event, _ := enterpriserule.NewTimerEvent(userId)
		event.IntervalMin = interval
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(event)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, userId, interval)
		assert.NoError(t, data.Result)
	})

	t.Run("ng:create", func(t *testing.T) {
		ctx := context.TODO()
		userId := ""

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(context.TODO(), userId, 0)
		assert.True(t, errors.Is(data.Result, ErrCreate))
	})

	t.Run("ng:update", func(t *testing.T) {
		ctx := context.TODO()
		userId := ""

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, errors.New("error"))

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(context.TODO(), userId, 0)
		assert.True(t, errors.Is(data.Result, ErrFind))
	})
}
