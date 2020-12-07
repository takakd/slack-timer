package notifyevent

import (
	"context"
	"slacktimer/internal/app/enterpriserule"
)

// Repository defines repository methods used in notification usecase.
type Repository interface {
	FindTimerEvent(ctx context.Context, userID string) (*enterpriserule.TimerEvent, error)
	SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error)
}
