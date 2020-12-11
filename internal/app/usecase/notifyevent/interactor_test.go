package notifyevent

import (
	"errors"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"testing"
	"time"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		ac := appcontext.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}
		caseEvent, err := enterpriserule.NewTimerEvent(caseInput.UserID)
		require.NoError(t, err)
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = time.Now()
		caseEvent.State = enterpriserule.TimerEventStateQueued

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(ac, OutputData{
			Result: nil,
			UserID: caseInput.UserID,
		})

		r.EXPECT().FindTimerEvent(caseInput.UserID).Return(caseEvent, nil)

		n.EXPECT().Notify(ac, caseInput.UserID, caseInput.Message).Return(nil)

		r.EXPECT().SaveTimerEvent(caseEvent).Return(caseEvent, nil)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err = i.NotifyEvent(ac, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ng:already notified", func(t *testing.T) {
		ac := appcontext.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}

		caseEvent, err := enterpriserule.NewTimerEvent(caseInput.UserID)
		require.NoError(t, err)
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = time.Now()
		caseEvent.State = enterpriserule.TimerEventStateWait

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(ac, OutputData{
			UserID: caseInput.UserID,
		})

		r.EXPECT().FindTimerEvent(caseInput.UserID).Return(caseEvent, nil)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err = i.NotifyEvent(ac, caseInput)
		assert.NoError(t, err)
	})

	t.Run("ng:notify", func(t *testing.T) {
		ac := appcontext.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}

		caseEvent, err := enterpriserule.NewTimerEvent(caseInput.UserID)
		require.NoError(t, err)
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = time.Now()
		caseEvent.State = enterpriserule.TimerEventStateQueued

		caseError := errors.New("notify error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(ac, OutputData{
			Result: caseError,
			UserID: caseInput.UserID,
		})

		r.EXPECT().FindTimerEvent(caseInput.UserID).Return(caseEvent, nil)

		n.EXPECT().Notify(ac, caseInput.UserID, caseInput.Message).Return(caseError)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err = i.NotifyEvent(ac, caseInput)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:find", func(t *testing.T) {
		ac := appcontext.TODO()
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

		o.EXPECT().Output(ac, OutputData{
			Result: caseError,
			UserID: caseInput.UserID,
		})

		r.EXPECT().FindTimerEvent(caseInput.UserID).Return(nil, caseError)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err := i.NotifyEvent(ac, caseInput)
		assert.Equal(t, caseError, err)
	})

	t.Run("ng:save", func(t *testing.T) {
		ac := appcontext.TODO()
		caseInput := InputData{
			UserID:  "test",
			Message: "message",
		}

		caseEvent, err := enterpriserule.NewTimerEvent(caseInput.UserID)
		require.NoError(t, err)
		caseEvent.IntervalMin = 10
		caseEvent.NotificationTime = time.Now()
		caseEvent.State = enterpriserule.TimerEventStateQueued

		caseError := errors.New("save error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		o := NewMockOutputPort(ctrl)
		r := NewMockRepository(ctrl)
		n := NewMockNotifier(ctrl)
		d := di.NewMockDI(ctrl)

		o.EXPECT().Output(ac, OutputData{
			Result: caseError,
			UserID: caseInput.UserID,
		})

		r.EXPECT().FindTimerEvent(caseInput.UserID).Return(caseEvent, nil)

		n.EXPECT().Notify(ac, caseInput.UserID, caseInput.Message).Return(nil)

		r.EXPECT().SaveTimerEvent(caseEvent).Return(nil, caseError)

		d.EXPECT().Get("notifyevent.OutputPort").Return(o)
		d.EXPECT().Get("notifyevent.Repository").Return(r)
		d.EXPECT().Get("notifyevent.Notifier").Return(n)

		di.SetDi(d)

		i := NewInteractor()
		err = i.NotifyEvent(ac, caseInput)
		assert.Equal(t, caseError, err)
	})
}
