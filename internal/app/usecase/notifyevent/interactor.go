package notifyevent

import (
	"context"
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

	logDetail := map[string]interface{}{
		"user_id": input.UserID,
	}

	var event *enterpriserule.TimerEvent
	event, outputData.Result = s.repository.FindTimerEvent(ctx, input.UserID)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info("found event", logDetail)

	// Check item to be notified
	if !event.Queued() {
		log.Info("already notified", logDetail)
		s.outputPort.Output(outputData)
		return nil
	}

	// Send notify.
	outputData.Result = s.notifier.Notify(input.UserID, input.Message)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info("notified", logDetail)

	event.IncrementNotificationTime()
	event.SetWait()

	_, outputData.Result = s.repository.SaveTimerEvent(ctx, event)
	if outputData.Result != nil {
		s.outputPort.Output(outputData)
		return outputData.Result
	}

	log.Info("updated event", logDetail)

	outputData.Result = nil
	s.outputPort.Output(outputData)

	return nil
}
