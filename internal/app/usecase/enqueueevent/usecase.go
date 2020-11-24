package enqueueevent

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/pkg/log"
	"time"
)

// Errors that this usecase returns.
var (
	ErrFind = errors.New("could not find")
)

// Input port
type Usecase interface {
	// Enqueue notification event, which notification time overs eventTime.
	EnqueueEvent(ctx context.Context, eventTime time.Time) error
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
		repository: di.Get("Repository").(Repository),
		outputPort: di.Get("OutputPort").(OutputPort),
		queue:      di.Get("Queue").(Queue),
	}
}

func (s *Interactor) EnqueueEvent(ctx context.Context, eventTime time.Time) error {
	outputData := &OutputData{}

	events, err := s.repository.FindTimerEventsByTime(ctx, eventTime)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("find %v: %w", eventTime, ErrFind)
	}

	//outputData.NotifiedUserIdList = make([]string, 0)
	//outputData.QueueMessageIdList= make([]string, 0)
	for _, e := range events {
		// Enqueue notification message, and send notify by other lambda corresponded queue.
		id, err := s.queue.Enqueue(&QueueMessage{
			UserId: e.UserId,
		})
		if err != nil {
			log.Error(fmt.Sprintf("failed to enqueue user_id=%s", e.UserId))
			continue
		}

		outputData.NotifiedUserIdList = append(outputData.NotifiedUserIdList, e.UserId)
		outputData.QueueMessageIdList = append(outputData.QueueMessageIdList, id)
	}
	s.outputPort.Output(outputData)

	return nil
}
