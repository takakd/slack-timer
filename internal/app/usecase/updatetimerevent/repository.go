package updatetimerevent

import (
	"context"
	"slacktimer/internal/app/enterpriserule"
	"time"
)

// Make entities permanent.
type Repository interface {
	FindTimerEvent(ctx context.Context, userId string) (*enterpriserule.TimerEvent, error)
	FindTimerEventByTime(ctx context.Context, from, to time.Time) ([]*enterpriserule.TimerEvent, error)
	SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error)
}
