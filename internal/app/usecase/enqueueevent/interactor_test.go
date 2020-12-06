package enqueueevent

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"testing"
	"time"
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
		ctx := context.TODO()
		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		caseEvents[0], _ = enterpriserule.NewTimerEvent("id1")
		caseEvents[1], _ = enterpriserule.NewTimerEvent("id2")
		caseQueueMsg := []QueueMessage{
			{
				caseEvents[0].UserId,
			},
			{
				caseEvents[1].UserId,
			},
		}
		caseOutputData := OutputData{
			NotifiedUserIdList: []string{
				caseEvents[0].UserId,
				caseEvents[1].UserId,
			},
			QueueMessageIdList: []string{
				"mid1", "mid2",
			},
		}

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(gomock.Eq(ctx), gomock.Eq(caseInput.EventTime)).
			Return(caseEvents, nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvents[0])).Return(caseEvents[0], nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvents[1])).Return(caseEvents[1], nil)

		q := NewMockQueue(ctrl)
		q.EXPECT().Enqueue(gomock.Eq(caseQueueMsg[0])).Return(caseOutputData.QueueMessageIdList[0], nil)
		q.EXPECT().Enqueue(gomock.Eq(caseQueueMsg[1])).Return(caseOutputData.QueueMessageIdList[1], nil)

		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(gomock.Eq(caseOutputData))

		interactor := &Interactor{
			repository: m,
			outputPort: o,
			queue:      q,
		}

		interactor.EnqueueEvent(ctx, caseInput)
	})

	t.Run("ng:failed find", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseFindError := errors.New("repository error")
		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(gomock.Eq(ctx), gomock.Eq(caseInput.EventTime)).
			Return(nil, caseFindError)

		caseOutputData := OutputData{
			Result: fmt.Errorf("find error time=%v: %w", caseInput.EventTime, caseFindError),
		}
		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(gomock.Eq(caseOutputData))

		interactor := &Interactor{
			repository: m,
			outputPort: o,
		}

		interactor.EnqueueEvent(ctx, caseInput)
	})

	t.Run("ng:exist failed enqueue", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		caseEvents[0], _ = enterpriserule.NewTimerEvent("id1")
		caseEvents[1], _ = enterpriserule.NewTimerEvent("id2")
		caseQueueMsg := []QueueMessage{
			{
				caseEvents[0].UserId,
			},
			{
				caseEvents[1].UserId,
			},
		}
		caseOutputData := OutputData{
			NotifiedUserIdList: []string{
				caseEvents[0].UserId,
			},
			QueueMessageIdList: []string{
				"mid1",
			},
		}
		caseError := errors.New("error")

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(gomock.Eq(ctx), gomock.Eq(caseInput.EventTime)).
			Return(caseEvents, nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvents[0])).Return(caseEvents[0], nil)

		q := NewMockQueue(ctrl)
		q.EXPECT().Enqueue(gomock.Eq(caseQueueMsg[0])).Return(caseOutputData.QueueMessageIdList[0], nil)
		q.EXPECT().Enqueue(gomock.Eq(caseQueueMsg[1])).Return("", caseError)

		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(gomock.Eq(caseOutputData))

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Error(fmt.Sprintf("enqueue error user_id=%s: %s", caseEvents[1].UserId, caseError))
		log.SetDefaultLogger(l)

		interactor := &Interactor{
			repository: m,
			outputPort: o,
			queue:      q,
		}

		interactor.EnqueueEvent(ctx, caseInput)
	})

	t.Run("ng:update error", func(t *testing.T) {
		ctx := context.TODO()
		caseInput := InputData{
			EventTime: time.Now().UTC(),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		caseEvents[0], _ = enterpriserule.NewTimerEvent("id1")
		caseEvents[1], _ = enterpriserule.NewTimerEvent("id2")
		caseQueueMsg := []QueueMessage{
			{
				caseEvents[0].UserId,
			},
			{
				caseEvents[1].UserId,
			},
		}
		caseOutputData := OutputData{
			NotifiedUserIdList: []string{
				caseEvents[0].UserId,
			},
			QueueMessageIdList: []string{
				"mid1",
			},
		}
		caseError := errors.New("error")

		m := NewMockRepository(ctrl)
		m.EXPECT().FindTimerEventsByTime(gomock.Eq(ctx), gomock.Eq(caseInput.EventTime)).
			Return(caseEvents, nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvents[0])).Return(caseEvents[0], nil)
		m.EXPECT().SaveTimerEvent(gomock.Eq(ctx), gomock.Eq(caseEvents[1])).Return(nil, caseError)

		q := NewMockQueue(ctrl)
		q.EXPECT().Enqueue(gomock.Eq(caseQueueMsg[0])).Return(caseOutputData.QueueMessageIdList[0], nil)
		q.EXPECT().Enqueue(gomock.Eq(caseQueueMsg[1])).Return("dummy", nil)

		o := NewMockOutputPort(ctrl)
		o.EXPECT().Output(gomock.Eq(caseOutputData))

		l := log.NewMockLogger(ctrl)
		l.EXPECT().Error(fmt.Sprintf("update error user_id=%s: %s", caseEvents[1].UserId, caseError))
		log.SetDefaultLogger(l)

		interactor := &Interactor{
			repository: m,
			outputPort: o,
			queue:      q,
		}

		interactor.EnqueueEvent(ctx, caseInput)
	})
}
