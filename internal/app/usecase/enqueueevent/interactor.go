package enqueueevent

import (
	"context"
	"fmt"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Interactor implements enqueueevent.InputPort.
type Interactor struct {
	repository Repository
	outputPort OutputPort
	queue      Queue
}

// NewInteractor create new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("enqueueevent.Repository").(Repository),
		outputPort: di.Get("enqueueevent.OutputPort").(OutputPort),
		queue:      di.Get("enqueueevent.Queue").(Queue),
	}
}

var _ InputPort = (*Interactor)(nil)

// EnqueueEvent enqueues a notification event.
func (s Interactor) EnqueueEvent(ctx context.Context, data InputData) {
	outputData := OutputData{}

	events, err := s.repository.FindTimerEventsByTime(ctx, data.EventTime)
	if err != nil {
		outputData.Result = fmt.Errorf("find error time=%v: %w", data.EventTime, err)
		s.outputPort.Output(outputData)
		return
	}

	for _, e := range events {
		if e.Queued() {
			// Skip if it has already queued.
			continue
		}

		// Enqueue notification message, and send notify by other lambda corresponded queue.
		id, err := s.queue.Enqueue(QueueMessage{
			UserID: e.UserID,
			Text:   "test",
		})
		if err != nil {
			log.Error(fmt.Sprintf("enqueue error user_id=%s: %v", e.UserID, err))
			continue
		}

		// Update state.
		e.SetQueued()
		if _, err := s.repository.SaveTimerEvent(ctx, e); err != nil {
			log.Error(fmt.Sprintf("update error user_id=%s: %v", e.UserID, err))
			continue
		}

		outputData.NotifiedUserIDList = append(outputData.NotifiedUserIDList, e.UserID)
		outputData.QueueMessageIDList = append(outputData.QueueMessageIDList, id)
	}

	s.outputPort.Output(outputData)
}
