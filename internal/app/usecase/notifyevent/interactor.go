package notifyevent

import (
	"context"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
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

func (s Interactor) NotifyEvent(ctx context.Context, input InputData) error {
	outputData := OutputData{
		UserId: input.UserId,
	}

	var event *enterpriserule.TimerEvent
	event, outputData.Result = s.repository.FindTimerEvent(ctx, input.UserId)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info(fmt.Sprintf("found event %s", input.UserId))

	// Check item to be notified
	if !event.Queued() {
		log.Info(fmt.Sprintf("already notified %s", input.UserId))
		s.outputPort.Output(outputData)
		return nil
	}

	// Send notify.
	outputData.Result = s.notifier.Notify(input.UserId, input.Message)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info(fmt.Sprintf("notified %s", input.UserId))

	event.IncrementNotificationTime()
	event.SetWait()

	_, outputData.Result = s.repository.SaveTimerEvent(ctx, event)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info(fmt.Sprintf("updated event %s", input.UserId))

	outputData.Result = nil
	s.outputPort.Output(outputData)

	return nil
}
