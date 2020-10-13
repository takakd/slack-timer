package usecase

import (
	"context"
	"github.com/pkg/errors"
	"proteinreminder/internal/app/apprule"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/log"
	"time"
)

// The type of Error that this usecase returns.
type SaveProteinEventError int

// Errors that this usecase returns.
var (
	ErrFind   = errors.New("could not find")
	ErrCreate = errors.New("failed to create event")
	ErrSave   = errors.New("failed to save")
)

type ProteinEventSaver interface {
	SaveTimeToDrink(ctx context.Context, userId string, timeToDrink time.Time) error
	SaveIntervalSec(ctx context.Context, userId string, minutes time.Duration) error
}

type SaveProteinEvent struct {
	repository apprule.Repository
}

func NewSaveProteinEvent(repository apprule.Repository) (*SaveProteinEvent, error) {
	if repository == nil {
		return nil, errors.New("repository must not be null")
	}
	return &SaveProteinEvent{
		repository: repository,
	}, nil
}

// Common processing.
func (s *SaveProteinEvent) saveProteinEventValue(ctx context.Context, userId string, toDrink *time.Time, remindInterval *time.Duration) error {

	event, err := s.repository.FindProteinEvent(ctx, userId)
	if err != nil {
		log.Error(err)
		return ErrFind
	}

	if event == nil {
		if event, err = enterpriserule.NewProteinEvent(userId); err != nil {
			return ErrCreate
		}
	}

	if toDrink != nil {
		event.UtcTimeToDrink = *toDrink
	}
	if remindInterval != nil {
		event.DrinkTimeIntervalSec = *remindInterval
	}

	if _, err = s.repository.SaveProteinEvent(ctx, []*enterpriserule.ProteinEvent{event}); err != nil {
		return ErrSave
	}

	return nil
}

// Save the time for user to drink.
func (s *SaveProteinEvent) SaveTimeToDrink(ctx context.Context, userId string, timeToDrink time.Time) error {
	return s.saveProteinEventValue(ctx, userId, &timeToDrink, nil)
}

// Save the remind interval second for user.
func (s *SaveProteinEvent) SaveIntervalSec(ctx context.Context, userId string, minutes time.Duration) error {
	return s.saveProteinEventValue(ctx, userId, nil, &minutes)
}
