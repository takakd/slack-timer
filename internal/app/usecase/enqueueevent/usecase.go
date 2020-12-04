package enqueueevent

import (
	"context"
	"fmt"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

// Input port
type Usecase interface {
	// Enqueue notification event, which notification time overs eventTime.
	EnqueueEvent(ctx context.Context, eventTime time.Time)
}

type OutputData struct {
	Result             error
	NotifiedUserIdList []string
	QueueMessageIdList []string
}

type OutputPort interface {
	Output(data *OutputData)
}

type Interactor struct {
	repository Repository
	outputPort OutputPort
	queue      Queue
}

func NewUsecase() Usecase {
	return &Interactor{
		repository: di.Get("enqueueevent.Repository").(Repository),
		outputPort: di.Get("enqueueevent.OutputPort").(OutputPort),
		queue:      di.Get("enqueueevent.Queue").(Queue),
	}
}

func (s *Interactor) EnqueueEvent(ctx context.Context, eventTime time.Time) {
	outputData := &OutputData{}

	events, err := s.repository.FindTimerEventsByTime(ctx, eventTime)
	if err != nil {
		outputData.Result = fmt.Errorf("find error time=%v: %w", eventTime, err)
		s.outputPort.Output(outputData)
		return
	}

	for _, e := range events {
		if e.Queued() {
			// Skip if it has aleady queued.
			continue
		}

		// Enqueue notification message, and send notify by other lambda corresponded queue.
		id, err := s.queue.Enqueue(&QueueMessage{
			UserId: e.UserId,
		})
		if err != nil {
			log.Error(fmt.Sprintf("enqueue error user_id=%s: %v", e.UserId, err))
			continue
		}

		// Update state.
		e.SetQueued()
		if _, err := s.repository.SaveTimerEvent(ctx, e); err != nil {
			log.Error(fmt.Sprintf("update error user_id=%s: %v", e.UserId, err))
			continue
		}

		outputData.NotifiedUserIdList = append(outputData.NotifiedUserIdList, e.UserId)
		outputData.QueueMessageIdList = append(outputData.QueueMessageIdList, id)
	}

	s.outputPort.Output(outputData)
}
