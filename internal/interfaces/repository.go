package interfaces

import (
	"context"
	"proteinreminder/internal/entity"
	"time"
)

type Repository interface {
	// ProteinEvent
	FindProteinEvent(ctx context.Context, userId string) (*entity.ProteinEvent, error)
	FindProteinEventByTime(ctx context.Context, from, to time.Time) ([]*entity.ProteinEvent, error)
	SaveProteinEvent(ctx context.Context, events []*entity.ProteinEvent) ([]*entity.ProteinEvent, error)
}
