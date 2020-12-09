package updatetimerevent

import (
	"slacktimer/internal/app/enterpriserule"
	"time"
)

// Repository defines repository methods used in updating timer events usecase.
type Repository interface {
	FindTimerEvent(userID string) (*enterpriserule.TimerEvent, error)
	FindTimerEventByTime(from, to time.Time) ([]*enterpriserule.TimerEvent, error)
	SaveTimerEvent(event *enterpriserule.TimerEvent) (*enterpriserule.TimerEvent, error)
}
