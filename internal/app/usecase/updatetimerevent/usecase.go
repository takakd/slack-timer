package updatetimerevent

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"slacktimer/internal/app/driver/di"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/pkg/log"
	"time"
)

// Errors that this usecase returns.
var (
	ErrFind   = errors.New("could not find")
	ErrCreate = errors.New("failed to create event")
	ErrSave   = errors.New("failed to save")
)

type Usecase interface {
	// Update next notification time.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	UpdateTimeToDrink(ctx context.Context, userId string, overWriteOutputPort OutputPort)

	// Save notification interval minutes.
	// Pass OutputPort interface if overwrite presenter implementation.
	//		e.g. HTTPResponse that needs http.ResponseWrite
	SaveIntervalMin(ctx context.Context, userId string, minutes int, overWriteOutputPort OutputPort)
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
		//outputPort: di.Get("UpdateTimerEventOutputPort").(OutputPort),
	}
}

// Common processing.
func (s *Interactor) saveTimerEventValue(ctx context.Context, userId string, remindInterval int) *OutputData {

	outputData := &OutputData{}

	event, err := s.repository.FindTimerEvent(ctx, userId)
	if err != nil {
		log.Error(err)
		outputData.Result = fmt.Errorf("find %v: %w", userId, ErrFind)
		return outputData
	}

	if event == nil {
		if event, err = enterpriserule.NewTimerEvent(userId); err != nil {
			log.Error(err)
			outputData.Result = fmt.Errorf("new %v: %w", userId, ErrCreate)
			return outputData
		}
		event.UtcTimeToDrink = time.Now().UTC()
	}

	if remindInterval != 0 {
		event.DrinkTimeIntervalMin = remindInterval
	} else {
		// Set next notify time.
		event.UtcTimeToDrink = event.UtcTimeToDrink.Add(time.Duration(event.DrinkTimeIntervalMin) * time.Minute)
	}

	if _, err = s.repository.SaveTimerEvent(ctx, []*enterpriserule.TimerEvent{event}); err != nil {
		log.Error(err)
		outputData.Result = fmt.Errorf("save %v: %w", userId, ErrSave)
		return outputData
	}

	outputData.SavedEvent = event
	return outputData
}

// Update the time for user to drink.
func (s *Interactor) UpdateTimeToDrink(ctx context.Context, userId string, overWriteOutputPort OutputPort) {
	data := s.saveTimerEventValue(ctx, userId, 0)
	if overWriteOutputPort != nil {
		overWriteOutputPort.Output(data)
		//} else {
		//	s.outputPort.Output(data)
	}
}

// Save the remind interval second for user.
func (s *Interactor) SaveIntervalMin(ctx context.Context, userId string, minutes int, overWriteOutputPort OutputPort) {
	data := s.saveTimerEventValue(ctx, userId, minutes)
	if overWriteOutputPort != nil {
		overWriteOutputPort.Output(data)
		//} else {
		//	s.outputPort.Output(data)
	}
}
