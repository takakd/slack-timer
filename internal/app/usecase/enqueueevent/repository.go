package enqueueevent

import (
	"context"
	"slacktimer/internal/app/enterpriserule"
	"time"
)

// Make entities permanent.
type Repository interface {
	FindTimerEventsByTime(ctx context.Context, eventTime time.Time) ([]*enterpriserule.TimerEvent, error)
	SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (saved *enterpriserule.TimerEvent, err error)
}
