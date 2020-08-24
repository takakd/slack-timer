package usecase

import (
	"context"
	"proteinreminder/internal/entity"
	"proteinreminder/internal/interfaces"
	"proteinreminder/internal/ioc"
	"time"
)

const (
	SetProteinEventNoError = iota
	SetProteinEventErrorFind
	SetProteinEventErrorCreate
	SetProteinEventErrorSave
)

type SetProteinEventError int

type SetProteinEvent struct {
	repository interfaces.Repository
}

func NewSetProteinEvent(repository interfaces.Repository) *SetProteinEvent {
	return &SetProteinEvent{
		repository: repository,
	}
}

type SetProteinEventValueArgs struct {
	timeToDrink         *time.Time
	remindIntervalInMin *time.Duration
}

// Common processing.
func (s *SetProteinEvent) setProteinEventValue(ctx context.Context, userId string, args *SetProteinEventValueArgs) SetProteinEventError {

	logger := ioc.GetLogger()

	// テストがしずらいので、FindProteinEventを持つinterfaceを受け取る
	event, err := s.repository.FindProteinEvent(ctx, userId)
	if err != nil {
		logger.Error("%v", err.Error())
		return SetProteinEventErrorFind
	}

	if event != nil {
		if event, err = entity.NewProteinEvent(userId); err != nil {
			return SetProteinEventErrorCreate
		}
	}

	if args.timeToDrink != nil {
		event.UtcTimeToDrink = *args.timeToDrink
	}
	if args.remindIntervalInMin != nil {
		event.DrinkTimeIntervalSec = *args.remindIntervalInMin
	}

	// Save
	if _, err = s.repository.SaveProteinEvent(ctx, []*entity.ProteinEvent{event}); err != nil {
		return SetProteinEventErrorSave
	}

	return SetProteinEventNoError
}

// Set the time to drink.
func (s *SetProteinEvent) SetProteinEventTimeToDrink(ctx context.Context, userId string, timeToDrink time.Time) SetProteinEventError {
	args := &SetProteinEventValueArgs{
		timeToDrink: &timeToDrink,
	}
	return s.setProteinEventValue(ctx, userId, args)
}

// Set the remind interval second.
func (s *SetProteinEvent) SetProteinEventIntervalSec(ctx context.Context, userId string, minutes time.Duration) SetProteinEventError {
	args := &SetProteinEventValueArgs{
		remindIntervalInMin: &minutes,
	}
	return s.setProteinEventValue(ctx, userId, args)
}
