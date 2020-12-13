package enqueueevent

import (
	"errors"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"
	"time"

	"slacktimer/internal/app/util/appcontext"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewInteractor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := NewMockRepository(ctrl)
	o := NewMockOutputPort(ctrl)
	q := NewMockQueue(ctrl)
	d := di.NewMockDI(ctrl)

	d.EXPECT().Get("enqueueevent.OutputPort").Return(o)
	d.EXPECT().Get("enqueueevent.Repository").Return(r)
	d.EXPECT().Get("enqueueevent.Queue").Return(q)
	di.SetDi(d)

	i := NewInteractor()
	assert.Equal(t, o, i.outputPort)
	assert.Equal(t, r, i.repository)
	assert.Equal(t, q, i.queue)
}

func TestInteractor_EnqueueEvent(t *testing.T) {
	t.Run("ok:enqueue", func(t *testing.T) {
		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		caseEvents[0], _ = enterpriserule.NewTimerEvent("id1")
		caseEvents[0].State = enterpriserule.TimerEventStateWait
		caseEvents[1], _ = enterpriserule.NewTimerEvent("id2")
		caseEvents[1].State = enterpriserule.TimerEventStateWait
		caseQueueMsg := []QueueMessage{
			{
				UserID: caseEvents[0].UserID(),
			},
			{
				UserID: caseEvents[1].UserID(),
			},
		}
		caseOutputData := OutputData{
			NotifiedUserIDList: []string{
				caseEvents[0].UserID(),
				caseEvents[1].UserID(),
			},
			QueueMessageIDList: []string{
				"mid1", "mid2",
			},
		}

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(caseInput.EventTime).
			Return(caseEvents, nil)
		m.EXPECT().SaveTimerEvent(caseEvents[0]).Return(caseEvents[0], nil)
		m.EXPECT().SaveTimerEvent(caseEvents[1]).Return(caseEvents[1], nil)

		q := NewMockQueue(ctrl)
		q.EXPECT().Enqueue(caseQueueMsg[0]).Return(caseOutputData.QueueMessageIDList[0], nil)
		q.EXPECT().Enqueue(caseQueueMsg[1]).Return(caseOutputData.QueueMessageIDList[1], nil)

		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(appcontext.TODO(), caseOutputData)

		interactor := &Interactor{
			repository: m,
			outputPort: o,
			queue:      q,
		}

		interactor.EnqueueEvent(appcontext.TODO(), caseInput)
	})

	t.Run("ng:failed find", func(t *testing.T) {
		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseFindError := errors.New("repository error")
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(caseInput.EventTime).
			Return(nil, caseFindError)

		caseOutputData := OutputData{
			Result: fmt.Errorf("find error time=%v: %w", caseInput.EventTime, caseFindError),
		}
		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(appcontext.TODO(), caseOutputData)

		interactor := &Interactor{
			repository: m,
			outputPort: o,
		}

		interactor.EnqueueEvent(appcontext.TODO(), caseInput)
	})

	t.Run("ng:exist failed enqueue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		caseEvents[0], _ = enterpriserule.NewTimerEvent("id1")
		caseEvents[0].State = enterpriserule.TimerEventStateWait
		caseEvents[1], _ = enterpriserule.NewTimerEvent("id2")
		caseEvents[1].State = enterpriserule.TimerEventStateWait
		caseQueueMsg := []QueueMessage{
			{
				UserID: caseEvents[0].UserID(),
			},
			{
				UserID: caseEvents[1].UserID(),
			},
		}
		caseOutputData := OutputData{
			NotifiedUserIDList: []string{
				caseEvents[0].UserID(),
			},
			QueueMessageIDList: []string{
				"mid1",
			},
		}
		caseError := errors.New("error")

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(caseInput.EventTime).
			Return(caseEvents, nil)
		m.EXPECT().SaveTimerEvent(caseEvents[0]).Return(caseEvents[0], nil)

		q := NewMockQueue(ctrl)
		q.EXPECT().Enqueue(caseQueueMsg[0]).Return(caseOutputData.QueueMessageIDList[0], nil)
		q.EXPECT().Enqueue(caseQueueMsg[1]).Return("", caseError)

		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(appcontext.TODO(), caseOutputData)

		l := log.NewMockLogger(ctrl)
		l.EXPECT().ErrorWithContext(ac, fmt.Sprintf("enqueue error user_id=%s: %v", caseEvents[1].UserID(), caseError))
		log.SetDefaultLogger(l)

		interactor := &Interactor{
			repository: m,
			outputPort: o,
			queue:      q,
		}

		interactor.EnqueueEvent(ac, caseInput)
	})

	t.Run("ng:update error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ac := appcontext.TODO()

		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		caseEvents[0], _ = enterpriserule.NewTimerEvent("id1")
		caseEvents[0].State = enterpriserule.TimerEventStateWait
		caseEvents[1], _ = enterpriserule.NewTimerEvent("id2")
		caseEvents[1].State = enterpriserule.TimerEventStateWait
		caseQueueMsg := []QueueMessage{
			{
				UserID: caseEvents[0].UserID(),
			},
			{
				UserID: caseEvents[1].UserID(),
			},
		}
		caseOutputData := OutputData{
			NotifiedUserIDList: []string{
				caseEvents[0].UserID(),
			},
			QueueMessageIDList: []string{
				"mid1",
			},
		}
		caseError := errors.New("error")

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(caseInput.EventTime).
			Return(caseEvents, nil)
		m.EXPECT().SaveTimerEvent(caseEvents[0]).Return(caseEvents[0], nil)
		m.EXPECT().SaveTimerEvent(caseEvents[1]).Return(nil, caseError)

		q := NewMockQueue(ctrl)
		q.EXPECT().Enqueue(caseQueueMsg[0]).Return(caseOutputData.QueueMessageIDList[0], nil)
		q.EXPECT().Enqueue(caseQueueMsg[1]).Return("dummy", nil)

		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(appcontext.TODO(), caseOutputData)

		l := log.NewMockLogger(ctrl)
		l.EXPECT().ErrorWithContext(ac, fmt.Sprintf("update error user_id=%s: %v", caseEvents[1].UserID(), caseError))
		log.SetDefaultLogger(l)

		interactor := &Interactor{
			repository: m,
			outputPort: o,
			queue:      q,
		}

		interactor.EnqueueEvent(ac, caseInput)
	})
}
