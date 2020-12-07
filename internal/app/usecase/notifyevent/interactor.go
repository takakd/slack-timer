package notifyevent

import (
	"context"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
)

// Interactor implements notifyevent.InputPort.
type Interactor struct {
	outputPort OutputPort
	repository Repository
	notifier   Notifier
}

var _ InputPort = (*Interactor)(nil)

// NewInteractor create new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		outputPort: di.Get("notifyevent.OutputPort").(OutputPort),
		repository: di.Get("notifyevent.Repository").(Repository),
		notifier:   di.Get("notifyevent.Notifier").(Notifier),
	}
}

// NotifyEvent notifies events to users.
func (s Interactor) NotifyEvent(ctx context.Context, input InputData) error {
	outputData := OutputData{
		UserID: input.UserID,
	}

	var event *enterpriserule.TimerEvent
	event, outputData.Result = s.repository.FindTimerEvent(ctx, input.UserID)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info(fmt.Sprintf("found event %s", input.UserID))

	// Check item to be notified
	if !event.Queued() {
		log.Info(fmt.Sprintf("already notified %s", input.UserID))
		s.outputPort.Output(outputData)
		return nil
	}

	// Send notify.
	outputData.Result = s.notifier.Notify(input.UserID, input.Message)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info(fmt.Sprintf("notified %s", input.UserID))

	event.IncrementNotificationTime()
	event.SetWait()

	_, outputData.Result = s.repository.SaveTimerEvent(ctx, event)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info(fmt.Sprintf("updated event %s", input.UserID))

	outputData.Result = nil
	s.outputPort.Output(outputData)

	return nil
}
