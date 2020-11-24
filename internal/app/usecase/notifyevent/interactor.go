package notifyevent

import (
	"context"
	"github.com/pkg/errors"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/enterpriserule"
)

// Errors this usecase returns.
var (
	ErrFind = errors.New("could not find")
)

type Interactor struct {
	outputPort OutputPort
	repository Repository
	notifier   Notifier
}

func NewInteractor() InputPort {
	return &Interactor{
		outputPort: di.Get("notifyevent.OutputPort").(OutputPort),
		repository: di.Get("notifyevent.Repository").(Repository),
		notifier:   di.Get("notifyevent.Notifier").(Notifier),
	}
}

func (s *Interactor) NotifyEvent(ctx context.Context, input *InputData) error {
	outputData := &OutputData{
		UserId: input.UserId,
	}

	// Send notify.
	outputData.Result = s.notifier.Notify(input.UserId, input.Message)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	// Update time.
	var event *enterpriserule.TimerEvent
	event, outputData.Result = s.repository.FindTimerEvent(ctx, input.UserId)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	event.IncrementNotificationTime()

	_, outputData.Result = s.repository.SaveTimerEvent(ctx, event)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	outputData.Result = nil
	s.outputPort.Output(outputData)

	return nil
}
