package updatetimerevent

import (
	"context"
	"slacktimer/internal/app/enterpriserule"
	"time"
)

// Repository defines repository methods used in updating timer events usecase.
type Repository interface {
	FindTimerEvent(ctx context.Context, userID string) (*enterpriserule.TimerEvent, error)
	FindTimerEventByTime(ctx context.Context, from, to time.Time) ([]*enterpriserule.TimerEvent, error)
	SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error)
}
