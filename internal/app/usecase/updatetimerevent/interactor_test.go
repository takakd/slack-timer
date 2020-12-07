package updatetimerevent

import (
	"context"
	"errors"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInteractor_saveTimerEventValue(t *testing.T) {

	t.Run("ok:create", func(t *testing.T) {
		ctx := context.TODO()
		userID := "abc"
		caseTime := time.Now().UTC()

		caseEvent, _ := enterpriserule.NewTimerEvent(userID)
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = caseTime.Add(time.Duration(caseEvent.IntervalMin) * time.Minute)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userID)).
			Return(nil, nil)

		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Any()).DoAndReturn(func(_ context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error) {
			return event, nil
		})

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, userID, caseTime, caseEvent.IntervalMin)
		assert.NoError(t, data.Result)
		assert.Equal(t, caseEvent, data.SavedEvent)
	})

	t.Run("ok:next notify", func(t *testing.T) {
		ctx := context.TODO()
		caseTime := time.Now().UTC()

		caseEvent, _ := enterpriserule.NewTimerEvent("abc")
		caseEvent.NotificationTime = caseTime
		caseEvent.IntervalMin = 10

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent.UserID)).
			Return(caseEvent, nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		want, _ := enterpriserule.NewTimerEvent(caseEvent.UserID)
		want.IntervalMin = caseEvent.IntervalMin
		want.NotificationTime = caseTime.Add(time.Duration(want.IntervalMin) * time.Minute)

		data := interactor.saveTimerEventValue(ctx, caseEvent.UserID, caseTime, 0)

		assert.NoError(t, data.Result)
		assert.Equal(t, caseEvent, data.SavedEvent)
		assert.Equal(t, want, data.SavedEvent)
	})

	t.Run("ok:interval", func(t *testing.T) {
		caseTime := time.Now().UTC()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userID := "abc"
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userID)).
			Return(nil, nil)

		interval := 1
		caseEvent, _ := enterpriserule.NewTimerEvent(userID)
		caseEvent.NotificationTime = caseTime.Add(time.Duration(interval) * time.Minute)
		caseEvent.IntervalMin = interval
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		data := interactor.saveTimerEventValue(ctx, userID, caseTime, interval)
		assert.NoError(t, data.Result)
		assert.Equal(t, caseEvent, data.SavedEvent)
	})

	t.Run("ng:create", func(t *testing.T) {
		ctx := context.TODO()
		userID := ""

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userID)).
			Return(nil, nil)

		interactor := &Interactor{
			repository: m,
		}

		noUse := time.Now()
		data := interactor.saveTimerEventValue(context.TODO(), userID, noUse, 0)
		assert.Equal(t, fmt.Errorf("creating timer event error userID=%v: %w", userID, errors.New("must set userID")), data.Result)
	})

	t.Run("ng:update", func(t *testing.T) {
		ctx := context.TODO()
		userID := ""
		err := errors.New("error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(userID)).
			Return(nil, err)

		interactor := &Interactor{
			repository: m,
		}

		noUse := time.Now()
		data := interactor.saveTimerEventValue(context.TODO(), userID, noUse, 0)
		assert.Equal(t, fmt.Errorf("finding timer event error userID=%v: %w", userID, err), data.Result)
	})
}
