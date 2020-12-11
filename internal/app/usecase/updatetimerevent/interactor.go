package updatetimerevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
	"time"
)

// Interactor implements updatetimerevent.InputPort.
type Interactor struct {
	repository Repository
}

var _ InputPort = (*Interactor)(nil)

// NewInteractor creates new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("updatetimerevent.Repository").(Repository),
	}
}

// Common processing.
func (s Interactor) saveTimerEventValue(userID string, notificationTime time.Time, remindInterval int, text string) *OutputData {

	outputData := &OutputData{}

	event, err := s.repository.FindTimerEvent(userID)
	if err != nil {
		outputData.Result = fmt.Errorf("finding timer event error userID=%v: %w", userID, err)
		return outputData
	}

	if event == nil {
		if event, err = enterpriserule.NewTimerEvent(userID); err != nil {
			outputData.Result = fmt.Errorf("creating timer event error userID=%v: %w", userID, err)
			return outputData
		}
	}

	if remindInterval != 0 {
		event.IntervalMin = remindInterval
	}

	// Update values.
	event.NotificationTime = notificationTime.Add(time.Duration(event.IntervalMin) * time.Minute)
	event.Text = text

	if _, err = s.repository.SaveTimerEvent(event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userID=%v: %w", userID, err)
		return outputData
	}

	outputData.SavedEvent = event
	return outputData
}

// UpdateNotificationTime sets notificationTime to the notification time of the event which corresponds to userID.
func (s Interactor) UpdateNotificationTime(ac appcontext.AppContext, input UpdateNotificationTimeInputData, presenter OutputPort) {
	data := s.saveTimerEventValue(input.UserID, input.NotificationTime, 0, "")
	if presenter != nil {
		presenter.Output(ac, *data)
	}
}

// SaveIntervalMin sets notification interval to the event which corresponds to userID.
func (s Interactor) SaveIntervalMin(ac appcontext.AppContext, input SaveEventInputData, presenter OutputPort) {
	data := s.saveTimerEventValue(input.UserID, input.CurrentTime, input.Minutes, input.Text)
	if presenter != nil {
		presenter.Output(ac, *data)
	}
}
