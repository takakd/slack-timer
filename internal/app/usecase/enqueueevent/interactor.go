package enqueueevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Interactor implements enqueueevent.InputPort.
type Interactor struct {
	repository Repository
	outputPort OutputPort
	queue      Queue
}

// NewInteractor creates new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("enqueueevent.Repository").(Repository),
		outputPort: di.Get("enqueueevent.OutputPort").(OutputPort),
		queue:      di.Get("enqueueevent.Queue").(Queue),
	}
}

var _ InputPort = (*Interactor)(nil)

// EnqueueEvent enqueues a notification event.
func (s Interactor) EnqueueEvent(ac appcontext.AppContext, data InputData) {
	outputData := OutputData{}

	events, err := s.repository.FindTimerEventsByTime(data.EventTime)
	if err != nil {
		outputData.Result = fmt.Errorf("find error time=%v: %w", data.EventTime, err)
		s.outputPort.Output(ac, outputData)
		return
	}

	for _, e := range events {
		if e.State != enterpriserule.TimerEventStateWait {
			// Skip if the event is enqueued or disabled.
			continue
		}

		// Enqueue notification message, and send notify by other lambda corresponded queue.
		id, err := s.queue.Enqueue(QueueMessage{
			UserID: e.UserID(),
			Text:   e.Text,
		})
		if err != nil {
			log.ErrorWithContext(ac, fmt.Sprintf("enqueue error user_id=%s: %v", e.UserID(), err))
			continue
		}

		// Update state.
		e.State = enterpriserule.TimerEventStateQueued
		if _, err := s.repository.SaveTimerEvent(e); err != nil {
			log.ErrorWithContext(ac, fmt.Sprintf("update error user_id=%s: %v", e.UserID(), err))
			continue
		}

		outputData.NotifiedUserIDList = append(outputData.NotifiedUserIDList, e.UserID())
		outputData.QueueMessageIDList = append(outputData.QueueMessageIDList, id)
	}

	s.outputPort.Output(ac, outputData)
}
