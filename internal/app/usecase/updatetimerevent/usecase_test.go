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
		caseTime := time.Now().UTC()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		caseEvent, _ := enterpriserule.NewTimerEvent(userId)
		caseEvent.NotificationTime = caseTime
		caseEvent.NotificationTime.Add(time.Duration(caseEvent.IntervalMin) * time.Minute)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(_ context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error) {
			return event, nil
		})

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, userId, caseTime, caseEvent.IntervalMin)
		assert.NoError(t, data.Result)
		assert.Equal(t, caseEvent, data.SavedEvent)
	})

	t.Run("ok:next notify", func(t *testing.T) {
		ctx := context.TODO()
		caseEvent, _ := enterpriserule.NewTimerEvent("abc")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent.UserId)).
			Return(caseEvent, nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		noUse := time.Now().UTC()
		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserId)
		want.NotificationTime = want.NotificationTime.Add(time.Duration(want.IntervalMin) * time.Minute)

		data := interactor.saveTimerEventValue(ctx, caseEvent.UserId, noUse, 0)

		assert.NoError(t, data.Result)
		assert.Equal(t, caseEvent, data.SavedEvent)
		assert.Equal(t, want, data.SavedEvent)
	})

	t.Run("ok:interval", func(t *testing.T) {
		caseTime := time.Now().UTC()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc"
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userId)).
			Return(nil, nil)

		interval := 1
		caseEvent, _ := enterpriserule.NewTimerEvent(userId)
		caseEvent.NotificationTime = caseTime.Add(time.Duration(interval) * time.Minute)
		caseEvent.IntervalMin = interval
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, userId, caseTime, interval)
		assert.NoError(t, data.Result)
		assert.Equal(t, caseEvent, data.SavedEvent)
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

		noUse := time.Now()
		data := interactor.saveTimerEventValue(context.TODO(), userId, noUse, 0)
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

		noUse := time.Now()
		data := interactor.saveTimerEventValue(context.TODO(), userId, noUse, 0)
		assert.True(t, errors.Is(data.Result, ErrFind))
	})
}
