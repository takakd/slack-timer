package notifyevent

import (
	"context"
	"errors"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewInteractor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	o := NewMockOutputPort(ctrl)
	r := NewMockRepository(ctrl)
	n := NewMockNotifier(ctrl)
	d := di.NewMockDI(ctrl)

	d.EXPECT().Get("notifyevent.OutputPort").Return(o)
	d.EXPECT().Get("notifyevent.Repository").Return(r)
	d.EXPECT().Get("notifyevent.Notifier").Return(n)
	di.SetDi(d)

	i := NewInteractor()
	assert.Equal(t, o, i.outputPort)
	assert.Equal(t, r, i.repository)
	assert.Equal(t, n, i.notifier)
}

func TestInteractor_NotifyEvent(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}
		caseEvent := &enterpriserule.TimerEvent{
			UserID:           caseInput.UserID,
			IntervalMin:      10,
			NotificationTime: time.Now(),
		}
		caseEvent.SetQueued()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(gomock.Eq(OutputData{
			Result: nil,
			UserID: caseInput.UserID,
		}))

		r.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseInput.UserID)).Return(caseEvent, nil)

		n.EXPECT().Notify(gomock.Eq(caseInput.UserID), gomock.Eq(caseInput.Message)).Return(nil)

		r.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent)).Return(caseEvent, nil)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err := i.NotifyEvent(ctx, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ng:already notified", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}

		caseEvent := &enterpriserule.TimerEvent{
			UserID:           caseInput.UserID,
			IntervalMin:      10,
			NotificationTime: time.Now(),
		}
		caseEvent.SetWait()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(gomock.Eq(OutputData{
			UserID: caseInput.UserID,
		}))

		r.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseInput.UserID)).Return(caseEvent, nil)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err := i.NotifyEvent(ctx, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ng:notify", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}

		caseEvent := &enterpriserule.TimerEvent{
			UserID:           caseInput.UserID,
			IntervalMin:      10,
			NotificationTime: time.Now(),
		}
		caseEvent.SetQueued()

		caseError := errors.New("notify error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(gomock.Eq(OutputData{
			Result: caseError,
			UserID: caseInput.UserID,
		}))

		r.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseInput.UserID)).Return(caseEvent, nil)

		n.EXPECT().Notify(gomock.Eq(caseInput.UserID), gomock.Eq(caseInput.Message)).Return(caseError)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err := i.NotifyEvent(ctx, caseInput)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:find", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}
		caseError := errors.New("find error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(gomock.Eq(OutputData{
			Result: caseError,
			UserID: caseInput.UserID,
		}))

		r.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseInput.UserID)).Return(nil, caseError)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err := i.NotifyEvent(ctx, caseInput)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:save", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}
		caseEvent := &enterpriserule.TimerEvent{
			UserID:           caseInput.UserID,
			IntervalMin:      10,
			NotificationTime: time.Now(),
		}
		caseEvent.SetQueued()

		caseError := errors.New("save error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(gomock.Eq(OutputData{
			Result: caseError,
			UserID: caseInput.UserID,
		}))

		r.EXPECT().FindTimerEvent(gomock.Eq(ctx), gomock.Eq(caseInput.UserID)).Return(caseEvent, nil)

		n.EXPECT().Notify(gomock.Eq(caseInput.UserID), gomock.Eq(caseInput.Message)).Return(nil)

		r.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvent)).Return(nil, caseError)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err := i.NotifyEvent(ctx, caseInput)
		assert.Equal(t, caseError, err)
	})
}
