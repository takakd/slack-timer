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

// NewInteractor create new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("updatetimerevent.Repository").(Repository),
	}
}

// Common processing.
func (s Interactor) saveTimerEventValue(userID string, notificationTime time.Time, remindInterval int) *OutputData {

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

	// Set next notify time.
	event.NotificationTime = notificationTime.Add(time.Duration(event.IntervalMin) * time.Minute)

	if _, err = s.repository.SaveTimerEvent(event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userID=%v: %w", userID, err)
		return outputData
	}

	outputData.SavedEvent = event
	return outputData
}

// UpdateNotificationTime sets notificationTime to the notification time of the event which corresponds to userID.
func (s Interactor) UpdateNotificationTime(ac appcontext.AppContext, userID string, notificationTime time.Time, presenter OutputPort) {
	data := s.saveTimerEventValue(userID, notificationTime, 0)
	if presenter != nil {
		presenter.Output(ac, *data)
	}
}

// SaveIntervalMin sets notification interval to the event which corresponds to userID.
func (s Interactor) SaveIntervalMin(ac appcontext.AppContext, userID string, currentTime time.Time, minutes int, presetner OutputPort) {
	data := s.saveTimerEventValue(userID, currentTime, minutes)
	if presetner != nil {
		presetner.Output(ac, *data)
	}
}
