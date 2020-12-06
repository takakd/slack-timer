package updatetimerevent

import (
	"context"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"time"
)

type Interactor struct {
	repository Repository
}

func NewInteractor() InputPort {
	return &Interactor{
		repository: di.Get("updatetimerevent.Repository").(Repository),
	}
}

// Common processing.
func (s Interactor) saveTimerEventValue(ctx context.Context, userId string, notificationTime time.Time, remindInterval int) *OutputData {

	outputData := &OutputData{}

	event, err := s.repository.FindTimerEvent(ctx, userId)
	if err != nil {
		outputData.Result = fmt.Errorf("finding timer event error userId=%v: %w", userId, err)
		return outputData
	}

	if event == nil {
		if event, err = enterpriserule.NewTimerEvent(userId); err != nil {
			outputData.Result = fmt.Errorf("creating timer event error userId=%v: %w", userId, err)
			return outputData
		}
	}

	if remindInterval != 0 {
		event.IntervalMin = remindInterval
	}

	// Set next notify time.
	event.NotificationTime = notificationTime.Add(time.Duration(event.IntervalMin) * time.Minute)

	if _, err = s.repository.SaveTimerEvent(ctx, event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userId=%v: %w", userId, err)
		return outputData
	}

	outputData.SavedEvent = event
	return outputData
}

func (s Interactor) UpdateNotificationTime(ctx context.Context, userId string, notificationTime time.Time, presenter OutputPort) {
	data := s.saveTimerEventValue(ctx, userId, notificationTime, 0)
	if presenter != nil {
		presenter.Output(*data)
	}
}

func (s Interactor) SaveIntervalMin(ctx context.Context, userId string, currentTime time.Time, minutes int, presetner OutputPort) {
	data := s.saveTimerEventValue(ctx, userId, currentTime, minutes)
	if presetner != nil {
		presetner.Output(*data)
	}
}
