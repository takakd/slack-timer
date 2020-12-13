package notifyevent

import (
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
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

// NewInteractor creates new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		outputPort: di.Get("notifyevent.OutputPort").(OutputPort),
		repository: di.Get("notifyevent.Repository").(Repository),
		notifier:   di.Get("notifyevent.Notifier").(Notifier),
	}
}

// NotifyEvent notifies events to users.
func (s Interactor) NotifyEvent(ac appcontext.AppContext, input InputData) error {
	outputData := OutputData{
		UserID: input.UserID,
	}

	logDetail := map[string]interface{}{
		"user_id": input.UserID,
	}

	var event *enterpriserule.TimerEvent
	event, outputData.Result = s.repository.FindTimerEvent(input.UserID)
	if outputData.Result != nil {
		s.outputPort.Output(ac, outputData)
		return outputData.Result
	}

	log.InfoWithContext(ac, "found event", logDetail)

	// Check item to be notified
	if event.State != enterpriserule.TimerEventStateQueued {
		log.InfoWithContext(ac, "already notified", logDetail)
		s.outputPort.Output(ac, outputData)
		return nil
	}

	// Send notify.
	outputData.Result = s.notifier.SendMessage(ac, input.UserID, input.Message)
	if outputData.Result != nil {
		s.outputPort.Output(ac, outputData)
		return outputData.Result
	}

	log.InfoWithContext(ac, "notified", logDetail)

	event.IncrementNotificationTime()
	event.State = enterpriserule.TimerEventStateWait

	_, outputData.Result = s.repository.SaveTimerEvent(event)
	if outputData.Result != nil {
		s.outputPort.Output(ac, outputData)
		return outputData.Result
	}

	log.InfoWithContext(ac, "updated event", logDetail)

	outputData.Result = nil
	s.outputPort.Output(ac, outputData)

	return nil
}
