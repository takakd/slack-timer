package timeronevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
)

const (
	// ReplySuccess is message to be replied on success.
	ReplySuccess = "OK, start the notification."
	// ReplyFailure is message to be replied on failure.
	ReplyFailure = "Failed, check command syntax."
)

// Interactor implements timeronevent.InputPort.
type Interactor struct {
	repository Repository
	replier    Replier
}

var _ InputPort = (*Interactor)(nil)

// NewInteractor creates new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("timeronevent.Repository").(Repository),
		replier:    di.Get("timeronevent.Replier").(Replier),
	}
}

// SetEventOn start to notify event to user which corresponds to userID.
func (s Interactor) SetEventOn(ac appcontext.AppContext, input InputData, presenter OutputPort) {
	outputData := OutputData{}

	present := func() {
		if input.UserID != "" {
			msg := ReplySuccess
			if outputData.Result != nil {
				msg = ReplyFailure
			}

			// Reply result.
			if err := s.replier.SendMessage(ac, input.UserID, msg); err != nil {
				outputData.Result = fmt.Errorf("reply error userID=%v: %w", input.UserID, err)
			}
		}

		if presenter != nil {
			presenter.Output(ac, outputData)
		}
	}

	event, err := s.repository.FindTimerEvent(input.UserID)
	if err != nil || event == nil {
		outputData.Result = fmt.Errorf("finding timer event error userID=%v: %w", input.UserID, err)
		present()
		return
	}

	// Set to enable to notify.
	event.State = enterpriserule.TimerEventStateWait

	if _, err = s.repository.SaveTimerEvent(event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userID=%v: %w", input.UserID, err)
		present()
		return
	}

	outputData.SavedEvent = event
	present()
}
