package timeroffevent

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/appcontext"
	"slacktimer/internal/app/util/di"
)

// Interactor implements timeroffevent.InputPort.
type Interactor struct {
	repository Repository
}

var _ InputPort = (*Interactor)(nil)

// NewInteractor creates new struct.
func NewInteractor() *Interactor {
	return &Interactor{
		repository: di.Get("timeroffevent.Repository").(Repository),
	}
}

// SetEventOff start to notify event to user which corresponds to userID.
func (s Interactor) SetEventOff(ac appcontext.AppContext, input InputData, presenter OutputPort) {
	outputData := OutputData{}

	present := func() {
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

	// Set to disable event to notify.
	event.State = enterpriserule.TimerEventStateDisabled

	if _, err = s.repository.SaveTimerEvent(event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userID=%v: %w", input.UserID, err)
		present()
		return
	}

	outputData.SavedEvent = event
	present()
}
