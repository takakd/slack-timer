package updatetimerevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"time"
)

const (
	// ReplySuccess is message to be replied on success.
	ReplySuccess = "OK, set message."
	// ReplyFailure is message to be replied on failure.
	ReplyFailure = "Failed, Check command syntax."
)

// Interactor implements updatetimerevent.InputPort.
type Interactor struct {
	repository Repository
	replier    Replier
}

var _ InputPort = (*Interactor)(nil)

// NewInteractor creates new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("updatetimerevent.Repository").(Repository),
		replier:    di.Get("updatetimerevent.Replier").(Replier),
	}
}

// SaveIntervalMin sets notification interval to the event which corresponds to userID.
func (s Interactor) SaveIntervalMin(ac appcontext.AppContext, input SaveEventInputData, presenter OutputPort) {
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
	if err != nil {
		outputData.Result = fmt.Errorf("finding timer event error userID=%v: %w", input.UserID, err)
		present()
		return
	}

	if event == nil {
		if event, err = enterpriserule.NewTimerEvent(input.UserID); err != nil {
			outputData.Result = fmt.Errorf("creating timer event error userID=%v: %w", input.UserID, err)
			present()
			return
		}
	}

	if input.Minutes != 0 {
		event.IntervalMin = input.Minutes
	}

	// Update values.
	event.NotificationTime = input.CurrentTime.Add(time.Duration(event.IntervalMin) * time.Minute)
	event.Text = input.Text

	if _, err = s.repository.SaveTimerEvent(event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userID=%v: %w", input.UserID, err)
		present()
		return
	}

	outputData.SavedEvent = event
	present()
}
