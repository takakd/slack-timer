package updatetimerevent

import (
	"context"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/di"
	"slacktimer/internal/app/util/log"
	"time"
)

type Usecase interface {
	// Set notificationTime to the notification time of the event which corresponds to userId.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	UpdateNotificationTime(ctx context.Context, userId string, notificationTime time.Time, overWriteOutputPort OutputPort)

	// Set notification interval to the event which corresponds to userId.
	// Use currentTime in calculating notification time if the event is not created.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	SaveIntervalMin(ctx context.Context, userId string, currentTime time.Time, minutes int, overWriteOutputPort OutputPort)
}

type OutputData struct {
	Result     error
	SavedEvent *enterpriserule.TimerEvent
}

type OutputPort interface {
	Output(data *OutputData)
}

type Interactor struct {
	repository Repository
	//outputPort OutputPort
}

func NewUsecase() Usecase {
	return &Interactor{
		repository: di.Get("Repository").(Repository),
	}
}

// Common processing.
func (s *Interactor) saveTimerEventValue(ctx context.Context, userId string, notificationTime time.Time, remindInterval int) *OutputData {

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

	log.Debug(event.NotificationTime, event.IntervalMin, notificationTime)

	if _, err = s.repository.SaveTimerEvent(ctx, event); err != nil {
		outputData.Result = fmt.Errorf("saving timer event error userId=%v: %w", userId, err)
		return outputData
	}

	outputData.SavedEvent = event
	return outputData
}

// See Usecase interface for details.
func (s *Interactor) UpdateNotificationTime(ctx context.Context, userId string, notificationTime time.Time, overWriteOutputPort OutputPort) {
	data := s.saveTimerEventValue(ctx, userId, notificationTime, 0)
	if overWriteOutputPort != nil {
		overWriteOutputPort.Output(data)
	}
}

// Save the remind interval second for user.
// See Usecase interface for details.
func (s *Interactor) SaveIntervalMin(ctx context.Context, userId string, currentTime time.Time, minutes int, overWriteOutputPort OutputPort) {
	data := s.saveTimerEventValue(ctx, userId, currentTime, minutes)
	if overWriteOutputPort != nil {
		overWriteOutputPort.Output(data)
	}
}
