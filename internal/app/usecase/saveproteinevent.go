package usecase

import (
	"context"
	"fmt"
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
	UpdateTimeToDrink(ctx context.Context, userId string) error
	SaveIntervalMin(ctx context.Context, userId string, minutes int) error
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
func (s *SaveProteinEvent) saveProteinEventValue(ctx context.Context, userId string, remindInterval int) error {

	event, err := s.repository.FindProteinEvent(ctx, userId)
	if err != nil {
		log.Error(err)
		return fmt.Errorf("find %v: %w", userId, ErrFind)
	}

	if event == nil {
		if event, err = enterpriserule.NewProteinEvent(userId); err != nil {
			log.Error(err)
			return fmt.Errorf("new %v: %w", userId, ErrCreate)
		}
		event.UtcTimeToDrink = time.Now().UTC()
	}

	if remindInterval != 0 {
		event.DrinkTimeIntervalMin = remindInterval
	} else {
		// Set next notify time.
		event.UtcTimeToDrink = event.UtcTimeToDrink.Add(time.Duration(event.DrinkTimeIntervalMin) * time.Minute)
	}

	if _, err = s.repository.SaveProteinEvent(ctx, []*enterpriserule.ProteinEvent{event}); err != nil {
		log.Error(err)
		return fmt.Errorf("save %v: %w", userId, ErrSave)
	}

	return nil
}

// Update the time for user to drink.
func (s *SaveProteinEvent) UpdateTimeToDrink(ctx context.Context, userId string) error {
	return s.saveProteinEventValue(ctx, userId, 0)
}

// Save the remind interval second for user.
func (s *SaveProteinEvent) SaveIntervalMin(ctx context.Context, userId string, minutes int) error {
	return s.saveProteinEventValue(ctx, userId, minutes)
}
