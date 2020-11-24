package notifyevent

import (
	"context"
	"slacktimer/internal/app/enterpriserule"
)

type Repository interface {
	FindTimerEvent(ctx context.Context, userId string) (*enterpriserule.TimerEvent, error)
	SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error)
}
