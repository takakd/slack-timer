package apprule

import (
	"context"
	"proteinreminder/internal/app/enterpriserule"
	"time"
)

// Make entities permanent.
type Repository interface {
	FindProteinEvent(ctx context.Context, userId string) (*enterpriserule.ProteinEvent, error)
	FindProteinEventByTime(ctx context.Context, from, to time.Time) ([]*enterpriserule.ProteinEvent, error)
	SaveProteinEvent(ctx context.Context, events []*enterpriserule.ProteinEvent) ([]*enterpriserule.ProteinEvent, error)
}
